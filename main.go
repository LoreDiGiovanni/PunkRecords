package main

import (
	"log"
    "flag"
	"github.com/LoreDiGiovanni/punkrecords/p2p"
	"github.com/LoreDiGiovanni/punkrecords/storage"
)

func main() {
    var port string
    flag.StringVar(&port, "p", ":8000", "La porta su cui ascoltare")
    flag.Parse()
    
    transportOpts := p2p.TransportOpts{
        ListenAddr: port, 
        Handshake: p2p.NoHandshakeFunc,
        Decoder: p2p.BytesDecoder{},
        OnPeer: nil,
    }
    transport := p2p.NewTCPTransport(transportOpts) 

    storageOpts := storage.StorageOpts{
        Root: "./db",
        PathTransform: storage.CASPathTransformFunc, 
    }
    storage := storage.NewDefaultStorage(storageOpts)
    server := NewServer(transport,storage)
    server.Transport.(*p2p.TCPTransport).OnPeer = server.onPeer
    server.KnowPeers = append(server.KnowPeers, ":8001")

    err := server.Start()

    if err != nil {
        log.Fatal(err)
    }
    select {}
}
