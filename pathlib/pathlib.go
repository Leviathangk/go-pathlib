package pathlib

import "os"

type Parser struct {
	Path string      // 路径
	info os.FileInfo // 存储详情
}

func New(p string) *Parser {
	pathNew := new(Parser)
	pathNew.Path = p
	return pathNew
}
