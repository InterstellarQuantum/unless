package main

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

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("input directory,like ./app /usr/local/bin")
		return
	}
	path := args[1]
	infos, e := ioutil.ReadDir(path)
	if e != nil {
		fmt.Println(e)
		return
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
			pathname := path + v.Name()
			sha, e := GetSHA256FromFile(pathname)
			if e == nil {
				fmt.Printf("fileName: %s,sha256: %s \n", v.Name(), sha)
			} else {
				fmt.Println(e)
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
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
