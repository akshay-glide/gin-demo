package utils

import (
	"log"
	"os"
	"sync/atomic"
	"time"
)

func WaitTillStopFile(stoppedflag *uint32, stopch chan string, stopfilepath string) {
	log.Println("Stopfile Path:", stopfilepath)
	for atomic.LoadUint32(stoppedflag) == 0 {
		if _, err := os.Stat(stopfilepath); err == nil {
			log.Println("Removing stopfile")
			err := os.Remove(stopfilepath)
			if err != nil {
				log.Println("Removing stopfile has failed")
			}
			break
		} else {
			time.Sleep(time.Millisecond * 1000)
		}
	}
	stopch <- "stopfilereceived"
}
