package p2p

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
)


func TestNewTCPTransport(t *testing.T) {
    Opts := TransportOpts{
        ListenAddr: "127.0.0.1:8000", 
        Handshake: NoHandshakeFunc,
        Decoder: BytesDecoder{},
        OnPeer: func(peer Peer) error {
            return nil
        },
    }
    tr := NewTCPTransport(Opts)
    
    assert.Equal(t,tr.ListenAddr,"127.0.0.1:8000") 
    assert.Nil(t,tr.ListenAndAccept())
}



