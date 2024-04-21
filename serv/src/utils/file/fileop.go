package fileop

import (
	"os"
	"path/filepath"
)

func OpenFile(filePath string) (file *os.File, err error) {
	file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	// file existed
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			// try create path but failed
			return
		}
		file, err = os.Create(filePath)
		return
	}
	return
}