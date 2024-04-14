package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Noah-Wilderom/foreverstore/p2p"

	"github.com/joho/godotenv"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
