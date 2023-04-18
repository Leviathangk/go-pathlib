# pathlib
一个类似 python 标准库的路径处理模块

# 使用
方法请查看源码

```
package main

import (
	"fmt"
	"log"

	"github.com/Leviathangk/go-pathlib/pathlib"
)

func main() {
	parser := pathlib.New("D:\\Go\\goTest3")
	// fmt.Println(parser.Rename("xxx.txt",true))
	// fmt.Println(parser.MoveTo("D:\\Go\\goTest3\\static", true))
	fmt.Println(parser.Name())
	fmt.Println(parser.CleanName())
	fmt.Println(parser.Suffix())
	fmt.Println(parser.ListDir())

	// 查找文件
	err := parser.FindFiles(".+\\.js", func(p *pathlib.Parser, err error) error {
		fmt.Println(p.Path)
		return err
	})

	if err != nil {
		log.Fatalln(err)
	}

	// 深度遍历文件夹
	err = parser.WalkDir(func(p *pathlib.Parser, err error) error {
		fmt.Println(p.Path)
		return err
	})

	if err != nil {
		log.Fatalln(err)
	}
}

```
