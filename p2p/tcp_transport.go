package p2p

import (
	"net"
	"sync"
    "log"
    "io"
    "encoding/gob"

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
    msgch        chan Message
    wg           sync.WaitGroup 
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

func (t *TCPTransport) Dial(addr string) (net.Conn,error) {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return nil,err
    }
    go t.handleConn(conn)
    return conn, nil
}

func (t *TCPTransport) handleConn(conn net.Conn){
    gob.Register(Message{})
    gob.Register(StoreFile{})
    gob.Register(GetFile{})
    var err error
    defer func() {
        if err != nil {
            log.Printf("[TCP] %+v: Drop connection Error: %+v\n", conn.RemoteAddr(), err)
        }
        conn.Close()
    }()
    msg := Message{From: t.GetAddr(),}
    for {
        err := gob.NewDecoder(conn).Decode(&msg)
        if err != nil {
            if err == io.EOF {
                log.Printf("[TCP] Connection closed by %+v\n", conn.RemoteAddr())
                err = nil 
                break
            }else {
                log.Printf("[TCP] Decoding error: %s\n", err)
            }
        }else{
            log.Printf("[TCP] Decoded Payload: %s\n",msg.Payload)
            t.msgch <-msg
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

func (t *TCPTransport) GetAddr() string {
    return t.ListenAddr
}
