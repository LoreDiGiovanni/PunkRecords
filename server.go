package main

import (
    "github.com/LoreDiGiovanni/punkrecords/p2p"
    "github.com/LoreDiGiovanni/punkrecords/storage"
    "log"
    "sync"
    "io"
    "encoding/gob"
    "bytes"
    "net"
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
    err := s.Storage.Writestreem(key, tee)
    if err != nil {
        if err.Error() == "ErrAlreadyExists" {
            return nil
        }
        return err
    }else {
        payload := p2p.StoreFile{Key: key, BufSize: int64(buf.Len())}
        return s.Broadcast(payload,buf)
    }

}

func (s *server) Broadcast(payload p2p.StoreFile, r io.Reader) error{
    buf := new(bytes.Buffer) ; 
    message := p2p.Message{
        From: s.Transport.GetAddr(),
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
            case msg := <-s.Transport.Consume():
                log.Printf("[MSG][%s] Received From [%s]\n",s.Transport.GetAddr(), msg.From)
                s.handleMessagePayload(&msg)
                //s.StoreData(p.Key, bytes.NewReader(p.Data))
            case <-s.quitch:
                return nil 
        }  
    }
}


func (s *server) handleMessagePayload(msg *p2p.Message) error {
    log.Printf("\t Received %s\n", msg.Payload)
    switch msg.Payload.(type) {
    case p2p.StoreFile:
        p := msg.Payload.(p2p.StoreFile)
        log.Printf("\t Received %s :)\n",p.Key,)
        return nil
    case p2p.GetFile:
        p := msg.Payload.(p2p.GetFile)
        log.Printf("\t Received %s :)\n",p.Key,)
        return nil
    }
    return nil
}

func (s *server) onPeer(peer p2p.Peer) error{
    s.ActivePeers[peer.RemoteAddr()] = peer
    log.Printf("[TCP] Exec onPear %s \n", peer.RemoteAddr())
    return nil
}
