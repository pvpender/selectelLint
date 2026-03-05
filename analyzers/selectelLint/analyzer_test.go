package selectelLint

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestSelectelLint(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		options map[string]string
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
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			an := NewAnalyzer()

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
