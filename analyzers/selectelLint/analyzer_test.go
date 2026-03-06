package selectelLint

import (
	"testing"

	"github.com/pvpender/selectelLint/config"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSelectelLint(t *testing.T) {
	testCases := []struct {
		name        string
		path        string
		options     map[string]string
		customRules []config.Rule
	}{
		{
			name: "capital",
			path: analysistest.TestData() + "/capital",
			options: map[string]string{
				"capitalLetter": "true",
			},
		},
		{
			name: "English",
			path: analysistest.TestData() + "/english",
			options: map[string]string{
				"englishLetter": "true",
			},
		},
		{
			name: "Special",
			path: analysistest.TestData() + "/special",
			options: map[string]string{
				"specialLetter": "true",
			},
		},
		{
			name: "Sensitive",
			path: analysistest.TestData() + "/sensitive",
			options: map[string]string{
				"sensitiveData": "true",
				"specialLetter": "false",
			},
		},
		{
			name: "Custom",
			path: analysistest.TestData() + "/custom",
			options: map[string]string{
				"enableCustomRules": "true",
			},
			customRules: []config.Rule{
				{
					Name:        "Digit special password",
					Description: "Digit detected!",
					Pattern:     "\\d+",
				},
				{
					Name:        "Author",
					Description: "Author detected!",
					Pattern:     "(\\s+)?pvpender(\\s+)?",
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			cfg := config.NewConfig()

			if len(test.customRules) > 0 {
				cfg.EnableCustomRules = true
				for _, rule := range test.customRules {
					cfg.Rules = append(cfg.Rules, config.Rule{
						Name:        rule.Name,
						Description: rule.Description,
						Pattern:     rule.Pattern,
					})
				}
			}

			an := NewAnalyzer(cfg)

			for k, v := range test.options {
				err := an.Flags.Set(k, v)
				if err != nil {
					t.Fatal(err)
				}
			}

			analysistest.Run(t, test.path, an)
		})
	}
}
