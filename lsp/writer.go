package lsp

import (
	"KamaiZen/rpc"
	"io"
)

func WriteResponse(writer io.Writer, response interface{}) {
	reply := rpc.EncodeMessage(response)
	writer.Write([]byte(reply))
}
