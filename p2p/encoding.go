package p2p

import (
    "encoding/gob"
    "io"
)

type Decoder interface {
    Decode(io.Reader,any) error
}

type GOBDecoder struct{}
func (d GOBDecoder) Decode(r io.Reader, v any) error {
    return gob.NewDecoder(r).Decode(v)
}

type NoDecoder struct{}
func (d NoDecoder) Decode(r io.Reader, v any) error {
    buffer := make([]byte, 1024)
    n, err := r.Read(buffer)
    if err != nil {
        return err
    }else{
        copy(v.([]byte), buffer[:n])
        return nil
    }
}


