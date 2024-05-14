package p2p

import (
    "net"
)

type Peer interface {
    net.Conn
    Send([]byte) error
}

type Transport interface {
    // Method that have to listen, accept and handle connections
    ListenAndAccept() error
    // Method return a channel for consume messages     
    Consume() <-chan net.Conn 
    Dial(addr string) (net.Conn,error)
    GetAddr() string
    //Close() error
}
type TransportOpts struct {
    ListenAddr string
    Decoder    Decoder
    OnPeer     func(Peer) error
}
