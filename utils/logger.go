package utils

import (
	"log"
	"path"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger(logDirPath string) {
	logFilePath := path.Join(logDirPath, "general.log")
	lumberjackLogger := &lumberjack.Logger{
		// Log file abbsolute path, os agnostic
		Filename:   filepath.ToSlash(logFilePath),
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	log.SetFlags(log.Ltime | log.Ldate | log.Lmicroseconds | log.Lshortfile)

	log.SetOutput(lumberjackLogger)
}

func AddLoggerFile(logDir string, logFile string, logger *log.Logger) {
	lumberjackLogger := &lumberjack.Logger{
		// Log file abbsolute path, os agnostic
		Filename:   filepath.ToSlash(path.Join(logDir, logFile)),
		MaxSize:    5, // MB
		MaxBackups: 5,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}
	logger.SetFlags(log.Ltime | log.Ldate | log.Lmicroseconds | log.Lshortfile)
	logger.SetOutput(lumberjackLogger)
}
