package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/LoreDiGiovanni/punkrecords/p2p"
    "github.com/LoreDiGiovanni/punkrecords/storage"
)

func TestServer(t *testing.T) {
    transportOpts := p2p.TransportOpts{
        ListenAddr: "127.0.0.1:8000", 
        Handshake: p2p.NoHandshakeFunc,
        Decoder: p2p.BytesDecoder{},
        OnPeer: nil,
    }
    transport := p2p.NewTCPTransport(transportOpts)
    storageOpts := storage.StorageOpts{
        Root: "./db",
        PathTransform: storage.CASPathTransformFunc, 
    }
    storage := storage.NewDefaultStorage(storageOpts)
    server := NewServer(transport,storage)
    assert.NotNil(t,server)

    err := server.Start()
    assert.Nil(t, err)

    server.Close()
}

