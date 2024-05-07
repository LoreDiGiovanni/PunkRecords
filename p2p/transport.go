package p2p


type Peer interface {
    RemoteAddr() string
    Send([]byte) error
    Close() error
}

type Transport interface {
    // Method that have to listen, accept and handle connections
    ListenAndAccept() error
    // Method return a channel for consume messages     
    Consume() <-chan Message
    Dial(addr string) error
    Close() error
}
type TransportOpts struct {
    ListenAddr string
    Handshake    HandshakeFunc
    Decoder      Decoder
    OnPeer      func(Peer) error
}
