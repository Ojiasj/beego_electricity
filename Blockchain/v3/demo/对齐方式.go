package main

import (
	"fmt"
	"unsafe" //go语言的sizeof
)
// 低 --------》 高
// 12 34  -> 大端 -> 高尾端
// 34 12  -> 小端 -> 低尾端

func main() {
	s := int16(0x1234)
	b := int8(s)
	fmt.Println("int16字节大小为", unsafe.Sizeof(s)) //结果为2
	if 0x34 == b {
		fmt.Println("little endian")
	} else {
		fmt.Println("big endian")
	}
}
