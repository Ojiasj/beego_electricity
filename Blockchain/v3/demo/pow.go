package main

import (
	"crypto/sha256"
	"fmt"
)

func main()  {

	//交易数据
	data := "helloworld"

	for i:= 0; i< 1000000; i++ {
		hash := sha256.Sum256([]byte(data + string(i)))

		fmt.Printf("hash : %x\n", hash[:])
	}
}
