package selectelLint

import (
	"flag"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
	"unicode"

	detector "github.com/kevinwang15/sensitive-data-detector"
	"github.com/pvpender/selectelLint/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type Analyzer struct {
	config *config.Config
}

func NewAnalyzer(conf ...*config.Config) *analysis.Analyzer {
	var cfg *config.Config

	if len(conf) == 0 {
		cfg = config.NewConfig()
	} else {
		cfg = conf[0]
	}

	an := &Analyzer{
		config: cfg,
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
	flags.BoolVar(&an.config.CapitalLetter, "capitalLetter", an.config.CapitalLetter, "Enable capital letter check")
	flags.BoolVar(&an.config.EnglishLetter, "englishLetter", an.config.EnglishLetter, "Enable english letter check")
	flags.BoolVar(&an.config.SpecialLetters, "specialLetter", an.config.SpecialLetters, "Enable special letter check")
	flags.BoolVar(&an.config.SensitiveData, "sensitiveData", an.config.SensitiveData, "Enable sensitive data check")
	flags.BoolVar(&an.config.EnableCustomRules, "enableCustomRules", an.config.EnableCustomRules, "Enable custom rules")

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

		an.runChecker(an.config.CapitalLetter, pass, call, an.checkCapital, message)
		an.runChecker(an.config.EnglishLetter, pass, call, an.checkEnglish, message)
		an.runChecker(an.config.SpecialLetters, pass, call, an.checkSpecialLetters, message)
		an.runChecker(an.config.SensitiveData, pass, call, an.checkSensitiveData, message)

		if an.config.EnableCustomRules {
			for _, rule := range an.config.Rules {
				if msg, fail := an.checkCustomRule(message, &rule); fail {
					pass.Reportf(call.Pos(), "%s", msg)
				}
			}
		}
	})

	return nil, nil //nolint:nilnil
}

func (an *Analyzer) isLogFunction(call *ast.CallExpr) bool {
	logMethods := map[string]bool{
		"Debug":        true,
		"DebugContext": true,
		"DebugCtx":     true,
		"Error":        true,
		"ErrorContext": true,
		"ErrorCtx":     true,
		"Info":         true,
		"InfoContext":  true,
		"InfoCtx":      true,
		"Log":          true,
		"LogAttrs":     true,
		"Warn":         true,
		"WarnContext":  true,
		"WarnCtx":      true,
		"DPanic":       true,
		"DPanicf":      true,
		"DPanicw":      true,
		"Fatal":        true,
		"Fatalf":       true,
		"Fatalw":       true,
		"Panic":        true,
		"Debugf":       true,
		"Debugw":       true,
		"Errorf":       true,
		"Errorw":       true,
		"Infof":        true,
		"Infow":        true,
		"Logf":         true,
		"Logw":         true,
		"Panicf":       true,
		"Panicw":       true,
		"Warnf":        true,
		"Warnw":        true,
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

	case *ast.BasicLit:
		if lit.Kind == token.STRING {
			return strings.Trim(lit.Value, `"'`), true
		}

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

func (an *Analyzer) runChecker(
	shouldRun bool,
	pass *analysis.Pass,
	call *ast.CallExpr,
	checker func(string) (string, bool),
	message string,
) {
	if !shouldRun {
		return
	}

	if msg, fail := checker(message); fail {
		pass.Reportf(call.Pos(), "%s", msg)
	}
}

func (an *Analyzer) checkCapital(msg string) (string, bool) {
	if unicode.IsUpper([]rune(msg)[0]) && unicode.IsLetter([]rune(msg)[0]) {
		return "Start uppercase detected!", true
	}

	return "", false
}

func (an *Analyzer) checkEnglish(msg string) (string, bool) {
	re := regexp.MustCompile(`[ -~]`)
	for _, char := range msg {
		if match := re.MatchString(string(char)); !match && unicode.IsLetter(char) {
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

func (an *Analyzer) checkCustomRule(msg string, rule *config.Rule) (string, bool) {
	d, err := detector.NewDetector(
		detector.WithoutDefaultPatterns(),
		detector.WithPatterns(detector.Pattern{
			Name:        rule.Name,
			Description: rule.Description,
			Expression:  rule.Pattern,
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
