package soak

import (
	"os"
	"path"

	"github.com/phillip-england/wir/internal/wherr"
)

type Vfs struct {
	isLocked bool
	Path     string
	Assets   []*VirtualAsset
}

func LoadVfs(isLocked bool, relPath ...string) (*Vfs, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, wherr.Consume(wherr.Here(), err, "")
	}
	outParts := append([]string{}, cwd)
	outParts = append(outParts, relPath...)
	outPath := path.Join(outParts...)
	assets, err := LoadVirtualAssets(isLocked, cwd, outPath)
	if err != nil {
		return nil, wherr.Consume(wherr.Here(), err, "")
	}
	return &Vfs{
		isLocked: isLocked,
		Path:     outPath,
		Assets:   assets,
	}, nil
}

func LoadVfsAbsolute(isLocked bool, absPath string) (*Vfs, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, wherr.Consume(wherr.Here(), err, "")
	}
	assets, err := LoadVirtualAssets(isLocked, cwd, absPath)
	if err != nil {
		return nil, wherr.Consume(wherr.Here(), err, "")
	}
	return &Vfs{
		isLocked: isLocked,
		Path:     absPath,
		Assets:   assets,
	}, nil
}

func (v *Vfs) IterAssets(fn func(a *VirtualAsset) bool) {
	for _, asset := range v.Assets {
		shouldContinue := fn(asset)
		if shouldContinue {
			continue
		}
		break
	}
}

func (v *Vfs) Sync() error {
	if v.isLocked {
		return wherr.Err(wherr.Here(), "attempted to sync a locked virtual file system")
	}
	var potErr error
	v.IterAssets(func(a *VirtualAsset) bool {
		err := a.Save()
		if err != nil {
			potErr = wherr.Consume(wherr.Here(), err, "")
			return false
		}
		return true
	})
	if potErr != nil {
		return potErr
	}
	return nil
}
