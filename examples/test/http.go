package main

import (
	"HOQ/logs"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	println(path)
	err := http.ListenAndServeTLS("localhost:6666", "/Users/hwc/go/prjs/HOQ/cert/cert.pem", "/Users/hwc/go/prjs/HOQ/cert/key.pem", MyHanler{})
	if err != nil {
		logs.Error(err)
	}
}

type MyHanler struct {
}

func (MyHanler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host + " " + r.Method)
	w.Write([]byte("OK"))
}
