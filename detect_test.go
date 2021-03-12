package hugobuildpack_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	hugobuildpack "github.com/fg-j/explorations/hugo-buildpack"
	"github.com/paketo-buildpacks/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir string
		detect     packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		Expect(os.MkdirAll(filepath.Join(workingDir, "content", "subdirectory"), os.ModePerm)).To(Succeed())
		Expect(err).NotTo(HaveOccurred())

		detect = hugobuildpack.Detect()
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("when the app contains a content directory with at least one *.md file in a subdirectory", func() {
		it.Before(func() {
			Expect(ioutil.WriteFile(filepath.Join(workingDir, "content", "subdirectory", "hello.md"), nil, os.ModePerm)).To(Succeed())
		})

		it("detects", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "hugo",
						Metadata: map[string]interface{}{
							"build": true,
						},
					},
				},
				Provides: []packit.BuildPlanProvision{
					{
						Name: "hugo",
					},
				},
			}))
		})
	})

	context("when the app contains a content directory with at least one *.html file", func() {
		it.Before(func() {
			Expect(ioutil.WriteFile(filepath.Join(workingDir, "content", "subdirectory", "index.html"), nil, os.ModePerm)).To(Succeed())
		})

		it("detects", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "hugo",
						Metadata: map[string]interface{}{
							"build": true,
						},
					},
				},
				Provides: []packit.BuildPlanProvision{
					{
						Name: "hugo",
					},
				},
			}))
		})
	})

	context("when the app contains no content/*.[md|html] files", func() {
		it("fails detection", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(ContainSubstring("no *.md or *.html files found in content dir")))
		})
	})
	context("failure cases", func() {
		context("search for *md and *html files fails", func() {
			it("returns the error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: `\/\/`,
				})
				Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
			})
		})
	})
}
