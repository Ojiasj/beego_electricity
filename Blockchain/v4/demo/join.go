package main

import (
	"strings"
	"fmt"
	"bytes"
)

func main() {
	str1 := []string{"hello", "world", "!"}
	res := strings.Join(str1, "")
	fmt.Printf("res : %s\n", res)

	//func Join(s [][]byte, sep []byte) []byte {
	res1 := bytes.Join([][]byte{[]byte("hello"), []byte("world")}, []byte(""))
	fmt.Printf("res1 : %s\n", res1)
}



