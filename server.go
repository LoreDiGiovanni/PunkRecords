package main

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"net"
	"sync"

	"github.com/LoreDiGiovanni/punkrecords/p2p"
	"github.com/LoreDiGiovanni/punkrecords/storage"
)

type server struct {
   Transport  p2p.Transport
   Storage    storage.Storage
   quitch     chan struct{}  

   KnowPeers  []string
   peerlock   sync.RWMutex
   ActivePeers map[net.Addr]p2p.Peer
}

func NewServer(transport p2p.Transport,storage storage.Storage) *server {
    return &server{
            Transport: transport,
            Storage: storage,
            quitch: make(chan struct{}),
            ActivePeers: make(map[net.Addr]p2p.Peer),
    }
}

func (s *server) Close(){
    close(s.quitch)
}

type payload struct {
    Key string
    Data []byte
}


func (s *server) StoreData(key string, r io.Reader) error {
    buf := new(bytes.Buffer)
    tee := io.TeeReader(r, buf)
    err := s.Storage.Write(key, tee)
    if err != nil {
        if err.Error() == "ErrAlreadyExists" {
            return nil
        }
        return err
    }else {
        payload := p2p.MessageStoreFile{Key: key, Size: int64(buf.Len())}
        return s.Broadcast(payload,buf)
    }

}

func (s *server) Broadcast(payload p2p.MessageStoreFile, r io.Reader) error{
    buf := new(bytes.Buffer) ; 
    message := p2p.Message{
        Payload: payload,
    }
    gob.NewEncoder(buf).Encode(message)
    for _, addr := range s.KnowPeers {
        conn, err := s.Transport.Dial(addr); 
        if err != nil {
            log.Printf("[SERVER] %s Offline\n",addr)
        }else {
            n, err := conn.Write(buf.Bytes())
            if err != nil{
                log.Printf("[SERVER] %s Unrichable\n",addr)
            }else {
                log.Printf("[SERVER] %s Sent %d bytes\n",addr, n)
            }
            n64,err := io.CopyN(conn, r, payload.Size)
        }
    }
    return nil
}

func (s *server) Start() error {
    if err := s.Transport.ListenAndAccept(); err != nil {
        return err
    }else {
        go s.loop()
        return nil
    } 
}

func (s *server) loop() error {
    for {
        select {
            case conn := <-s.Transport.Consume():
                log.Printf("[SERVER][%s]  New connection [%s -> %s]",s.Transport.GetAddr(),conn.LocalAddr().String(),conn.RemoteAddr().String())
                go s.handleConn(conn)
            case <-s.quitch:
                return nil 
        }  
    }
}
func (s *server) handleConn(conn net.Conn) error {
    defer conn.Close()
    msg := p2p.Message{}
    err := gob.NewDecoder(conn).Decode(&msg)
    if err != nil {
        if err.Error() == "EOF" {
            log.Printf("[SERVER][%s] Connection end with %s\n",s.Transport.GetAddr(),conn.RemoteAddr())
            return nil
        }else{
            log.Printf("[SERVER][%s] Read error: %s\n",s.Transport.GetAddr() ,err)
            return err
        }
    }
    switch msg.Payload.(type) {
        case p2p.MessageStoreFile:
            p := msg.Payload.(p2p.MessageStoreFile)
            log.Printf("[SERVER][%s]\t Received %s size %d \n",s.Transport.GetAddr(),p.Key,p.Size,)
            return s.StoreStream(p.Key, conn)
        case p2p.MessageGetFile:
            p := msg.Payload.(p2p.MessageGetFile)
            log.Printf("[SERVER][%s]\t Received %s \n",s.Transport.GetAddr(),p.Key,)
            return nil
        default:
            log.Printf("[SERVER][%s]\t Undefined payload type",s.Transport.GetAddr())
            return nil
    }
}



func (s *server) onPeer(peer p2p.Peer) error{
    s.ActivePeers[peer.RemoteAddr()] = peer
    log.Printf("[TCP] Exec onPear %s \n", peer.RemoteAddr())
    return nil
}

func init() {
	gob.Register(p2p.MessageStoreFile{})
	gob.Register(p2p.MessageGetFile{})
}
