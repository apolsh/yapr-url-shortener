package exitmain

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestExitMain(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "./...")
}
