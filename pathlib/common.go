package pathlib

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// Exists 判断路径是否存在
func (p *Parser) Exists() bool {
	_, err := p.Info()

	return err == nil
}

// Info 获取详情
func (p *Parser) Info() (os.FileInfo, error) {
	if p.info != nil {
		return p.info, nil
	}

	info, err := os.Stat(p.Path)
	if err != nil {
		return nil, err
	}
	p.info = info

	return p.info, nil
}

// IsFile 是否是文件
func (p *Parser) IsFile() bool {
	info, err := p.Info()
	if err != nil {
		return false
	}

	return !info.IsDir()
}

// IsDir 是否是文件夹
func (p *Parser) IsDir() bool {
	info, err := p.Info()
	if err != nil {
		return false
	}

	return info.IsDir()
}

// Join 合并路径
func (p *Parser) Join(paths ...string) *Parser {
	newPath := New(p.Path)

	for _, p := range paths {
		newPath.Path = filepath.Join(newPath.Path, p)
	}

	return newPath
}

// Parent 获取父路径
func (p *Parser) Parent() *Parser {
	return New(filepath.Dir(p.Path))
}

// Rename 重命名：真的只针对名字，输入新名字即可（含后缀）
// name 新的名字
// override 是否存在即覆盖，为 false 时，重复将会报 err
func (p *Parser) Rename(name string, override bool) (err error) {
	if p.Exists() {
		newPath := p.Parent().Join(name)
		if !override && newPath.Exists() {
			return HasExistsErr(newPath.Path)
		} else {
			return os.Rename(p.Path, newPath.Path)
		}
	}
	return NotExistsErr(p.Path)
}

// MkdirAll 创建路径：包含父路径，一般 0777
func (p *Parser) MkdirAll(mode os.FileMode) error {
	return os.MkdirAll(p.Path, mode)
}

// MoveTo 移动：包含路径及名字（名字不一样会被重命名）
// toPath：全路径，含有名字
// override 是否存在即覆盖，为 false 时，重复将会报 err
// 注意：如果是文件夹，那么源文件夹将会消失
func (p *Parser) MoveTo(toPath string, override bool) error {
	newPath := New(toPath)

	if newPath.Exists() {
		if !override {
			return HasExistsErr(toPath)
		}
	} else {
		err := newPath.Parent().MkdirAll(0777)
		if err != nil {
			return err
		}
	}

	return os.Rename(p.Path, toPath)
}

// Delete 删除：是文件夹则会整个文件夹及内部文件都被删除
func (p *Parser) Delete() error {
	if p.IsFile() {
		return os.Remove(p.Path)
	} else {
		return os.RemoveAll(p.Path)
	}
}

// Name 获取名字：是文件的话含后缀
func (p *Parser) Name() string {
	_, name := filepath.Split(p.Path)
	return name
}

// CleanName 获取名字：不含后缀
func (p *Parser) CleanName() string {
	name := p.Name()

	if p.IsDir() {
		return name
	}

	return name[0 : len(name)-len(p.Suffix())]
}

// Suffix 获取后缀：如 .txt
func (p *Parser) Suffix() string {
	return filepath.Ext(p.Path)
}

// ListDir 返回文件夹列表：如果文件夹很大，建议 Walk
func (p *Parser) ListDir() (allPaths []*Parser, err error) {
	var dir *os.File
	var names []string

	if !p.Exists() {
		return nil, NotExistsErr(p.Path)
	} else if !p.IsDir() {
		return nil, NotDirErr(p.Path)
	}

	dir, err = os.Open(p.Path)
	defer dir.Close()
	if err != nil {
		return
	}

	names, err = dir.Readdirnames(0) // <=0 返回所有
	if err != nil {
		return
	}

	for _, f := range names {
		allPaths = append(allPaths, New(filepath.Join(p.Path, f)))
	}

	return
}

// FindFunc 配合 FindFiles 使用
type FindFunc func(p *Parser, err error) error

// FindFiles 查找指定文件，返回路径处理器
func (p *Parser) FindFiles(pattern string, f FindFunc) error {
	var re *regexp.Regexp
	var err error

	re, err = regexp.Compile(pattern)
	if err != nil {
		return err
	}

	// 遍历查找
	err = p.WalkDir(func(p *Parser, err error) error {
		if re.MatchString(p.Path) {
			err = f(New(p.Path), err)
			if err != nil {
				if errors.Is(err, io.EOF) {
					err = nil
				} else {
					return err
				}
			}
		}
		return err
	})

	return err
}
