package selectelLint

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"log/slog"
	"pvpender/selectelLint/config"
	"unicode"
)

type Analyzer struct {
	config *config.Config
}

func NewAnalyzer() *analysis.Analyzer {
	ca := &Analyzer{
		config: config.NewConfig(),
	}

	return &analysis.Analyzer{
		Name:     "CapitalAnalyzer",
		Doc:      "Checks that there are no capitals in logs",
		Flags:    ca.Flags(),
		Run:      ca.Run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func (ca *Analyzer) Flags() flag.FlagSet {
	flags := flag.NewFlagSet("CapitalAnalyzer", flag.ExitOnError)
	flags.BoolVar(&ca.config.CapitalLetter, "capitalLetter", true, "Enable capital letter check")
	flags.BoolVar(&ca.config.EnglishLetter, "englishLetter", true, "Enable english letter check")
	flags.BoolVar(&ca.config.SpecialLetters, "specialLetter", true, "Enable special letter check")
	flags.BoolVar(&ca.config.SensitiveData, "sensitiveData", true, "Enable sensitive data")
	flags.BoolVar(&ca.config.EnableCustomRules, "enableCustomRules", false, "Enable custom rules")

	return *flags
}

func (ca *Analyzer) Run(pass *analysis.Pass) (interface{}, error) {
	inspectorVar := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspectorVar.Preorder(nodeFilter, func(node ast.Node) {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return
		}

		if !ca.isLogFunction(call) {
			return
		}

		message, exists := ca.getMessage(call)

		if !exists {
			return
		}

		if ca.config.CapitalLetter {
			if msg, fail := ca.checkCapital(message); fail {
				pass.Reportf(call.Pos(), "%s", msg)
			}
		}

		if ca.config.EnglishLetter {
			ca.checkEnglish(message)
		}

		if ca.config.SpecialLetters {
			ca.checkSpecialLetters(message)
		}

		if ca.config.EnableCustomRules {
			ca.checkSensitiveData(message, call)
		}

	})

	return nil, nil
}

func (ca *Analyzer) isLogFunction(call *ast.CallExpr) bool {
	logMethods := map[string]bool{
		"Info": true,
	}

	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		return logMethods[fun.Sel.Name]

	case *ast.Ident:
		slog.Info(fmt.Sprintf("%t", call.Fun))

		return logMethods[fun.Name]
	}

	return false
}

// Возвращаемый тип у функции проверить
func (ca *Analyzer) getMessage(expr ast.Expr) (string, bool) {
	switch lit := expr.(type) {
	case *ast.CallExpr:
		if len(lit.Args) == 0 {
			return "", false
		}

		for _, arg := range lit.Args {
			if msg, ok := ca.getMessage(arg); ok {
				return msg, true
			}
		}

		break
	case *ast.BasicLit:
		if lit.Kind == token.STRING {
			return lit.Value, true
		}

		break

	case *ast.BinaryExpr:
		if lit.Op == token.AND {
			left, leftOk := ca.getMessage(lit.X)
			right, rightOk := ca.getMessage(lit.Y)

			if leftOk && rightOk {
				return left + right, true
			}
		}

		break
	}

	return "", false
}

func (ca *Analyzer) checkCapital(msg string) (string, bool) {
	if unicode.IsUpper([]rune(msg)[1]) && unicode.IsLetter([]rune(msg)[1]) {
		return "Start uppercase detected!", true
	}

	return "", false
}

func (ca *Analyzer) checkEnglish(msg string) {

}

func (ca *Analyzer) checkSpecialLetters(msg string) {

}

func (ca *Analyzer) checkSensitiveData(msg string, call *ast.CallExpr) {

}
