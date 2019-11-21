package model

import (
	"fmt"
	"os"
	"testing"
)

func Test_init_page_tree(t *testing.T) {
	z, err := Init_pages("../../config/resource/menu")
	if err != nil {
		fmt.Printf("error %v", err)
	}
	fmt.Println("hello ", z)
	fmt.Println("...", z.Branches[0].Branches[0].Branches[0].Branches[0].Curdir)
	fmt.Println("...", z.Branches[1].Branches[0].Branches[0].Branches[0].Curdir)
}

func Test_init_role_page_tree(t *testing.T) {
	z, err := Init_pages("../../config/resource/role")
	if err != nil {
		fmt.Printf("error %v", err)
	}
	fmt.Println(z)
	fmt.Println("...", z.Branches[0].Branches[0].Curdir)
}
func Test_stat(t *testing.T) {
	dirInfo, _ := os.Stat("view_page_tree_test.go")
	fmt.Printf("%s\n", dirInfo.Name())
}
