package log

import (
	"golang.org/x/exp/slog"
	"os"
	"encoding/json"
	"io/ioutil"
	"serv/utils/file"
)

var cfg LogCfg
var L *slog.Logger
var file *os.File

type LogCfg struct {
	MaxSize int `json:"max_size"`
	PersistSz int `json:"persist_sz"`
	Path string `json:"path"`
	PersistLv string `json:"persist_level"`
}

func LoadCfg(path string) {
	fs, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(fs, &cfg)
	if err != nil {
		panic(err.Error())
	}
}

func InitLog() {
	fs, err := fileop.OpenFile(cfg.Path)
	if err != nil { panic(err.Error()) }
	var lv slog.Level
	switch cfg.PersistLv {
	case "LevelDebug":
		lv = slog.LevelDebug
	case "LevelInfo":
		lv = slog.LevelInfo
	case "LevelWarn":
		lv = slog.LevelWarn
	case "LevelError":
		lv = slog.LevelError
	}

	file = fs
    L = slog.New(slog.NewTextHandler(fs, &slog.HandlerOptions{Level: lv}))
}

// 后续用来存聊天记录
// func (h *logHandler) Write(buf []byte) (n int, err error) {
// 	if len(h.buffer) + len(buf) >= h.capacity {
// 		_, err = h.fs.Write(h.buffer)
// 		if err != nil { return }
// 		n, err = h.fs.Write(buf)
		
// 		h.buffer = h.buffer[:0]
// 		return
// 	}

// 	h.buffer = append(h.buffer, buf...)
// 	return len(buf), nil
// }