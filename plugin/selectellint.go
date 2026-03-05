package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/pvpender/selectelLint/analyzers/selectelLint"
	"github.com/pvpender/selectelLint/config"
	"golang.org/x/tools/go/analysis"
)

type SelectelLintPlugin struct {
	analyzer *analysis.Analyzer
}

func init() {
	register.Plugin("sclint", New)
}

func New(conf any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[config.Config](conf)
	if err != nil {
		return nil, err
	}

	an := selectelLint.NewAnalyzer(&settings)

	return &SelectelLintPlugin{analyzer: an}, nil
}

func (plug *SelectelLintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{plug.analyzer}, nil
}

func (plug *SelectelLintPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
