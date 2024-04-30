package main

import (
    "log"
    "github.com/LoreDiGiovanni/punkrecords/p2p"
)

func main() {
    tcpOpts := p2p.TCPTransportOpts{
        ListenAddr: "127.0.0.1:8000", 
        Handshake: p2p.NoHandshakeFunc,
        Decoder: p2p.NoDecoder{},
    }
    tr := p2p.NewTCPTransport(tcpOpts) 
    if err := tr.ListenAndAccept(); err != nil { 
        log.Fatal(err) 
    } 
    
    select {}

}
