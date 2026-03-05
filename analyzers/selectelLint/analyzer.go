package selectelLint

import (
	"flag"
	detector "github.com/kevinwang15/sensitive-data-detector"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"pvpender/selectelLint/config"
	"regexp"
	"strings"
	"unicode"
)

type Analyzer struct {
	config *config.Config
}

func NewAnalyzer() *analysis.Analyzer {
	an := &Analyzer{
		config: config.NewConfig(),
	}

	return &analysis.Analyzer{
		Name:     "CapitalAnalyzer",
		Doc:      "Checks that there are no capitals in logs",
		Flags:    an.Flags(),
		Run:      an.Run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func (an *Analyzer) Flags() flag.FlagSet {
	flags := flag.NewFlagSet("CapitalAnalyzer", flag.ExitOnError)
	flags.BoolVar(&an.config.CapitalLetter, "capitalLetter", true, "Enable capital letter check")
	flags.BoolVar(&an.config.EnglishLetter, "englishLetter", true, "Enable english letter check")
	flags.BoolVar(&an.config.SpecialLetters, "specialLetter", true, "Enable special letter check")
	flags.BoolVar(&an.config.SensitiveData, "sensitiveData", true, "Enable sensitive data check")
	flags.BoolVar(&an.config.EnableCustomRules, "enableCustomRules", false, "Enable custom rules")

	return *flags
}

func (an *Analyzer) Run(pass *analysis.Pass) (interface{}, error) {
	inspectorVar := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspectorVar.Preorder(nodeFilter, func(node ast.Node) {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return
		}

		if !an.isLogFunction(call) {
			return
		}

		message, exists := an.getMessage(call)

		if !exists {
			return
		}

		if an.config.CapitalLetter {
			if msg, fail := an.checkCapital(message); fail {
				pass.Reportf(call.Pos(), "%s", msg)
			}
		}

		if an.config.EnglishLetter {
			if msg, fail := an.checkEnglish(message); fail {
				pass.Reportf(call.Pos(), "%s", msg)
			}
		}

		if an.config.SpecialLetters {
			if msg, fail := an.checkSpecialLetters(message); fail {
				pass.Reportf(call.Pos(), "%s", msg)
			}
		}

		if an.config.SensitiveData {
			if msg, fail := an.checkSensitiveData(message); fail {
				pass.Reportf(call.Pos(), "%s", msg)
			}
		}

	})

	return nil, nil
}

func (an *Analyzer) isLogFunction(call *ast.CallExpr) bool {
	logMethods := map[string]bool{
		"Info": true,
	}

	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:

		return logMethods[fun.Sel.Name]

	case *ast.Ident:
		return logMethods[fun.Name]
	}

	return false
}

func (an *Analyzer) getMessage(expr ast.Expr) (string, bool) {
	switch lit := expr.(type) {
	case *ast.CallExpr:
		if len(lit.Args) == 0 {
			return "", false
		}

		for _, arg := range lit.Args {
			if msg, ok := an.getMessage(arg); ok {
				return strings.Trim(msg, `"'`), true
			}
		}

		break
	case *ast.BasicLit:
		if lit.Kind == token.STRING {
			return strings.Trim(lit.Value, `"'`), true
		}

		break

	case *ast.BinaryExpr:
		if lit.Op == token.ADD {
			left, leftOk := an.getMessage(lit.X)
			right, rightOk := an.getMessage(lit.Y)

			if leftOk && rightOk {
				return left + right, true
			}

			if leftOk {
				return left, true
			}
		}

		break

	case *ast.Ident:
		if lit.Obj == nil {
			return "", false
		}

		if value, ok := lit.Obj.Decl.(*ast.ValueSpec); ok {
			for _, v := range value.Values {
				return an.getMessage(v)
			}
		}
	}

	return "", false
}

func (an *Analyzer) checkCapital(msg string) (string, bool) {
	if unicode.IsUpper([]rune(msg)[0]) && unicode.IsLetter([]rune(msg)[0]) {
		return "Start uppercase detected!", true
	}

	return "", false
}

func (an *Analyzer) checkEnglish(msg string) (string, bool) {
	for _, char := range msg {
		if match, _ := regexp.MatchString(`[ -~]`, string(char)); !match && unicode.IsLetter(char) {
			return "Not english character detected!", true
		}
	}

	return "", false
}

func (an *Analyzer) checkSpecialLetters(msg string) (string, bool) {
	for _, char := range msg {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) && !unicode.IsDigit(char) {
			return "Special letters detected!", true
		}
	}

	return "", false
}

func (an *Analyzer) checkSensitiveData(msg string) (string, bool) {
	d, err := detector.NewDetector(
		detector.WithPatterns(detector.Pattern{
			Name:        "Custom sensitive",
			Description: "Value that looks like a password, secret, or API key assignment",
			Expression:  `(?i)(token|secret|password|passwd|api[_\-]?key)\s*[:=]\s*(?:['"]?(?:%s|[A-Za-z0-9_\-/+=)]{6,})['"]?|$)`,
			Severity:    detector.SeverityHigh,
			Types:       []string{"credential"},
		}),
	)

	if err != nil {
		return "", false
	}

	violations, err := d.Scan(msg)

	if err != nil {
		return "", false
	}

	for _, violation := range violations {
		return violation.Description, true
	}

	return "", false
}
