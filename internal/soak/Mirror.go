package soak

import (
	"os"
	"path"

	"github.com/phillip-england/wir/internal/wherr"
)

type Mirror struct {
	targetVFS  *Vfs
	compareVFS *Vfs
}

func NewMirror(targetDirRelPath []string, compareToDirRelPath []string) (Mirror, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Mirror{}, wherr.Consume(wherr.Here(), err, "")
	}
	targetDirPath := path.Join(append([]string{cwd}, targetDirRelPath...)...)
	targetVFS, err := LoadVfsAbsolute(true, targetDirPath)
	if err != nil {
		return Mirror{}, wherr.Consume(wherr.Here(), err, "")
	}
	compareToDirPath := path.Join(append([]string{cwd}, compareToDirRelPath...)...)
	compareVFS, err := LoadVfsAbsolute(true, compareToDirPath)
	if err != nil {
		return Mirror{}, wherr.Consume(wherr.Here(), err, "")
	}
	err = ensureDirsMirror(targetVFS, compareVFS)
	if err != nil {
		return Mirror{}, wherr.Consume(wherr.Here(), err, "")
	}
	return Mirror{
		targetVFS:  targetVFS,
		compareVFS: compareVFS,
	}, nil
}

func ensureDirsMirror(targetVFS *Vfs, compareVFS *Vfs) error {
	var potErr error
	targetVFS.IterAssets(func(targetAsset *VirtualAsset) bool {
		foundCompareAsset := false
		shouldBreak := false
		compareVFS.IterAssets(func(compareAsset *VirtualAsset) bool {
			if targetAsset.FileNameNoExt != compareAsset.FileNameNoExt {
				return true
			}
			foundCompareAsset = true
			return false
		})
		if !foundCompareAsset {
			potErr = wherr.Err(wherr.Here(), "failed to locate comparison asset for %s", targetAsset.Path)
			return false
		}
		if shouldBreak {
			return false
		}
		return true
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

func (d Mirror) Iter(fn func(target *VirtualAsset, compare *VirtualAsset) bool) {
	d.targetVFS.IterAssets(func(targetAsset *VirtualAsset) bool {
		d.compareVFS.IterAssets(func(compareAsset *VirtualAsset) bool {
			if targetAsset.FileNameNoExt == compareAsset.FileNameNoExt {
				shouldContinue := fn(targetAsset, compareAsset)
				if shouldContinue {
					return true
				}
				return false
			}
			return true
		})
		return true
	})
}
