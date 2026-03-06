package english

import "log/slog"

func TestEnglish() {
	slog.Info("english test success")
}

func NotEnglish() {
	slog.Info("не английский") // want `Not english character detected!`
}

func German() {
	slog.Warn("deutsch ähnelt Englisch") // want `Not english character detected!`
}

func Chinese() {
	slog.Warn("中国人") // want `Not english character detected!`
}
