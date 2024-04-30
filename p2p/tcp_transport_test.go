package p2p

import (
  "testing"
  "github.com/stretchr/testify/assert"
)


func TestNewTCPTransport(t *testing.T) {
    listenAddres := "127.0.0.1:8000"
    tr := NewTCPTransport(listenAddres)
    assert.Equal(t,tr.listenAddres,listenAddres) 
    
    assert.Nil(t,tr.ListenAndAccept())

    select {}
}



