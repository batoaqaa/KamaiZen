package lsp

import (
	"KamaiZen/logger"
	"KamaiZen/rpc"
	"io"
	"os"
	"sync"
)

var writer io.Writer

const buffered_channel_size = 24

var writer_channel = make(chan []byte, buffered_channel_size)

func WriteResponse(response interface{}) {
	reply := rpc.EncodeMessage(response)
	writer_channel <- []byte(reply)
}

func Write(message []byte) {
	writer.Write(message)
}

func Start(wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info("Starting writer")
	for {
		select {
		case message := <-writer_channel:
			Write(message)
		}
	}
}

func Initialise() {
	writer = os.Stdout
}
