package main

import (
	"github.com/pvpender/selectellint/analyzers/selectelLint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(selectelLint.NewAnalyzer())
}
