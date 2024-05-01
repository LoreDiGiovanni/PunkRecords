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

func (p *TCPPeer) Close() error {
    return p.conn.Close()
}

type TCPTransportOpts struct {
    ListenAddr string
    Handshake    HandshakeFunc
    Decoder      Decoder
    OnPeer      func(Peer) error
}

type TCPTransport struct {
    TCPTransportOpts
    listenAddres string  
    listener     net.Listener
    msgch       chan Message
    
    mu           sync.RWMutex
    peers        map[net.Addr]Peer
}




func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
    return &TCPTransport{
        TCPTransportOpts: opts,
        msgch: make(chan Message),
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

func (t *TCPTransport) Consume() <-chan Message {
    return t.msgch
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
    var err error
    defer func() {
        fmt.Printf("%+v: Drop connection Error: %+v\n", peer.conn.RemoteAddr(), err)
        peer.Close()
    }()
    if  t.OnPeer != nil{
        err = t.OnPeer(peer)
    }
    if err == nil {
        ok ,err := t.Handshake(conn)
        if err != nil {
            err = fmt.Errorf("TCP handshake error: %s", err)
        }else if !ok{
            err = fmt.Errorf("TCP handshake denied")
        }else{
            msg := Message{From: conn.RemoteAddr()}
            fmt.Printf("TCP connection from %+v\n", peer)
            for {
                if err := t.Decoder.Decode(conn,&msg); err != nil {
                    fmt.Printf("TCP decoding error: %s\n", err)
                }else{
                    t.msgch <-msg
                }
            }
        }
    } 

    
}
