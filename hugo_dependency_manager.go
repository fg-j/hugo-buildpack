package hugobuildpack

import (
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/cargo"
	"github.com/paketo-buildpacks/packit/fs"
	"github.com/paketo-buildpacks/packit/postal"
)

type HugoDependencyManager struct {
	postal.Service
}

func NewHugoDependencyManager(transport cargo.Transport) HugoDependencyManager {
	return HugoDependencyManager{
		Service: postal.NewService(transport),
	}
}

func (h HugoDependencyManager) Install(dependency postal.Dependency, cnbPath, layerPath string) error {
	err := h.Service.Install(dependency, cnbPath, layerPath)
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(layerPath, "bin"), os.ModePerm)
	if err != nil {
		return err
	}

	err = fs.Move(filepath.Join(layerPath, "hugo"), filepath.Join(layerPath, "bin", "hugo"))
	if err != nil {
		return err
	}

	return nil
}
