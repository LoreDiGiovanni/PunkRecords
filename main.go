package main

import (
    "log"
    "github.com/LoreDiGiovanni/punkrecords/p2p"
    "fmt"
)

func main() {
    tcpOpts := p2p.TCPTransportOpts{
        ListenAddr: "127.0.0.1:8000", 
        Handshake: p2p.NoHandshakeFunc,
        Decoder: p2p.BytesDecoder{},
        OnPeer: nil,
    }
    tr := p2p.NewTCPTransport(tcpOpts) 

    go func() {
        for {
            msg := <-tr.Consume()
            fmt.Printf("%s: %s\n", msg.From, msg.Payload)
        }
    }()

    if err := tr.ListenAndAccept(); err != nil { 
        log.Fatal(err) 
    } 
    
    select {}

}
