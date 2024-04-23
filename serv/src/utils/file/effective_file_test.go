package fileop

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"time"
	"path/filepath"
	"regexp"
	"strings"
	"os"
)

func TestOpenAndClose(t *testing.T) {
	// remove test directory if existed
	os.RemoveAll("./test")
	MaxSize = 5 * 1024 // change to 5KB
	MaxBufferSz = 5 * 1024

	path := "./test/tfile"
	f, e := OpenEFile(path)
	defer os.RemoveAll("./test")
	assert.Equal(t, e, nil)

	// we suppose msg will write into cache
	writeMsg := make([]byte, MaxBufferSz)
	f.Write(writeMsg)

	// that's check it
	f.fs.Close()

	fi, _ := os.Stat(path)
	assert.Equal(t, fi.Size(), int64(0))

	f.fs, _ = OpenFile(path)
	f.Close()
	// !TODO timeout panic
	for f.fs != nil { time.Sleep(time.Millisecond * 10) }
	fi, _ = os.Stat(path)
	assert.Equal(t, fi.Size(), MaxSize)
}

func TestAutoBackUp(t *testing.T) {
	MaxSize = 5 * 1024 // change to 5KB
	MaxBufferSz = 5 * 1024

	path := "./test/tfile"
	f, e := OpenEFile(path)
	defer os.RemoveAll("./test")
	assert.Equal(t, e, nil)

	// we suppose msg will write into cache
	writeMsg := make([]byte, MaxBufferSz + (MaxBufferSz) / 2)
	f.Write(writeMsg[:MaxBufferSz + 1])
	//!TODO need to fix bug
	time.Sleep(time.Millisecond * 10)
	f.Write(writeMsg[:MaxBufferSz / 2])

	f.Close()
	for f.fs != nil { time.Sleep(time.Millisecond * 10) }

	dir := filepath.Dir(path)
	files, err := ioutil.ReadDir(dir)
	assert.Equal(t, err, nil)

	// only 2 file
	for _, file := range files {
		assert.Equal(t, strings.HasPrefix(file.Name(), "tfile"), true)
		if file.Name() == "tfile" {
			assert.Equal(t, file.Size(), MaxSize / 2)
		} else {
			re := regexp.MustCompile(`tfile\.[0-9]{14}`)
			assert.Equal(t, re.MatchString(file.Name()), true)
			assert.Equal(t, file.Size(), MaxSize + 1)
		}
	}
}