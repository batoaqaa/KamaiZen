package rpc_test

import (
	"KamaiZen/rpc"
	"testing"
)

type EncodingExample struct {
	Method string
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 18\r\n\r\n{\"Method\":\"hello\"}"
	value := EncodingExample{Method: "hello"}
	actual := rpc.EncodeMessage(value)
	if actual != expected {
		t.Fatalf("Expected: %s,\ngot: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	expectedContentLength := 18
	value := []byte("Content-Length: 18\r\n\r\n{\"method\":\"hello\"}")
	method, content, err := rpc.DecodeMessage(value)
	contentLength := len(content)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if contentLength != expectedContentLength {
		t.Fatalf("Expected: %d,\ngot: %d", expectedContentLength, contentLength)
	}
	if method != "hello" {
		t.Fatalf("Expected: hello,\ngot: %s", method)
	}
}
