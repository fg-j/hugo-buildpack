package hugobuildpack_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitHugoBuildpack(t *testing.T) {
	suite := spec.New("hugo-buildpack", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Build", testBuild)
	suite("Detect", testDetect)
	suite.Run(t)
}
