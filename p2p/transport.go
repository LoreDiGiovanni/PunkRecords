package p2p


type Peer interface {
    Close() error
}

type Transport interface {
    // Method that have to listen, accept and handle connections
    ListenAndAccept() error
    // Method return a channel for consume messages     
    Consume() <-chan Message
}
type TransportOpts struct {
    ListenAddr string
    Handshake    HandshakeFunc
    Decoder      Decoder
    OnPeer      func(Peer) error
}
