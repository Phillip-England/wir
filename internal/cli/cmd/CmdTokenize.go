package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/phillip-england/wir/internal/mood"
	"github.com/phillip-england/wir/internal/soak"
	"github.com/phillip-england/wir/internal/wherr"
	"github.com/phillip-england/wir/internal/wirtokenizer"
)

type CmdTokenize struct{
	argInPath string
	argOutPath string
	inPathAbs string
	outPathAbs string
	shouldOverwrite bool
	isTargetingDir bool
}

func NewCmdTokenize(cli *mood.Cli) (mood.Cmd, error) {
	argInPath, err := cli.ArgGetByPositionForce(2, "missing <INPUT_FILE> for wir tokenize")
	if err != nil {
		return CmdTokenize{}, wherr.Consume(wherr.Here(), err, "")
	}
	argOutPath, err := cli.ArgGetByPositionForce(3, "missing <OUTPUT_FILE> for wir tokenize")
	if err != nil {
		return CmdTokenize{}, wherr.Consume(wherr.Here(), err, "")
	}
	inPathAbs := path.Join(cli.Cwd, argInPath)
	if !mood.FileExists(inPathAbs) {
		return CmdTokenize{}, wherr.Err(wherr.Here(), "<INPUT_FILE> does not exist in wir tokenize")
	}
	isTargetingDir := false
	if mood.IsDir(inPathAbs) {
		isTargetingDir = true
	}
	if isTargetingDir {
		if !mood.IsDir(argOutPath) {
			return CmdTokenize{}, wherr.Err(wherr.Here(), "<OUTPUT_FILE> must be a directory if <INPUT_FILE> is a directory")
		}
	} else {
		if !mood.IsFile(argOutPath) {
			return CmdTokenize{}, wherr.Err(wherr.Here(), "<OUTPUT_FILE> must be a file path if <INPUT_FILE> is a file")
		}
	}
	return CmdTokenize{
		argInPath: argInPath,
		argOutPath: argOutPath,
		inPathAbs: inPathAbs,
		outPathAbs: path.Join(cli.Cwd, argOutPath),
		shouldOverwrite: cli.FlagExists("-o"),
		isTargetingDir: isTargetingDir,
	}, nil
}

func (cmd CmdTokenize) Execute(cli *mood.Cli) error {
	if cmd.isTargetingDir {
		err := tokenizeDir(cmd)
		if err != nil {
			return wherr.Consume(wherr.Here(), err, "")
		}
	} else {
		err := tokenizeFile(cmd)
		if err != nil {
			return wherr.Consume(wherr.Here(), err, "")
		}
	}
	return nil
}

func tokenizeFile(cmd CmdTokenize) error {
	if cmd.shouldOverwrite {
		os.RemoveAll(cmd.argOutPath)
	}
	if mood.FileExists(cmd.argOutPath) {
		return wherr.Err(wherr.Here(), "file already exists at %s", cmd.outPathAbs)
	}
	fBytes, err := os.ReadFile(cmd.inPathAbs)
	if err != nil {
		return wherr.Consume(wherr.Here(), err, "")
	}
	fStr := string(fBytes)
	tk, err := wirtokenizer.TokenizerNewFromString(fStr)
	if err != nil {
		return wherr.Consume(wherr.Here(), err, "")
	}
	err = os.WriteFile(cmd.outPathAbs, []byte(tk.Str()), 0775)
	if err != nil {
		return wherr.Consume(wherr.Here(), err, "")
	}
	return nil
}

func tokenizeDir(cmd CmdTokenize) error {
	if cmd.shouldOverwrite {
		os.RemoveAll(cmd.outPathAbs)
	}
	if mood.FileExists(cmd.outPathAbs) {
		return wherr.Err(wherr.Here(), "dir already exists at %s", cmd.outPathAbs)
	}
	vfs, err := soak.LoadVfsAbsolute(true, cmd.inPathAbs)
	if err != nil {
		return wherr.Consume(wherr.Here(), err, "")
	}
	var potErr error
	vfs.IterAssets(func(a *soak.VirtualAsset) bool {
		outDirPath := path.Join(cmd.outPathAbs, a.FileNameNoExt+".tok")
		tk, err := wirtokenizer.TokenizerNewFromString(a.Text)
		if err != nil {
			potErr = wherr.Consume(wherr.Here(), err, "")
		}
		err = os.MkdirAll(path.Dir(outDirPath), 0755)
		if err != nil {
			potErr =  wherr.Consume(wherr.Here(), err, "")
		}
		err = os.WriteFile(outDirPath, []byte(tk.Str()), 0644)
		if err != nil {
			fmt.Println(err.Error())
			potErr =  wherr.Consume(wherr.Here(), err, "")
		}
		return true
	})
	if potErr != nil {
		return potErr
	}
	return nil
}