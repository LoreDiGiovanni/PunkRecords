package main

import (
	"bytes"
    "time"
	"github.com/LoreDiGiovanni/punkrecords/p2p"
	"github.com/LoreDiGiovanni/punkrecords/storage"
)

func main() {
    s1 := makeServer(":8000","./db2")
    s2 := makeServer(":9000","./db1")
    s1.Transport.(*p2p.TCPTransport).OnPeer = s1.onPeer

    go s2.Start()
    time.Sleep(2 * time.Second)
    go s1.Start()
    time.Sleep(2 * time.Second)
    s1.StoreData("test", bytes.NewReader([]byte("hello world ")))
    s1.KnowPeers = append(s1.KnowPeers, ":9000")


    select {}
}

func makeServer(port string,root string,knowPeers ...string) *server {
    transportOpts := p2p.TransportOpts{
        ListenAddr: port, 
        Decoder: p2p.BytesDecoder{},
        OnPeer: nil,
    }
    storageOpts := storage.StorageOpts{
        Root: root,
        PathTransform: storage.CASPathTransformFunc, 
    }
    stor := storage.NewDefaultStorage(storageOpts)
    transport := p2p.NewTCPTransport(transportOpts) 
    return NewServer(transport,stor)
}

