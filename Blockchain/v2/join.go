package main

import (
	"strings"
	"fmt"
	"bytes"
)

func test() {
	str1 := []string{"hello", "world", "!"}
	res := strings.Join(str1, "-")
	fmt.Printf("%s\n", res)

	res2 := bytes.Join([][]byte{[]byte("hello"), []byte("world")}, []byte(""))
	fmt.Printf("%s\n", string(res2))
}
