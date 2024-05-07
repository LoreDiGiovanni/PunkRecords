package main

import (
    "github.com/LoreDiGiovanni/punkrecords/p2p"
    "github.com/LoreDiGiovanni/punkrecords/storage"
    "log"
    "sync"
)

type server struct {
   Transport  p2p.Transport
   Storage    storage.Storage
   quitch     chan struct{}  
   KnowPeers  []string
   peerlock   sync.RWMutex
   ActivePeers map[string]p2p.Peer
}

func NewServer(transport p2p.Transport,storage storage.Storage) *server {
    return &server{
            Transport: transport,
            Storage: storage,
            quitch: make(chan struct{}),
            ActivePeers: make(map[string]p2p.Peer),
    }
}

func (s *server) Close(){
    close(s.quitch)
}

func (s *server) BootstrapKnownPeers() {
    for _, addr := range s.KnowPeers {
        if err := s.Transport.Dial(addr); err != nil {
            log.Printf("[TCP] Pear %s Offline\n", addr) 
        }else {
            log.Printf("[TCP] Pear %s Online\n", addr)
        }
    }
}

func (s *server) onPeer(peer p2p.Peer) error{
    s.peerlock.Lock()
    defer s.peerlock.Unlock()
    s.ActivePeers[peer.RemoteAddr()] = peer
    log.Printf("[TCP] OnPear %s \n", peer.RemoteAddr())
    return nil
    
}

func (s *server) Start() error {
    s.BootstrapKnownPeers() 
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
                log.Printf("%s: %s\n", msg.From, msg.Payload)
            case <-s.quitch:
                return nil 
        }  
    }
}
