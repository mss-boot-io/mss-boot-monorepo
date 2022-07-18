/*
 * @Author: lwnmengjing
 * @Date: 2021/12/16 7:39 下午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/12/16 7:39 下午
 */

package pkg

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// PathCreate create path
func PathCreate(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// PathExist path exist
func PathExist(addr string) bool {
	s, err := os.Stat(addr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// FileCreate create file
func FileCreate(content bytes.Buffer, name string, mode os.FileMode) error {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	changeStr := strings.ReplaceAll(content.String(), `\$`, `$`)
	changeStr = strings.ReplaceAll(changeStr, `\}`, "}")
	changeStr = strings.ReplaceAll(changeStr, `\{`, "{")
	_, err = file.WriteString(changeStr)
	if err != nil {
		log.Println(err)
	}
	return err
}

// FileCopy copy file
func FileCopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceFileStat.Mode())
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// GetSubPath get directory's subject path
func GetSubPath(directory string) ([]string, error) {
	dirs, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	subPath := make([]string, 0)
	for i := range dirs {
		if dirs[i].IsDir() {
			subPath = append(subPath, dirs[i].Name())
		}
	}
	return subPath, nil
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
