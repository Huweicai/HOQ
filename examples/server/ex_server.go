package main

import (
	"HOQ/hoq"
	"log"
)

func main() {
	server, _ := hoq.NewServer(hoq.EngineQuic, hoq.EchoHandler)
	err := server.Run("127.0.0.1:8787")
	if err != nil {
		log.Fatalln(err.Error())
	}
}
