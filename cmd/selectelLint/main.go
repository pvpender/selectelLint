package main

import (
	"github.com/pvpender/selectelLint/analyzers/selectelLint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(selectelLint.NewAnalyzer())
}
