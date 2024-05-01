package storage

import (
    "bytes"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestDefaultStorage(t *testing.T) {
    opts := StorageOpts{
        PathTransformFunc: CASPathTransformFunc, 
    }
    storage := NewDefaultStorage(opts)
    assert.NotNil(t, storage)
    expextedPath := "8b5bf/c68fe/7e1e7/575c7/93bc4/d74ab/b8ff6/0ebb1"
    expectedFilename := "8b5bfc68fe7e1e7575c793bc4d74abb8ff60ebb1"
    data := bytes.NewReader([]byte("hello world"))
    if err := storage.writestreem("FileName", data); err != nil {
        t.Fatal(err)
    }else {
        path, filename := storage.PathTransformFunc("FileName")
        assert.Equal(t,path, expextedPath) 
        assert.Equal(t,filename, expectedFilename)
    }
}
