package soak

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/phillip-england/wir/internal/wherr"
)

type VirtualAsset struct {
	Path          string
	Dirname       string
	Text          string
	Ext           string
	RelPath       string
	FileName      string
	isLocked      bool
	FileNameNoExt string
}

func UnixRelPath(basePath, targetPath string) (string, error) {
	relPath, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return "", wherr.Consume(wherr.Here(), err, "")
	}
	unixPath := filepath.ToSlash(relPath)
	if !strings.HasPrefix(unixPath, "./") && !strings.HasPrefix(unixPath, "../") {
		unixPath = "./" + unixPath
	}
	return unixPath, nil
}

func LoadVirtualAssets(isLocked bool, cwd string, pth string) ([]*VirtualAsset, error) {
	var assets []*VirtualAsset
	err := filepath.WalkDir(pth, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return wherr.Consume(wherr.Here(), err, "")
		}
		if d.IsDir() {
			return nil
		}
		dirname := path.Dir(p)
		fBytes, err := os.ReadFile(p)
		if err != nil {
			return wherr.Consume(wherr.Here(), err, "")
		}
		relPath, err := UnixRelPath(cwd, p)
		if err != nil {
			return wherr.Consume(wherr.Here(), err, "")
		}
		ext := path.Ext(p)
		filename := path.Base(p)
		asset := &VirtualAsset{
			Path:          p,
			Dirname:       dirname,
			Text:          string(fBytes),
			Ext:           ext,
			RelPath:       relPath,
			FileName:      filename,
			isLocked:      isLocked,
			FileNameNoExt: strings.TrimSuffix(filename, ext),
		}
		assets = append(assets, asset)
		return nil
	})
	if err != nil {
		return assets, wherr.Consume(wherr.Here(), err, "")
	}
	return assets, nil
}

func (a *VirtualAsset) Save() error {
	if a.isLocked {
		return wherr.Err(wherr.Here(), "attemped to save a locked virtual asset")
	}
	err := os.WriteFile(a.Path, []byte(a.Text), 0644)
	if err != nil {
		return wherr.Consume(wherr.Here(), err, "")
	}
	return nil
}

func (a *VirtualAsset) OverWrite(s string) error {
	if a.isLocked {
		return wherr.Err(wherr.Here(), "attemped to overwrite a locked virtual asset")
	}
	a.Text = s
	return nil
}

func (a *VirtualAsset) Append(s string) error {
	if a.isLocked {
		return wherr.Err(wherr.Here(), "attemped to append to a locked virtual asset")
	}
	a.Text = a.Text + s
	return nil
}

func (a *VirtualAsset) Prepend(s string) error {
	if a.isLocked {
		return wherr.Err(wherr.Here(), "attemped to prepend to a locked virtual asset")
	}
	a.Text = s + a.Text
	return nil
}
