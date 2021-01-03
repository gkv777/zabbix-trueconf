package main

// Config ...
type Config struct {
	TrueConfURL      string `yaml:"trueconf_url"`
	TrueConfClient   string `yaml:"trueconf_client"`
	TrueConfSecret   string `yaml:"trueconf_secret"`
	TrueConfTLSInsec bool   `yaml:"trueconf_tls_insecure"`
	UserAgent        string `yaml:"user_agent"`
	LogFile          string `yaml:"log_file"`
	LogDebug         bool   `yaml:"log_debug"`
}

