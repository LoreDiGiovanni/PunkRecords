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



func NewDefaultStorage(options StorageOpts) *DefaultStorage {
    return &DefaultStorage{
        StorageOpts: options,
    }
}

func (s *DefaultStorage) writestreem(key string, r io.Reader) error {
    path,name := s.PathTransform(key)
    if err := os.MkdirAll(s.Root+"/"+path, os.ModePerm); err != nil {
        return err
    }else {
        buf := new(bytes.Buffer)
        io.Copy(buf, r)
        f, err := os.Create(MakePathToFile(s.Root,path,name))
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

func (s *DefaultStorage) reedstreem(key string) ([]byte, error) { 
    path,name := s.PathTransform(key)
    pathTofile := MakePathToFile(s.Root,path,name)
    f, err := os.Open(pathTofile)
    if err != nil {
        return nil, err
    }else {
        b, err := io.ReadAll(io.Reader(f))
        if err != nil {
            return nil, err
        }else {
            return b, nil
        }
    }
    
}
func (s *DefaultStorage) exists(key string) bool {
    path,name := s.PathTransform(key)
    pathTofile := MakePathToFile(s.Root,path,name)
    _, err := os.Stat(pathTofile)
    return err == nil
}

func (s *DefaultStorage) delete(key string) error {
    path,name := s.PathTransform(key)
    pathTofile := MakePathToFile(s.Root,path,name)
    return os.Remove(pathTofile)
    
}

func (s *DefaultStorage) deleteAll() error {
   return os.RemoveAll(s.Root)
}



