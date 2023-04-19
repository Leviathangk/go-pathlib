package pathlib

import (
	"errors"
	"io"
	"os"
)

// WalkFunc 深度遍历的函数
type WalkFunc func(p *Parser, err error) error

// WalkDir 深度遍历文件夹
func (p *Parser) WalkDir(f WalkFunc) error {
	if !p.IsDir() {
		return NotDirErr(p.Path)
	}

	return p.walkDir(f)
}

// walkDir 遍历程序
func (p *Parser) walkDir(f WalkFunc) error {
	var file *os.File
	var names []string
	var err error

	// 打开文件夹
	file, err = os.Open(p.Path)
	if err != nil {
		return err
	}

	for {
		// 每次读一个文件
		names, err = file.Readdirnames(1)

		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				break
			} else {
				return err
			}
		}

		// 构造新 Path
		newPath := p.Join(names[0])

		// 读完立即处理
		err = f(newPath, err)
		if err != nil {
			return err
		}

		// 判断是否是文件夹继续处理
		if newPath.IsDir() {
			err = newPath.walkDir(f)
			if err != nil {
				if errors.Is(err, io.EOF) {
					err = nil
					break
				} else {
					return err
				}
			}
		}
	}

	return err
}
