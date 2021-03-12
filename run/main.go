package main

import (
	"os"

	hugobuildpack "github.com/fg-j/explorations/hugo-buildpack"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/cargo"
	"github.com/paketo-buildpacks/packit/draft"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/scribe"
)

func main() {
	logEmitter := scribe.NewEmitter(os.Stdout)
	entryResolver := draft.NewPlanner()
	dependencyManager := hugobuildpack.NewHugoDependencyManager(cargo.NewTransport())
	packit.Run(hugobuildpack.Detect(),
		hugobuildpack.Build(entryResolver, dependencyManager, pexec.NewExecutable("hugo"), logEmitter),
	)
}
