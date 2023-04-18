package pathlib

import "fmt"

// HasExistsErr 文件已存在错误
type HasExistsErr string

func (f HasExistsErr) Error() string {
	return fmt.Sprintf("文件已存在：%s", string(f))
}

// NotExistsErr 文件不存在错误
type NotExistsErr string

func (f NotExistsErr) Error() string {
	return fmt.Sprintf("文件不存在：%s", string(f))
}

// NotDirErr 非文件夹错误
type NotDirErr string

func (f NotDirErr) Error() string {
	return fmt.Sprintf("该路径不是文件夹：%s", string(f))
}
