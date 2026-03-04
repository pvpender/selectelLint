package selectelLint

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestSelectelLint(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NewAnalyzer())
}
