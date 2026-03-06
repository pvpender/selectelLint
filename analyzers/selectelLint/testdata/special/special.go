package special

import (
	"log/slog"
)

func Emoji() {
	slog.Info("server started 🚀") // want `Special letters detected!`
	slog.Error("🤬")               // want `Special letters detected!`
	slog.Debug("🤖")               // want `Special letters detected!`
	slog.Warn("☂️")               // want `Special letters detected!`
}

func SpecialLetter() {
	slog.Info("server started!")   // want `Special letters detected!`
	slog.Warn("%asda")             // want `Special letters detected!`
	slog.Error("2543^4223=123412") // want `Special letters detected!`
	slog.Debug("/**/")             // want `Special letters detected!`
}

func InLogger() {
	logger := slog.Logger{}
	logger.Warn("this is incorrect...") // want `Special letters detected!`
}

func NoSpecial() {
	slog.Log(nil, slog.LevelDebug, "server started at 8080")
}
