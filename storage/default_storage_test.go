package storage

import (
    "bytes"
    "testing"
    "github.com/stretchr/testify/assert"
)




func TestDefaultStorage(t *testing.T) {
    strOpts := StorageOpts{
        root: "./db",
        pathTransform: CASPathTransformFunc,
    }
    storage := NewDefaultStorage(strOpts)
    assert.NotNil(t, storage)
    expextedPath := "8b5bf/c68fe/7e1e7/575c7/93bc4/d74ab/b8ff6/0ebb1"
    expectedFilename := "8b5bfc68fe7e1e7575c793bc4d74abb8ff60ebb1"
    bytesstring := []byte("hello world")
    data := bytes.NewReader(bytesstring)
    if err := storage.writestreem("FileName", data); err != nil {
        t.Fatal(err)
    }else {
        path, filename := CASPathTransformFunc("FileName")
        assert.Equal(t,path, expextedPath) 
        assert.Equal(t,filename, expectedFilename)
        b, err := storage.reedstreem("FileName")
        if err != nil {
            t.Fatal(err)
        }else {
            if !assert.Equal(t, b, bytesstring){
                t.Fatal("not equal\n expected: ", bytesstring, "\n got: ", b) 
            }
        }
        if err := storage.delete("FileName"); err != nil {
            t.Fatal(err)
        }else {
            
        }
    }
}


    



