package config

type Rule struct {
}

type Config struct {
	CapitalLetter     bool
	EnglishLetter     bool
	SpecialLetters    bool
	SensitiveData     bool
	EnableCustomRules bool
	Rules             []Rule
}

func NewConfig() *Config {
	return &Config{
		CapitalLetter:     true,
		EnglishLetter:     true,
		SpecialLetters:    true,
		SensitiveData:     true,
		EnableCustomRules: false,
		Rules:             nil,
	}
}
