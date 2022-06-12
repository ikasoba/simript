package main

import (
	"fmt"

	"github.com/ikasoba/simript"
)

func main(){
	root, err := simript.CreateASTFromTokens(simript.Tokenize(`hoge(1234,"moji",3.14)`))
	fmt.Println(root,err)
}