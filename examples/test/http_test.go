package main

import (
	"HOQ/logs"
	"net/http"
	"testing"
)

func TestMyHanler_ServeHTTP(t *testing.T) {
	resp, err := http.Get("https://localhost")
	logs.Error(resp, err)
}
