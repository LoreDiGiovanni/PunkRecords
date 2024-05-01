package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
    "strings"
)

type Storage interface {
    writestreem(key string , r io.Reader) error    
}


type PathTransformFunc func(string) (string, string)

var NoPathTransformFunc = func(filename string) (string, string) {
    return "./",filename 
}
var CASPathTransformFunc = func(filename string) (string,string) {
    hash := sha1.Sum([]byte(filename))
    hashStr := hex.EncodeToString(hash[:])
    
    blockSize := 5
    numBlocks := len(hashStr) / blockSize
    path := make([]string, numBlocks)

    for i := 0; i < numBlocks; i++ {
        start := i*blockSize
        end := start+blockSize
        path[i] = hashStr[start:end]
    }
    return strings.Join(path,"/"),hashStr
}

type StorageOpts struct {
    PathTransformFunc PathTransformFunc
}


