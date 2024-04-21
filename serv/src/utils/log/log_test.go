package log

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"regexp"
)

func TestLog(test *testing.T) {
	// cmdLineOutPutHandler(os.Stderr)
	// path := "./test.log"
	// defer os.Remove(path)
	// f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	// assert.Equal(test, err, nil)
	// cmdLineOutPutHandler(f)
	// f.Close()
	// f, err = os.OpenFile(path, os.O_RDONLY, 0600)
	// assert.Equal(test, err, nil)
	// buf := make([]byte, 175)
	// n1, err := f.Read(buf)
	// fmt.Println(n1, ":", string(buf))

	// var logFormat = regexp.MustCompile(`time=.* level=.* msg=.*`)
	// fmt.Println(logFormat.MatchString(string(buf)))
	// f.Close()
	{
		LoadCfg("log_cfg.json")
		err := os.MkdirAll("test", os.ModePerm)
		assert.Equal(test, err, nil)
		defer os.RemoveAll("test")
	}
	{
		InitLog()
		L.Warn("this is a warn log")
		file.Close()
	}
	{
		fs, err := os.OpenFile(cfg.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		assert.Equal(test, err, nil)
		f, _ := fs.Stat(); size := f.Size()
		buf := make([]byte, size)
		fs.Read(buf)
		var logFormat = regexp.MustCompile(`time=.* level=.* msg=.*`)
		assert.Equal(test, logFormat.MatchString(string(buf)), true)
		fs.Close()
	}
}