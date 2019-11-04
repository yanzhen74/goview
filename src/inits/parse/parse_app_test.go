package parse

import (
	"fmt"
	"os"
	"testing"
)

func Test_app_other_parse(t *testing.T) {
	fmt.Println(os.Getwd())
	os.Chdir("..")
	fmt.Println(os.Getwd())
	os.Chdir("..")
	fmt.Println(os.Getwd())
	os.Chdir("..")
	fmt.Println(os.Getwd())
	AppOtherParse()
	fmt.Println(O)
}
