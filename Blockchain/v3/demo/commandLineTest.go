package main

import (
	"os"
	"fmt"
)

func main() {
	len1 := len(os.Args)

	fmt.Printf("cmd len : %d\n", len1)
	for i, cmd := range os.Args {
		fmt.Printf("arg[%d] : %s\n", i, cmd)
	}
}
