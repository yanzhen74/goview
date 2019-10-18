package test

import (
	"fmt"
	"strings"
	"testing"
)

func Test_strings(t *testing.T) {
	s := "RTM_WYG_PK-CEH2_Result	00	0	02:03:04.3333"
	frame := fmt.Sprintf("%s\n1 aa 333;2 cc 222", s)

	lines := strings.Split(frame, "\n")
	fmt.Println(lines[0], lines[1])

	title := strings.Split(lines[0], "\t")
	fmt.Println(title[0])

	frame_type := strings.Split(title[0], "_")
	fmt.Println(frame_type)
}
