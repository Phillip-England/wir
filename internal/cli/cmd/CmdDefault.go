package cmd

import (
	"fmt"

	"github.com/phillip-england/wir/internal/mood"
)

type CmdDefault struct{}

func NewCmdDefault(cli *mood.Cli) (mood.Cmd, error) {
	return CmdDefault{}, nil
}

func (cmd CmdDefault) Execute(cli *mood.Cli) error {
	fmt.Println(`[webIR]: a language for expressing reactive web UI's across multiple platforms
[tokenize example/usage]:
  -wir tokenize <INPUT_FILE> <OUTPUT_FILE>
  -wir tokenize ./input.wir ./output.txt`)
	return nil
}