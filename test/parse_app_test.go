package test

import (
	"fmt"
	"testing"

	"github.com/yanzhen74/goview/src/inits/parse"
)

func Test_app_other_parse(t *testing.T) {
	parse.AppOtherParse()
	fmt.Println(parse.O)
}
