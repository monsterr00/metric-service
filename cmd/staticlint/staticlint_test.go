package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/printf"
)

func TestOsExitCheckAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), printf.Analyzer, "./testdata.go")
}
