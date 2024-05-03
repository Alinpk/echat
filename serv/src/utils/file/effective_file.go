package fileop

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var MaxSize int64 = 500 * 1024 * 1024 // 500m
var MaxBufferSz = 1024 * 5            // 5k

type EFile struct {
	// this is not strictly
	MaxSize     int64 // file size
	MaxBufferSz int   // how many strings can buffer store

	buffer []byte
	Name   string // fileName
	Dir    string // where to persist

	fs *os.File
	ch chan ([]byte)
	mu sync.Mutex

	isClosed bool
}

func (ef *EFile) Close() {
	// if channel close is done and goroutine change to recv
	// and run persist function, some unexpected behaviour may happened
	ef.mu.Lock()
	defer ef.mu.Unlock()
	close(ef.ch)
	ef.isClosed = true
}

func OpenEFile(path string) (ef *EFile, err error) {
	dir := filepath.Dir(path)
	name := filepath.Base(path)

	var fs *os.File
	fs, err = OpenFile(filepath.Join(dir, name))
	if err != nil {
		return
	}
	ef = &EFile{
		MaxSize:     MaxSize,
		MaxBufferSz: MaxBufferSz,

		buffer: make([]byte, 0, MaxBufferSz),
		Name:   name,
		Dir:    dir,

		fs: fs,
		ch: make(chan ([]byte), 10), // maybe can base in some param
	}

	go ef.Recv()
	return
}

func (ef *EFile) Write(buf []byte) (n int, err error) {
	defer func() {
		// maybe close
		if r := recover(); r != nil {
			n = 0
			err = fmt.Errorf("%v", r)
		}
	}()
	ef.ch <- buf
	return len(buf), nil
}

func (ef *EFile) Recv() {
	for bytes := range ef.ch {
		ef.buffer = append(ef.buffer, bytes...)
		if len(ef.buffer) < ef.MaxBufferSz {
			continue
		}

		old := ef.buffer
		ef.buffer = make([]byte, 0, MaxBufferSz)
		go ef.Persist(old)
	}
	go ef.Persist(ef.buffer)
}

func (ef *EFile) Persist(buf []byte) {
	ef.mu.Lock()
	defer ef.mu.Unlock()
	// ! this func should be call before unlock
	defer func() {
		if ef.isClosed {
			ef.fs = nil
		}
	}()
	// for situation: recvNormal->close->persistNormal->recvClose->persistClose
	if ef.fs == nil {
		return
	}
	for len(buf) != 0 {
		n, err := ef.fs.Write(buf)
		if err != nil {
			panic(err)
		}
		buf = buf[n:]
	}

	fi, _ := ef.fs.Stat()
	if fi.Size() < ef.MaxSize {
		if ef.isClosed {
			ef.fs.Close()
		}
		return
	}

	// file is too large, rename file and open a new one
	ef.fs.Close()
	if ef.isClosed {
		return
	}

	name := ef.Name + "." + time.Now().Format("20060102150405")
	oldname := filepath.Join(ef.Dir, ef.Name)
	newname := filepath.Join(ef.Dir, name)
	os.Rename(oldname, newname)

	ef.fs, _ = OpenFile(oldname)
}
