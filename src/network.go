package main

import (
	"context"
	// "encoding/binary"
	// "flag"
	"fmt"
	// "os"
	// "os/signal"
	// "syscall"
	// "time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func initialize_node() host.Host {
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

func server(node host.Host) {
	node.SetStreamHandler(protocolID, func(s network.Stream) {
		go writeCounter(s)
		go readCounter(s)
	})
	key := fmt.Sprintf("%s/p2p/%s", node.Addrs()[0], node.ID()) // Format key (it's the address + "/p2p/" + id)
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

// func main() {
// 	node := initialize_node() // treat this as the 'connection parameters'
// 	defer node.Close()

// 	// parse given address (TODO this will be replaced with database)
// 	peerAddr := flag.String("p", "", "peer address")
// 	flag.Parse()

// 	if *peerAddr == "" {
// 		server(node)
// 	} else {
// 		// if peer address was provided, connect to it
// 		client(node, *peerAddr)
// 	}

// 	sigCh := make(chan os.Signal)                         // create channel
// 	signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT) // notify the channel when SIGKILL or SIGINT received
// 	<-sigCh                                               // block until a signal is received
// 	// basically wait here until I do Ctrl+C
// }

// func writeCounter(s network.Stream) {
// 	// TODO write the file contents
// 	var counter uint64

// 	// infinite writing loop
// 	for {
// 		<-time.After(time.Second)
// 		counter++

// 		err := binary.Write(s, binary.BigEndian, counter)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

// func readCounter(s network.Stream) {
// 	// TODO read the file contents

// 	// infinite reading loop
// 	for {
// 		var counter uint64

// 		err := binary.Read(s, binary.BigEndian, &counter)
// 		if err != nil {
// 			panic(err)
// 		}

// 		fmt.Printf("Received %d from %s\n", counter, s.ID())
// 	}
// }
