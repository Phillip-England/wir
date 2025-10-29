package main

import (
	"fmt"

	"github.com/phillip-england/wir/internal/cli/cmd"
	"github.com/phillip-england/wir/internal/mood"
)

func main() {


	cli, err := mood.New(cmd.NewCmdDefault)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	cli.At("tokenize", cmd.NewCmdTokenize)

	err = cli.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

