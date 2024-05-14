package p2p

import (
	"encoding/gob"
	"io"
	"log"
	"net"
	"sync"
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
    connch        chan net.Conn
    wg           sync.WaitGroup 
    peers        map[net.Addr]Peer
}




func NewTCPTransport(opts TransportOpts) *TCPTransport {
    return &TCPTransport{
        TransportOpts: opts,
        connch: make(chan net.Conn),
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
        t.connch<- conn
    }
}

func (t *TCPTransport) Dial(addr string) (net.Conn,error) {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return nil,err
    }
    t.connch<- conn
    return conn, nil
}



func (t *TCPTransport) Decode(r io.Reader, msg *Message) error {
    return t.Decoder.Decode(r, msg)
}

func (t *TCPTransport) Consume() <-chan net.Conn{
    return t.connch
}

func (t *TCPTransport) Close() error {
    return t.listener.Close()
}

func (t *TCPTransport) GetAddr() string {
    return t.ListenAddr
}
func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}


