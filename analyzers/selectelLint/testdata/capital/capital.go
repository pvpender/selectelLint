package capital

import "log/slog"

func TestSimple() {
	slog.Info("Test")       // want `Start uppercase detected!`
	slog.Debug("I am here") // want `Start uppercase detected!`
}

func Normal() {
	slog.Info("normal log")
}

func TestLogger() {
	logger := slog.Logger{}

	logger.Warn("Warn")   // want `Start uppercase detected!`
	logger.Debug("Debug") // want `Start uppercase detected!`
	logger.Info("Info")   // want `Start uppercase detected!`
	logger.Error("Error") // want `Start uppercase detected!`
	logger.Info("ok")
}
