package p2p

type Message struct{
    From string 
    Payload interface{} 
}

type StoreFile struct {
    Key string
    BufSize int64
}

type GetFile struct {
    Key string
}



