package p2p

import (
    "fmt"
	"net"
	"sync"
    "log"
    "io"

)
type TCPPeer struct{
    net.Conn
}

func NewTCPPeer(conn net.Conn) *TCPPeer {
    return &TCPPeer{
        Conn: conn,
    } 
}



func (p *TCPPeer) Send(b []byte) error {
    _, err := p.Write(b)
    return err
}


type TCPTransport struct {
    TransportOpts
    listenAddres string  
    listener     net.Listener
    msgch       chan Message
    
    mu           sync.RWMutex
    peers        map[net.Addr]Peer
}




func NewTCPTransport(opts TransportOpts) *TCPTransport {
    return &TCPTransport{
        TransportOpts: opts,
        msgch: make(chan Message),
    }
}

func (t *TCPTransport) ListenAndAccept() error {
    log.Printf("[TCP] Transport listening on %s\n", t.ListenAddr)
    var err error
    t.listener, err = net.Listen("tcp", t.ListenAddr)
    if err != nil {
        log.Printf("[TCP] Listen error: %s\n", err)
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
            log.Printf("[TCP] Accepting error: %s\n", err)
            break
        }else {
            log.Printf("[TCP] Accepted connection from %+v\n", conn.RemoteAddr())
        }
        go t.handleConn(conn)
    }
}

func (t *TCPTransport) Dial(addr string) error {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return err
    }
    go t.handleConn(conn)
    return nil
}

func (t *TCPTransport) handleConn(conn net.Conn){
    peer := NewTCPPeer(conn)
    var err error
    defer func() {
        if err != nil {
            log.Printf("[TCP] %+v: Drop connection Error: %+v\n", peer.RemoteAddr(), err)
        }
        peer.Conn.Close()
    }()
    if  t.OnPeer != nil{
        err = t.OnPeer(peer)
    }
    if err == nil {
        ok ,err := t.Handshake(conn)
        if err != nil {
            err = fmt.Errorf("[TCP] Handshake error: %s", err)
        }else if !ok{
            err = fmt.Errorf("[TCP] Handshake denied")
        }else{
            msg := Message{From: conn.RemoteAddr().String()}
            log.Printf("[TCP] Connection from %s\n", peer.RemoteAddr().String())
            for {
                // When the connection is closed, we get an EOF error
                // TODO: make a better error handling 
                if err = t.Decoder.Decode(conn,&msg); err != nil {
                    if err == io.EOF {
                        log.Printf("[TCP] Connection closed by %+v\n", peer)
                        err = nil 
                        break
                    }else {
                        log.Printf("[TCP] Decoding error: %s\n", err)
                    }
                }else{
                    t.msgch <-msg
                }
            }
        }
    } 
}
func (t *TCPTransport) Decode(r io.Reader, msg *Message) error {
    return t.Decoder.Decode(r, msg)
}

func (t *TCPTransport) Consume() <-chan Message {
    return t.msgch
}

func (t *TCPTransport) Close() error {
    return t.listener.Close()
}

