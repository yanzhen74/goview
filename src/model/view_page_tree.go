package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type View_page_tree struct {
	Id       string
	Curname  string
	Curdir   string
	Branches []View_page_tree
	Isleaf   bool
}

var pid = 0

func Init_pages(dir string) (data View_page_tree, err error) {
	dir, _ = filepath.Abs(dir)
	//判断文件或目录是否存在
	file, err := os.Stat(dir)
	if err != nil {
		return data, err
	}
	data = View_page_tree{}

	//如果不是目录，直接返回文件信息
	if !file.IsDir() {
		data.Curdir = path.Dir(dir)
		data.Curname = file.Name()
		data.Isleaf = true
		data.Id = string(pid)
		pid += 1
		return data, err
	}
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(fileInfo)
		return data, err
	}

	//目录为空
	if len(fileInfo) == 0 {
		return
	}

	data.Curdir = dir
	_, data.Curname = path.Split(dir)
	data.Isleaf = false
	data.Id = string(pid)
	pid += 1

	for _, v := range fileInfo {
		if v.IsDir() {
			if subDir, err := Init_pages(dir + "/" + v.Name()); err != nil {
				return data, err
			} else {
				data.Branches = append(data.Branches, subDir)
				data.Isleaf = false
				data.Id = string(pid)
				pid += 1
			}
		} else {
			f := View_page_tree{}
			f.Curname = v.Name()
			f.Curdir = dir + "/" + v.Name()
			f.Isleaf = true
			f.Id = string(pid)
			pid += 1
			data.Branches = append(data.Branches, f)
		}
	}
	return data, err
}
