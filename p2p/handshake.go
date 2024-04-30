package p2p

import "net"

//The HandshakeFunc function is a prototype for a function that, 
//based on parameters defined by the function's author,
//decides whether to add the peer to the list of peers to consider.
//It can evaluate factors such as latency, connection quality, etc.
type HandshakeFunc func(net.Conn) (bool,error)

func NoHandshakeFunc(net.Conn) (bool,error) {return true,nil}

