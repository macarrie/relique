package config

type HTTPConfig struct {
	BindAddr string `mapstructure:"bind_addr" json:"bind_addr" toml:"bind_addr"`
	Port     int    `mapstructure:"port" json:"port" toml:"port"`
	SSLCert  string `mapstructure:"ssl_cert" json:"ssl_cert" toml:"ssl_cert"`
	SSLKey   string `mapstructure:"ssl_key" json:"ssl_key" toml:"ssl_key"`
}
