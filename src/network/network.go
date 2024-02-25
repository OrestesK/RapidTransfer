package network

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const protocolID = "RapidTransfer" // this is just a unique id, can be whatever, keeps the heckers away

func Initialize_node() host.Host {
	// Create p2p node
	// Listen only on ( ipv4 and tcp )
	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}
	return node
}

func Server(node host.Host, done chan bool) string {
	node.SetStreamHandler(protocolID, func(s network.Stream) {
		go writeCounter(s, done)
	})

	key := fmt.Sprintf("%s/p2p/%s", node.Addrs()[1], node.ID())
	return key
}

func Client(node host.Host, peerAddr string) {
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

	go writeCounter(s, nil) // Start Write thread TODO FIX THIS
}

func writeCounter(s network.Stream, done chan bool) {
	// TODO write the file contents
	var counter uint64

	<-time.After(time.Second)
	counter++

	err := binary.Write(s, binary.BigEndian, counter)
	if err != nil {
		panic(err)
	}
	done <- true
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
