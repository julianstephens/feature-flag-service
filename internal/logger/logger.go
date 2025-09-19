package logger

import (
	"log"
	"os"
	"sync"
)

var (
	logger *log.Logger
	once   sync.Once
)

func GetLogger() *log.Logger {
	once.Do(func() {
		logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	})
	return logger
}

func SetOutput(output *os.File) {
	GetLogger().SetOutput(output)
}

func SetPrefix(prefix string) {
	GetLogger().SetPrefix(prefix)
}

func SetFlags(flags int) {
	GetLogger().SetFlags(flags)
}
