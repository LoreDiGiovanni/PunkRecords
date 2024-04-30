package p2p


type Peer interface {}

type Transport interface {
    // Method that have to listen, accept and handle connections
    ListenAndAccept() error
}
