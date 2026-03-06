package custom

import "log/slog"

func TestCustomPassword() {
	slog.Warn("12314534059438438534935483454380348509340943i8093409") // want `Digit detected!`
}

func TestAuthor() {
	slog.Info("pvpender") // want `Author detected!`
}
