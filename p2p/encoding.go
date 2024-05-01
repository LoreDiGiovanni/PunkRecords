package p2p

import (
    "encoding/gob"
    "io"
)

type Decoder interface {
    Decode(io.Reader, *Message) error
}

type GOBDecoder struct{}
func (d GOBDecoder) Decode(r io.Reader, msg *Message) error {
    return gob.NewDecoder(r).Decode(&msg.Payload)
}

//Default decoder
type BytesDecoder struct{}
func (d BytesDecoder) Decode(r io.Reader, msg *Message) error {
    buffer := make([]byte, 1024)
    n, err := r.Read(buffer)
    if err != nil {
        return err
    }else{
        msg.Payload = buffer[:n]
        return nil
    }
}


