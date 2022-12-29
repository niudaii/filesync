package filesync

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var cwd string

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		log.Println(err)
	}
}

// Traverse walk要同步的文件, 生成md5, 并返回列表
func Traverse(path string) ([]string, error) {
	f, fErr := os.Lstat(path)
	if fErr != nil {
		return nil, fErr
	}
	var dir string
	var base string
	if f.IsDir() {
		dir = path
		base = "."
	} else {
		dir = filepath.Dir(path)
		base = filepath.Base(path)
	}
	// 相对目录遍历，需要cd到指定目录下；遍历完后必须返回原来的当前位置
	fErr = os.Chdir(dir)
	if fErr != nil {
		return nil, fErr
	}
	defer os.Chdir(cwd)

	md5List := make([]string, 10)
	var md5Str string
	WalkFunc := func(path string, info os.FileInfo, err error) error {
		//文件路径分隔符统一处理为 /
		if runtime.GOOS == "windows" {
			path = strings.ReplaceAll(path, "\\", "/")
		}
		if checkFileIsSyncList(path) == false {
			return nil
		}
		if info.IsDir() {
			md5Str = path + ",,Directory"
			md5List = append(md5List, md5Str)
		}
		if info.Mode().IsRegular() {
			md5Str, fErr = Md5OfAFile(path)
			if fErr != nil {
				return fErr
			}
			md5Str = path + ",," + md5Str
			md5List = append(md5List, md5Str)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			filename, err := os.Readlink(path)
			if err != nil {
				return err
			}
			md5Str = "symbolLink&&" + filename
			md5Str = path + ",," + md5Str
			md5List = append(md5List, md5Str)
		}
		return nil
	}
	fErr = filepath.Walk(base, WalkFunc)
	if fErr != nil {
		return nil, fErr
	}
	for i, v := range md5List {
		if v == "" || v == ".,,Directory" {
			continue
		}
		md5List = md5List[i:]
		break
	}
	if md5List[len(md5List)-1] == ".,,Directory" {
		return []string{}, nil
	}
	return md5List, nil
}

func Md5OfAFile(f string) (string, error) {
	fi, fiErr := os.Open(f)
	if fiErr != nil {
		return "", fiErr
	}
	defer fi.Close()
	r := bufio.NewReader(fi)
	h := md5.New()
	var s string
	var e error
	for {
		s, e = r.ReadString('\n')
		io.WriteString(h, s)
		if e != nil {
			if e == io.EOF {
				break
			} else {
				return "", e
			}
		}
	}
	s = fmt.Sprintf("%x", h.Sum(nil))
	return s, nil
}
