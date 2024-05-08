package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type StorageOpts struct {
    Root string
    PathTransform PathTransformFunc
}

type Storage interface {
    // key is the file name
    Writestreem(key string , r io.Reader) error
    Reedstreem(key string) ([]byte,error)
    Delete(key string ) error
    DeleteAll() error 
    Exists(key string) bool

}

// function that transforms a filename into a path
// returns the path and the filename
//you can write your own pathfunc but they are provided by default
type PathTransformFunc func(string) (string, string)

var  NoPathTransformFunc = func(filename string) (string, string) {
    return "./",filename 
}
//Content Addrasseble Storage
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


func MakePathToFile(root string,path string,filename string) string{
    return fmt.Sprintf("%s/%s/%s",root,path,filename)

}


