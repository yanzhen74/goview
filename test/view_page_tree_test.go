package test

import (
	"fmt"
	"testing"

	"github.com/yanzhen74/goview/src/model"
)

func Test_init_page_tree(t *testing.T) {
	z, err := model.Init_pages("../config/resource/menu")
	if err != nil {
		fmt.Printf("error %v", err)
	}
	fmt.Println("hello ", z)
}
