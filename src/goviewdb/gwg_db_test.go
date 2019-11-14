package goviewdb

import (
	"fmt"
	"testing"
)

func TestGwgDbNew(t *testing.T) {
	db := NewGWGDb("../../db/gwg.db")
	db.ListPic(10, 1, "")
}

func TestSaveGwgPic(t *testing.T) {
	db := NewGWGDb("../../db/gwg.db")
	db.SavePic("/data/gwg/1_17.jpeg")
	db.SavePic("/data/gwg/1_17.jpeg")
	db.SavePic("/data/gwg/1_17.jpeg")
	picList := db.ListPic(10, 1, "")
	for _, pic := range picList {
		fmt.Println(pic)
	}
}

func TestListGwgPic(t *testing.T) {
	db := NewGWGDb("../../db/gwg.db")
	picList := db.ListPic(10, 1, "")
	for _, pic := range picList {
		fmt.Println(pic)
	}
}

func TestClearGwgPic(t *testing.T) {
	db := NewGWGDb("../../db/gwg.db")
	db.ORM.Where("Id > ?", "0").Delete(GWGPic{})
}

func TestCountGwgPic(t *testing.T) {
	db := NewGWGDb("../../db/gwg.db")
	fmt.Printf("%v\n", db.CountPic())
}
