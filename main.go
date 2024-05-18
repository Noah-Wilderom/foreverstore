package main

import (
	"bytes"
	"fmt"
	"io"
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
		EncKey:            newEncryptionKey(),
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
	s3 := makeServer(":5000", ":3000", ":4000")

	go func() {
		log.Fatal(s1.Start())
	}()
	go func() {
		log.Fatal(s2.Start())
	}()

	time.Sleep(2 * time.Second)

	go s3.Start()
	time.Sleep(2 * time.Second)

	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("picture_%d.png", i)
		data := bytes.NewReader([]byte("my big data file here"))
		s3.Store(key, data)
		time.Sleep(5 * time.Millisecond)

		if err := s3.store.Delete(s3.Id, key); err != nil {
			log.Fatal(err)
		}

		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}
}
