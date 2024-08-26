package lsp

import (
	"KamaiZen/rpc"
	"io"
	"os"
)

var writer io.Writer

func WriteResponse(response interface{}) {
	reply := rpc.EncodeMessage(response)
	writer.Write([]byte(reply))
}

func Write(message []byte) {
	writer.Write(message)
}

func Initialise() {
	writer = os.Stdout
}
