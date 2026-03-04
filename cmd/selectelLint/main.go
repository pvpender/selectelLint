package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"pvpender/selectelLint/analyzers/selectelLint"
)

func main() {
	singlechecker.Main(selectelLint.NewAnalyzer())
}
