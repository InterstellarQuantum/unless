package text

import (
	"os"
)

//操作文本,主要目的是写入,把内容换行追加在文本里面

type Mt func(text string) error

func Preserve(path string) (Mt, error) {
	_, err := os.Stat(path)
	if err != nil {
		_, err := os.Create(path)
		if err != nil {
			return nil, err
		}
	}
	return func(text string) error {
		openFile, e := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0666)
		defer openFile.Close()
		if e != nil {
			return e
		}
		_, err := openFile.WriteString(text + "\n")
		if err != nil {
			return err
		}
		return nil
	}, nil
}
