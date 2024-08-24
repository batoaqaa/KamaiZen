package utils

import (
	"log"
	"os"
)

// global logger
var logger *log.Logger

func GetLogger() *log.Logger {
	if logger == nil {
		filename := "/home/ibrahim/work/KamaiZen/kamaizen.log"
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		logger = log.New(file, "[KamaiZen]", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return logger
}
