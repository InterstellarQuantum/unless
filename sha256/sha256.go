package uses

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

//返回文件名和sha256
func UseSha256(path string) (map[string]string, error) {
	var result = make(map[string]string, 0)
	infos, e := ioutil.ReadDir(path)
	if e != nil {
		return result, e
	}
	wg := sync.WaitGroup{}
	wg.Add(len(infos))
	for _, v := range infos {
		go func(v fs.FileInfo) {
			if !(v.IsDir()) {
				pathname := filepath.Join(path, v.Name())
				sha, e := GetSHA256FromFile(pathname)
				if e != nil {
					result[v.Name()] = e.Error()
				} else {
					result[v.Name()] = sha
				}
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	return result, nil
}

func GetSHA256FromFile(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, nil
}
