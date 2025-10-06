package main

import (
	"context"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"sync/atomic"

	"gin-demo/config"
	"gin-demo/database"
	"gin-demo/handlers/healthcheckhdlr"
	"gin-demo/handlers/userhdlr"
	"gin-demo/middlewares"
	"gin-demo/server"
	"gin-demo/services"
	"gin-demo/services/healthchecksvc"
	"gin-demo/services/usersvc"
	"gin-demo/utils"

	"github.com/akshay-glide/bivo-utils/kafka"
	"github.com/akshay-glide/bivo-utils/postgres"
)

func main() {
	// Use only 2 OS threads
	runtime.GOMAXPROCS(2)

	// Load server config
	serverConfig, err := config.GetServerConfig()
	if err != nil {
		log.Fatal("Server Config Error: ", err)
	}

	// Setup logging
	logDirPath := path.Join(*serverConfig.ScratchDir, "logs")
	utils.SetupLogger(logDirPath)

	apilogger := log.New(os.Stdout, "", log.Lmicroseconds|log.LstdFlags|log.Llongfile)
	utils.AddLoggerFile(logDirPath, "apiserver.log", apilogger)

	servlogger := log.New(os.Stdout, "", log.Lmicroseconds|log.LstdFlags|log.Llongfile)
	utils.AddLoggerFile(logDirPath, "service.log", servlogger)

	hdlrlogger := log.New(os.Stdout, "", log.Lmicroseconds|log.LstdFlags|log.Llongfile)
	utils.AddLoggerFile(logDirPath, "handler.log", hdlrlogger)

	// ########## Database ##########
	//postgresDSN := utils.GetPostgresDSN(serverConfig.PostgresConfig)

	// 1. ######## Services ########
	postgresDbSvcI, err := postgres.NewPostgresDB(serverConfig.PostgresConfig)
	if err != nil {
		apilogger.Fatal("DB instantiation failed")
	}
	err = database.AutoMigrateAll(postgresDbSvcI)
	if err != nil {
		apilogger.Fatal("DB AutoMigrate failed: ", err)
	}
	apilogger.Println("Postgres connected and migrated")

	// ########## Kafka Producer and Consumer ##########
	kafkaProducer, err := kafka.NewKafkaProducer(&serverConfig.KafkaConfig)
	if err != nil {
		apilogger.Fatal("Kafka Producer instantiation failed: ", err)
	}
	defer kafkaProducer.Close()

	consumer, err := kafka.NewKafkaConsumer(&serverConfig.KafkaConfig)
	if err != nil {
		apilogger.Fatal("Kafka Consumer instantiation failed: ", err)
	}
	defer consumer.Close()

	if err := consumer.Subscribe(); err != nil {
		apilogger.Fatal("Failed to subscribe to topics:", err)
	}

	handler := services.NewMessageHandler()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer.Start(ctx, handler.HandleMessage)

	// ########## Services ##########
	healthCheckSvcI := healthchecksvc.NewHealthCheckSvc()
	userSvcI := usersvc.NewUserService(postgresDbSvcI, *kafkaProducer, serverConfig.KafkaConfig.Topic, servlogger)

	// ########## Handlers ##########
	healthCheckHdlrI := healthcheckhdlr.NewHealthCheckHdlr(healthCheckSvcI)
	userHdlrI := userhdlr.NewUserHdlr(userSvcI, hdlrlogger)

	handlers := []*server.ServerHandlerMap{
		server.NewServerHandlerMap("/api/v1/health", healthCheckHdlrI),
		server.NewServerHandlerMap("/api/v1/user", userHdlrI),
	}

	// ########## Gin App + Middlewares ##########
	ginApp := server.GetGinApplication()
	middlewares.AddGinMiddlewares(ginApp)

	// ########## Start Server ##########
	serv := server.GetServer(serverConfig.APIServerConfig, ginApp, handlers)
	waitForShutdownInterrupt := serv.StartServer()

	// Setup stop file listener
	stoppedFlag := uint32(0)
	stopFileCh := make(chan string, 1)
	pendingThreads := &sync.WaitGroup{}
	pendingThreads.Add(1)

	go func() {
		utils.WaitTillStopFile(&stoppedFlag, stopFileCh, path.Join(*serverConfig.ScratchDir, "stopfile"))
		atomic.StoreUint32(&stoppedFlag, 1)
		pendingThreads.Done()
	}()

	// Wait for signal or stop file
	select {
	case <-waitForShutdownInterrupt:
		apilogger.Println("Interrupt received, initiating shutdown...")
	case <-stopFileCh:
		apilogger.Println("Stop file detected, initiating shutdown...")
	}

	atomic.StoreUint32(&stoppedFlag, 1)

	apilogger.Println("Shutting Down Server...")
	serv.ShutdownGracefully()
	apilogger.Println("Server Down")
}
