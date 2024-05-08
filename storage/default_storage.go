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

func (s *DefaultStorage) Writestreem(key string, r io.Reader) error {
    if s.Exists(key) {
       log.Printf("[DefaultStorage] Already exists in %s\n",s.Root)
       return nil  
    }
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
                log.Printf("[DefaultStorage] written %d bytes \n", n)
                return nil
            }
        }
    }
}   

func (s *DefaultStorage) Reedstreem(key string) ([]byte, error) { 
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
func (s *DefaultStorage) Exists(key string) bool {
    path,name := s.PathTransform(key)
    pathTofile := MakePathToFile(s.Root,path,name)
    _, err := os.Stat(pathTofile)
    return err == nil
}

func (s *DefaultStorage) Delete(key string) error {
    path,name := s.PathTransform(key)
    pathTofile := MakePathToFile(s.Root,path,name)
    return os.Remove(pathTofile)
    
}

func (s *DefaultStorage) DeleteAll() error {
   return os.RemoveAll(s.Root)
}



