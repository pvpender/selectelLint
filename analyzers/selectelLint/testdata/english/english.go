package english

import "log/slog"

func TestEnglish() {
	slog.Info("english test success")
}

func NotEnglish() {
	slog.Info("не английский") // want `Not english character detected!`
}
