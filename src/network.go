package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const protocolID = "RapidTransfer"

func server(node host.Host) {
	node.SetStreamHandler(protocolID, func(s network.Stream) {
		go writeCounter(s)
		go readCounter(s)
	})
	key := fmt.Sprintf("%s/p2p/%s", node.Addrs()[1], node.ID())
	fmt.Println(key)
}

func client(node host.Host, peerAddr string) {
	// Parse the multiaddr string.
	peerMA, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		panic(err)
	}
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	if err != nil {
		panic(err)
	}

	// Connect to given address
	if err := node.Connect(context.Background(), *peerAddrInfo); err != nil {
		panic(err)
	}
	fmt.Println("Connected to", peerAddrInfo.String())

	// Create a stream with peer
	s, err := node.NewStream(context.Background(), peerAddrInfo.ID, protocolID)
	if err != nil {
		panic(err)
	}

	go writeCounter(s) // Start Write thread
	go readCounter(s)  // Start Read thread
}

func writeCounter(s network.Stream) {
	// TODO write the file contents
	var counter uint64

	// infinite writing loop
	for {
		<-time.After(time.Second)
		counter++

		err := binary.Write(s, binary.BigEndian, counter)
		if err != nil {
			panic(err)
		}
	}
}

func readCounter(s network.Stream) {
	// TODO read the file contents

	// infinite reading loop
	for {
		var counter uint64

		err := binary.Read(s, binary.BigEndian, &counter)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Received %d from %s\n", counter, s.ID())
	}
}