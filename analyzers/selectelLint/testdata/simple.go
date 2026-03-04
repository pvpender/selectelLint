package testdata

import "log/slog"

func TestSimple() {
	slog.Info("Test") // want `Start uppercase detected!`
}

func Normal() {
	slog.Info("normal log")
}
