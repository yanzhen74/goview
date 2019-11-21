package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type View_page_tree struct {
	Id       string
	Curname  string
	Curdir   string
	Branches []View_page_tree
	Isleaf   bool
}

var page_id = 0

func Init_pages(dir string) (data View_page_tree, err error) {
	//判断文件或目录是否存在
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return data, err
	}

	data = View_page_tree{}

	//直接返回文件信息
	if !dirInfo.IsDir() {
		data.Curdir = dir
		data.Curname = dirInfo.Name()
		data.Isleaf = true
		data.Id = string(page_id)
		page_id += 1
		return data, err
	}

	// 获取目录包含的子目录和文件信息
	subInfos, err := ioutil.ReadDir(dir)
	if err != nil || len(subInfos) == 0 {
		fmt.Println(subInfos)
		return data, err
	}

	// 本目录信息
	data.Curdir = dir
	_, data.Curname = path.Split(dir)
	data.Isleaf = false
	data.Id = fmt.Sprintf("%d", page_id)
	page_id += 1

	// 子目录及文件信息
	for _, v := range subInfos {
		if v.IsDir() {
			if subDir, err := Init_pages(dir + "/" + v.Name()); err != nil {
				return data, err
			} else if subDir.Id == "" {
				continue
			} else {
				data.Branches = append(data.Branches, subDir)
			}
		} else {
			f := View_page_tree{}
			f.Curdir = dir + "/" + v.Name()
			f.Curname = v.Name()
			f.Isleaf = true
			f.Id = fmt.Sprintf("%d", page_id)
			page_id += 1
			data.Branches = append(data.Branches, f)
		}
	}
	return data, err
}
