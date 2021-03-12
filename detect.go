package hugobuildpack

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
)

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		var files []string

		err := filepath.Walk(filepath.Join(context.WorkingDir, "content"), func(path string, info os.FileInfo, err error) error {
			_, result := os.Stat(path)
			if err != nil {
				return result
			}
			if ext := filepath.Ext(path); ext == ".md" || ext == ".html" {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			return packit.DetectResult{}, fmt.Errorf("searching for *.md and *.html files in %s: %w", filepath.Join(context.WorkingDir, "content"), err)
		}

		if len(files) == 0 {
			return packit.DetectResult{}, packit.Fail.WithMessage("no *.md or *.html files found in content dir")
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
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
			},
		}, nil
	}
}
