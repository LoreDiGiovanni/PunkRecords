package main

import (
    "github.com/LoreDiGiovanni/punkrecords/p2p"
    "github.com/LoreDiGiovanni/punkrecords/storage"
    "log"
)

func main() {
    transportOpts := p2p.TransportOpts{
        ListenAddr: "127.0.0.1:8000", 
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
    err := server.Start()
    if err != nil {
        log.Fatal(err)
    }


     
    
    select {}

}
