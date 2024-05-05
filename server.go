package main

import (
    "github.com/LoreDiGiovanni/punkrecords/p2p"
    "github.com/LoreDiGiovanni/punkrecords/storage"
)

type server struct {
   Transport p2p.Transport
   Storage storage.Storage
}

func NewServer(transport p2p.Transport,storage storage.Storage) *server {
    return &server{
            Transport: transport,
            Storage: storage,
    }
}

func (s *server) Start() error {
    if err := s.Transport.ListenAndAccept(); err != nil {
        return err
    }else {
        go s.Transport.acceptLoop()
        return nil
    }
}

