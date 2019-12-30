package main

//import "fmt"

func main() {
	bc := NewBlockChain("班长")
	cli := CLI{bc}
	cli.Run()
}
	//bc.AddBlock("111111111111111")
	//bc.AddBlock("222222222222222")
	//

