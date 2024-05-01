package storage

import (
    "io"
    "os"
    "log"
    "bytes"
)

type DefaultStorage struct {
    StorageOpts 
}

func NewDefaultStorage(opts StorageOpts) *DefaultStorage {
    return &DefaultStorage{
        StorageOpts: opts,
    }
}

func (s *DefaultStorage) writestreem(key string , r io.Reader) error {
    path,name := s.PathTransformFunc(key)
    if err := os.MkdirAll(path, os.ModePerm); err != nil {
        return err
    }else {
        buf := new(bytes.Buffer)
        io.Copy(buf, r)
        f, err := os.Create(path+"/"+name)
        if err != nil {
            return err
        }else {
            n, err := io.Copy(f, buf)
            if err != nil {
                return err
            }else {
                log.Printf("written %d bytes => %s\n", n, path)
                return nil
            }
        }
    }
}   


