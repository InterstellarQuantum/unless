package uses

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
)

//返回文件名和sha256
func UseSha256(path string) (map[string]string, error) {
	var result = make(map[string]string, 0)
	infos, e := ioutil.ReadDir(path)
	if e != nil {
		return result, e
	}
	if !strings.HasSuffix(path, `/`) || !strings.HasSuffix(path, `\\`) {
		sysType := runtime.GOOS
		if sysType == `linux` {
			path = path + `/`
		}
		if sysType == `windows` {
			path = path + `\\`
		}
	}
	wg := sync.WaitGroup{}
	wg.Add(len(infos))
	for _, v := range infos {
		go func(v fs.FileInfo) {
			if !(v.IsDir()) {
				pathname := path + v.Name()
				sha, e := GetSHA256FromFile(pathname)
				if e == nil {
					//fmt.Printf("fileName: %s,sha256: %s \n", v.Name(), sha)
					result[v.Name()] = sha
				} else {
					result[v.Name()] = e.Error()
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
