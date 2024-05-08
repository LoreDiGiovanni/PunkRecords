package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/LoreDiGiovanni/punkrecords/p2p"
	"github.com/LoreDiGiovanni/punkrecords/storage"
)

func main() {
    
    transportOpts := p2p.TransportOpts{
        ListenAddr: ":9000", 
        Handshake: p2p.NoHandshakeFunc,
        Decoder: p2p.BytesDecoder{},
        OnPeer: nil,
    }
    transport := p2p.NewTCPTransport(transportOpts) 

    storageOpts := storage.StorageOpts{
        Root: "./db1",
        PathTransform: storage.CASPathTransformFunc, 
    }
    stor := storage.NewDefaultStorage(storageOpts)
    server1 := NewServer(transport,stor)
    server1.Transport.(*p2p.TCPTransport).OnPeer = server1.onPeer
    go func() {
        server1.Start()
    }()

    time.Sleep(time.Second*3)
    fmt.Printf("\n\n")
    
    stor.Root = "./db2"
    transport.ListenAddr = ":8000"
    server2 := NewServer(transport,stor)
    server2.KnowPeers = append(server1.KnowPeers, ":9000")
    server2.Transport.(*p2p.TCPTransport).OnPeer = server2.onPeer
    server2.Start()

    time.Sleep(time.Second*3)
    server2.StoreData("test", bytes.NewReader([]byte("hello world ")))

    select {}
}

