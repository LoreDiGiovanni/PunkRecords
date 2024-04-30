package p2p

import (
	"fmt"
	"net"
	"sync"
)
type TCPPeer struct{
    conn net.Conn
    outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
    return &TCPPeer{
        conn: conn,
        outbound: outbound,
    } 
}

type TCPTransportOpts struct {
    ListenAddr string
    Handshake    HandshakeFunc
    Decoder      Decoder
}

type TCPTransport struct {
    TCPTransportOpts
    listenAddres string  
    listener     net.Listener
    
    mu           sync.RWMutex
    peers        map[net.Addr]Peer
}




func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
    return &TCPTransport{
        TCPTransportOpts: opts,
    }
}

func (t *TCPTransport) ListenAndAccept() error {
    var err error
    t.listener, err = net.Listen("tcp", t.ListenAddr)
    if err != nil {
        return err
    }else{
        go t.acceptLoop()
        return nil
    }
}

func (t *TCPTransport) acceptLoop() {
    for {
        conn, err := t.listener.Accept()
        if err != nil {
            fmt.Printf("TCP accepting error: %s\n", err)
        }
        go t.handleConn(conn)
    }
}

func (t *TCPTransport) handleConn(conn net.Conn){
    peer := NewTCPPeer(conn, true)
    ok ,err := t.Handshake(conn)
    if err != nil {
        fmt.Printf("TCP handshake error: %s\n", err)
        conn.Close()
    }else if !ok{
        fmt.Printf("TCP handshake failed\n")
        conn.Close() 
    }else{
        msg := make([]byte, 1024) 
        fmt.Printf("TCP connection from %+v\n", peer)
        for {
            if err := t.Decoder.Decode(conn,msg); err != nil {
                fmt.Printf("TCP decoding error: %s\n", err)
            }else{
                fmt.Printf("TCP received: %s\n", msg)
            }
        }
    }
}
