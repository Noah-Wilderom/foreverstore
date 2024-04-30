package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Noah-Wilderom/foreverstore/p2p"

	"github.com/joho/godotenv"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr[1:] + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s1 := makeServer(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(2 * time.Second)

	go s2.Start()
	time.Sleep(2 * time.Second)

	for i := 0; i < 10; i++ {
		data := bytes.NewReader([]byte("my big data file here"))
		s2.Store("myprivatekey", data)
		time.Sleep(500 * time.Millisecond)
	}

	// r, err := s2.Get("myprivatekey")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// b, err := io.ReadAll(r)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Println(string(b))

	select {}
}
