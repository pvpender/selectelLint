package capital

import "log/slog"

func TestSimple() {
	slog.Info("Test")       // want `Start uppercase detected!`
	slog.Debug("I am here") // want `Start uppercase detected!`
}

func Normal() {
	slog.Info("normal log")

}
