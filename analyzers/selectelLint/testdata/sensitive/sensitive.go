package sensitive

import "log/slog"

const PASSWORD = "SuperSecretPassword18234!"

const API_KEY = "SJLDSF7798ASAAABDSFOT1231"

const TOKEN = "AAJJJGI1648349087IFSSD0"

func Password() {
	slog.Info("user password=" + PASSWORD) // want `Value that looks like a password, secret, or API key assignment`
}

func APIKey() {
	slog.Info("password:" + API_KEY) // want `Value that looks like a password, secret, or API key assignment`
}

func Token() {
	slog.Info("token:" + TOKEN) // want `Value that looks like a password, secret, or API key assignment`
}

func SprintToken() {
	token := TOKEN
	slog.Info("token:" + token) // want `Value that looks like a password, secret, or API key assignment`
}
