package services

import (
	"fmt"
	"testing"
)

func Test_Page(t *testing.T) {
    list:=MapItems{
    	{"a1",1},
		{"a2",2},
		{"a3",3},
		{"a4",4},
		{"a5",5},
		{"a6",6},
	}
   // fmt.Println(list[2:4])
    fmt.Println(list.Page(2,2))
}


