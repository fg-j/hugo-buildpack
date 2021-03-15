package hugobuildpack

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/postal"
	"github.com/paketo-buildpacks/packit/scribe"
)

//go:generate faux --interface EntryResolver --output fakes/entry_resolver.go
type EntryResolver interface {
	Resolve(name string, entries []packit.BuildpackPlanEntry, priorities []interface{}) (packit.BuildpackPlanEntry, []packit.BuildpackPlanEntry)
	MergeLayerTypes(name string, entries []packit.BuildpackPlanEntry) (bool, bool)
}

//go:generate faux --interface DependencyManager --output fakes/dependency_manager.go
type DependencyManager interface {
	Resolve(path, id, version, stack string) (postal.Dependency, error)
	Install(dependency postal.Dependency, cnbPath, layerPath string) error
}

//go:generate faux --interface Executable --output fakes/executable.go
type Executable interface {
	Execute(pexec.Execution) error
}

func Build(entryResolver EntryResolver, dependencyManager DependencyManager, hugo Executable, logs scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logs.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		entry, _ := entryResolver.Resolve("hugo", context.Plan.Entries, nil)
		entryVersion, _ := entry.Metadata["version"].(string)

		dependency, err := dependencyManager.Resolve(
			filepath.Join(context.CNBPath, "buildpack.toml"),
			entry.Name,
			entryVersion,
			context.Stack)
		if err != nil {
			return packit.BuildResult{}, err
		}

		// install hugo
		hugoLayer, err := context.Layers.Get("hugo")
		if err != nil {
			return packit.BuildResult{}, err
		}

		hugoLayer, err = hugoLayer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		launch, build := entryResolver.MergeLayerTypes("hugo", context.Plan.Entries)
		hugoLayer.Build = build
		hugoLayer.Cache = build
		hugoLayer.Launch = launch

		logs.Process("Installing Hugo")

		err = dependencyManager.Install(dependency, context.CNBPath, hugoLayer.Path)
		if err != nil {
			return packit.BuildResult{}, err
		}

		logs.Break()

		// build static assets
		logs.Process("Executing build process")
		buffer := bytes.NewBuffer(nil)
		args := []string{
			"--destination", "public",
		}

		logs.Subprocess("Running 'hugo %s'", strings.Join(args, " "))

		err = hugo.Execute(pexec.Execution{
			Args:   args,
			Dir:    context.WorkingDir,
			Stdout: buffer,
			Stderr: buffer,
			Env:    append(os.Environ(), fmt.Sprintf("PATH=%s:%s", filepath.Join(hugoLayer.Path, "bin"), os.Getenv("PATH"))),
		})

		if err != nil {
			return packit.BuildResult{}, err
		}

		return packit.BuildResult{
			Plan:   packit.BuildpackPlan{},
			Layers: []packit.Layer{hugoLayer},
		}, nil
	}
}
