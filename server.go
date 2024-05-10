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

func (s *server) BootstrapKnownPeers() {
    for _, addr := range s.KnowPeers {
        err := s.Transport.Dial(addr); 
        if err != nil {
            log.Printf("[TCP] Pear %s Offline\n", addr) 
        }else {
            log.Printf("[TCP] Pear %s Online\n", addr)
        }
    }
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
        payload := payload{Key: key, Data: buf.Bytes()}
        return s.Broadcast(payload)
    }

}

func (s *server) Broadcast(payload payload) error{
    buf := new(bytes.Buffer) ; 
    gob.NewEncoder(buf).Encode(payload)
    for _, peer := range s.ActivePeers {
        log.Printf("[MSG][%s] Send to[%s]: [%s,%s]\n",s.Transport.GetAddr(), peer.RemoteAddr(), payload.Key,payload.Data)
        peer.Send(buf.Bytes())
    }
    
    return nil
}

func (s *server) Start() error {
    if err := s.Transport.ListenAndAccept(); err != nil {
        return err
    }else {
        s.BootstrapKnownPeers()
        go s.loop()
        return nil
    } 
}

func (s *server) loop() error {
    for {
        select {
            case msg := <-s.Transport.Consume():
                var p payload 
                err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&p)
                if err != nil {
                    log.Fatal("[MSG] Failed to decode message: ", err)
                }else {
                    log.Printf("[MSG][%s] Received From [%s] Message: [%s,%s]\n",s.Transport.GetAddr(), msg.From,p.Key,p.Data)
                    s.StoreData(p.Key, bytes.NewReader(p.Data))
                }
            case <-s.quitch:
                return nil 
        }  
    }
}

func (s *server) onPeer(peer p2p.Peer) error{
    s.ActivePeers[peer.RemoteAddr()] = peer
    log.Printf("[TCP] Exec onPear %s \n", peer.RemoteAddr())
    return nil
}
