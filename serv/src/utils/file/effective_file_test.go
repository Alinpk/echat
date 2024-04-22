package fileop

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"fmt"
	"os"
)

func TestOpenFile(t *testing.T) {
	MaxSize = 5 * 1024 // change to 5KB
	f, e := OpenEFile("./test/tfile")
	assert.Equal(t, e, nil)

	writeMsg := make([]byte, MaxBufferSz)
	f.Write(writeMsg)

	time.Sleep(time.Second * 1)
	fi, _ := f.fs.Stat()
	fmt.Println("fi size:", fi.Size())

	f.Write(writeMsg)

	time.Sleep(time.Second * 1)
	fi, _ = f.fs.Stat()
	fmt.Println("fi size:", fi.Size())

	fi, _ = os.Stat("./test/tfile_20240422233532")
	fmt.Println("fi size:", fi.Size())

}