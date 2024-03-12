package network

import (
	"context"
	"encoding/binary"
	"fmt"

	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

func Server(node host.Host, file_name string, done chan bool) string {
	node.SetStreamHandler(protocolID, func(s network.Stream) {
		go write(s, file_name, done)
	})

	key := fmt.Sprintf("%s/p2p/%s", node.Addrs()[1], node.ID())
	return key
}

func Client(node host.Host, peerAddr string, file_name string, done chan bool, delete bool) {
	// Parse the multiaddr string.
	peerMA, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		panic(err)
	}
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	fmt.Println(peerAddrInfo)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Info: %s PeerAddr: %s filename %s\n", peerAddrInfo, peerAddr, file_name)

	// Connect to given address
	if err := node.Connect(context.Background(), *peerAddrInfo); err != nil {
		panic(err)
	}
	// fmt.Println("Connected to", peerAddrInfo.String())

	// Create a stream with peer
	s, err := node.NewStream(context.Background(), peerAddrInfo.ID, protocolID)
	if err != nil {
		panic(err)
	}

	//TODO THIS GETS HERE
	go read(s, file_name, done, delete) // Start Read thread
}

func write(s network.Stream, file_name string, done chan bool) {
	pwd, _ := os.Getwd()
	path := filepath.FromSlash(pwd + "/" + strings.TrimSuffix(file_name, "\r\n"))
	bytes, err := os.ReadFile(path)

	if err != nil {
		done <- true
		fmt.Println(err)
		return
	}

	err = binary.Write(s, binary.BigEndian, int64(len(bytes))) // Send the length of data
	if err != nil {
		done <- true
		fmt.Println(err)
		return
	}

	err = binary.Write(s, binary.BigEndian, bytes) // Send the actual data
	if err != nil {
		done <- true
		fmt.Println(err)
		return
	}

	done <- true
}

func read(s network.Stream, file_name string, done chan bool, delete bool) {
	var dataLength int64
	err := binary.Read(s, binary.BigEndian, &dataLength)
	if err != nil {
		done <- true
		fmt.Println(err)
		return
	}

	byteArray := make([]byte, dataLength)
	err = binary.Read(s, binary.BigEndian, &byteArray)
	if err != nil {
		done <- true
		fmt.Println(err)
		return
	}

	if delete {
		done <- true
		return
	}
	err = os.WriteFile(file_name, byteArray, fs.FileMode(0755))
	if err != nil {
		// Handle access denied by creating a new file with a modified name
		count := 1
		for {
			newFileName := fmt.Sprintf("%d%s", count, file_name)
			err = os.WriteFile(newFileName, byteArray, fs.FileMode(0755))
			if err == nil {
				fmt.Printf("File saved as %s\n", newFileName)
				break
			}
			count++
		}
	} else {
		done <- true
		//fmt.Println(err)
		return
	}
	fmt.Println("File saved as file.txt")

	done <- true
}
