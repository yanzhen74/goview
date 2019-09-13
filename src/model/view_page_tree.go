package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type View_page_tree struct {
	Curdir   string
	Branches []View_page_tree
	Leaves   []string
}

func Init(dir string) (data View_page_tree, err error) {
	//判断文件或目录是否存在
	file, err := os.Stat(dir)
	if err != nil {
		return data, err
	}
	data = View_page_tree{}

	//如果不是目录，直接返回文件信息
	if !file.IsDir() {
		data.Curdir = path.Dir(dir)
		data.Leaves = append(data.Leaves, file.Name())
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

	data.Curdir = path.Dir(dir)

	for _, v := range fileInfo {
		if v.IsDir() {
			if subDir, err := Init(dir + "/" + v.Name()); err != nil {
				return data, err
			} else {
				data.Branches = append(data.Branches, subDir)
			}
		} else {
			data.Leaves = append(data.Leaves, v.Name())
		}
	}
	return data, err
}
