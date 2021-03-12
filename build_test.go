package hugobuildpack_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	hugobuildpack "github.com/fg-j/explorations/hugo-buildpack"
	"github.com/fg-j/explorations/hugo-buildpack/fakes"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/postal"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		buffer            *bytes.Buffer
		cnbDir            string
		dependencyManager *fakes.DependencyManager
		entryResolver     *fakes.EntryResolver
		executable        *fakes.Executable
		layersDir         string
		workingDir        string

		build packit.BuildFunc
	)

	it.Before(func() {
		var err error
		layersDir, err = ioutil.TempDir("", "layers")
		Expect(err).NotTo(HaveOccurred())

		cnbDir, err = ioutil.TempDir("", "cnb")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		entryResolver = &fakes.EntryResolver{}
		entryResolver.ResolveCall.Returns.BuildpackPlanEntry = packit.BuildpackPlanEntry{
			Name: "hugo",
		}

		dependencyManager = &fakes.DependencyManager{}
		dependencyManager.ResolveCall.Returns.Dependency = postal.Dependency{
			ID:      "hugo",
			Name:    "hugo-dependency-name",
			SHA256:  "hugo-dependency-sha",
			Stacks:  []string{"some-stack"},
			URI:     "hugo-dependency-uri",
			Version: "hugo-dependency-version",
		}

		executable = &fakes.Executable{}

		buffer = bytes.NewBuffer(nil)
		logEmitter := scribe.NewEmitter(buffer)

		build = hugobuildpack.Build(entryResolver, dependencyManager, executable, logEmitter)
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(cnbDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	it("returns a result that installs hugo and builds correctly", func() {
		result, err := build(packit.BuildContext{
			WorkingDir: workingDir,
			CNBPath:    cnbDir,
			Stack:      "some-stack",
			BuildpackInfo: packit.BuildpackInfo{
				Name:    "Some Buildpack",
				Version: "some-version",
			},
			Plan: packit.BuildpackPlan{
				Entries: []packit.BuildpackPlanEntry{
					{
						Name: "hugo",
					},
				},
			},
			Layers: packit.Layers{Path: layersDir},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(result).To(Equal(packit.BuildResult{
			Plan: packit.BuildpackPlan{},
			Layers: []packit.Layer{
				{
					Name:      "hugo",
					Path:      filepath.Join(layersDir, "hugo"),
					SharedEnv: packit.Environment{},
					BuildEnv:  packit.Environment{},
					LaunchEnv: packit.Environment{},
					Build:     false,
					Launch:    false,
					Cache:     false,
				},
			},
		}))
	})

	context("failure cases", func() {
		context("dependency resolution fails", func() {
			it.Before(func() {
				dependencyManager.ResolveCall.Returns.Error = errors.New("dependency resolution error")
			})

			it("fails to build with the appropriate error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					CNBPath:    cnbDir,
					Stack:      "some-stack",
					BuildpackInfo: packit.BuildpackInfo{
						Name:    "Some Buildpack",
						Version: "some-version",
					},
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{
								Name: "hugo",
							},
						},
					},
					Layers: packit.Layers{Path: layersDir},
				})

				Expect(err).To(MatchError(ContainSubstring("dependency resolution error")))
			})
		})

		context("layer cannot be read/written", func() {
			it.Before(func() {
				Expect(os.Chmod(layersDir, 0000)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Chmod(layersDir, os.ModePerm)).To(Succeed())
			})

			it("fails to build with the appropriate error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					CNBPath:    cnbDir,
					Stack:      "some-stack",
					BuildpackInfo: packit.BuildpackInfo{
						Name:    "Some Buildpack",
						Version: "some-version",
					},
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{
								Name: "hugo",
							},
						},
					},
					Layers: packit.Layers{Path: layersDir},
				})

				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})

		context("dependency installation fails", func() {
			it.Before(func() {
				dependencyManager.InstallCall.Returns.Error = errors.New("dependency installation error")
			})

			it("fails to build with the appropriate error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					CNBPath:    cnbDir,
					Stack:      "some-stack",
					BuildpackInfo: packit.BuildpackInfo{
						Name:    "Some Buildpack",
						Version: "some-version",
					},
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{
								Name: "hugo",
							},
						},
					},
					Layers: packit.Layers{Path: layersDir},
				})

				Expect(err).To(MatchError(ContainSubstring("dependency installation error")))
			})
		})

		context("hugo process execution fails", func() {
			it.Before(func() {
				executable.ExecuteCall.Returns.Error = errors.New("hugo error")
			})

			it("fails to build with the appropriate error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					CNBPath:    cnbDir,
					Stack:      "some-stack",
					BuildpackInfo: packit.BuildpackInfo{
						Name:    "Some Buildpack",
						Version: "some-version",
					},
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{
								Name: "hugo",
							},
						},
					},
					Layers: packit.Layers{Path: layersDir},
				})

				Expect(err).To(MatchError(ContainSubstring("hugo error")))
			})
		})
	})
}
