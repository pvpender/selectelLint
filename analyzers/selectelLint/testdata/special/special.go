package special

import (
	"log/slog"
)

func Emoji() {
	slog.Info("server started 🚀") // want `Special letters detected!`
}

func SpecialLetter() {
	slog.Info("server started!") // want `Special letters detected!`
}
