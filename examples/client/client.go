package main

import (
	"bufio"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"log"
)

func main() {
	sess, err := quic.DialAddr("127.0.0.1:8787", &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	stream, err := sess.OpenStreamSync()
	if err != nil {
		log.Fatal(err)
		return
	}
	sendText := "Hello World\n"
	_, err = stream.Write([]byte(sendText))
	log.Println("SEND: " + sendText)
	if err != nil {
		log.Fatal(err)
		return
	}
	nr := bufio.NewReader(stream)
	got, _, err := nr.ReadRune()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("GOT: " + string(got))
}
