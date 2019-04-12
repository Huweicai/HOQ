package hoq

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/lucas-clemente/quic-go"
	_ "github.com/lucas-clemente/quic-go"
	"io/ioutil"
	"log"
	"math/big"
)

type Server struct {
}

func handleRequest(stream quic.Stream) {
	got, err := ioutil.ReadAll(stream)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	stream.Write(testText)
	log.Println("GOT" + string(got))
	stream.Close()
}

func (s *Server) Run(addr string) error {
	listen, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	log.Println("HTTP server started ,listen at " + addr)
	for {
		sess, err := listen.Accept()
		if err != nil {
			return err
		}
		stream, err := sess.AcceptStream()
		if err != nil {
			return err
		}
		go handleRequest(stream)
	}
}

//todo learn more about it's encryption
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
