package config

type Rule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Pattern     string `json:"pattern"`
}

type Config struct {
	CapitalLetter     bool   `json:"capitalLetter"`
	EnglishLetter     bool   `json:"englishLetter"`
	SpecialLetters    bool   `json:"specialLetters"`
	SensitiveData     bool   `json:"sensitiveData"`
	EnableCustomRules bool   `json:"enableCustomRules"`
	Rules             []Rule `json:"rules"`
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
