package guard

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"serv/utils/log"
	"time"
)

type GuardCfg struct {
	TidyUpInterval time.Duration // 整理文件的周期
	GuardPath      string        // 守护的路径
	FlowAddr       string        // 任务队列的地址
	ProcessDir     string        // 处理路径
}

var gcfg GuardCfg

func StartGuard(guardCfg GuardCfg) {
	gcfg = guardCfg

	go startGuard()
}

func startGuard() {
	ticker := time.NewTicker(gcfg.TidyUpInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fileList, err := Traverse(gcfg.GuardPath)
			fmt.Println(fileList)
			if err != nil {
				log.L.Warn("traverse failed", "err", err.Error())
				continue
			}
			err = Backup(fileList)
			fmt.Println("err:", err)
		}
	}
}

func Traverse(path string) ([]string, error) {
	fmt.Println("Traverse")
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0)
	re := regexp.MustCompile(`.*\.[0-9]{14}`)
	for _, file := range dir {
		if re.MatchString(file.Name()) {
			list = append(list, file.Name())
		}
	}
	return list, nil
}

func Backup(filelist []string) error {
	for _, filename := range filelist {
		src := filepath.Join(gcfg.GuardPath, filename)
		dst := filepath.Join(gcfg.ProcessDir, filename)

		// move file
		if err := os.Rename(src, dst); err != nil {
			// if move failed, only record a log and continue to process other file
			log.L.Warn("move file failed", "err", err.Error())
			continue
		}

		CompressAndTar(gcfg.ProcessDir, filename)
	}
	return nil
}

func CompressAndTar(dir, filename string) error {
	textFile, err := os.Open(filepath.Join(dir, filename))
	if err != nil {
		return err
	}
	defer textFile.Close()

	// 最终压缩产物
	outputfile := filename[:len(filename)-15] + ".tar.gz"
	tarGzFile, err := os.OpenFile(filepath.Join(dir, outputfile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer tarGzFile.Close()

	// 创建 gzip writer
	gzipWriter := gzip.NewWriter(tarGzFile)
	defer gzipWriter.Close()

	// 创建 tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// 创建文件头部信息
	fileInfo, err := textFile.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return err
	}
	header.Name = filename + "gz"

	// 写入文件头部信息
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	// 压缩文本文件内容并写入 tar 文件
	if _, err := io.Copy(tarWriter, textFile); err != nil {
		return err
	}
	os.Remove(filepath.Join(dir, filename))
	return nil
}
