package p2p

import (
    "net"
    "sync"
)

type TCPTransport struct {
    listenAddres string  
    listener     net.Listener
    mu           sync.RWMutex
    peers        map[net.Addr]Peer
}
