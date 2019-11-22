package config

type MailConfig struct {
	Host       string `json:"host"`
	Port       uint   `json:"port"`
	SenderName string `json:"sender_name"`
	Address    string `json:"address"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	UseTLS     bool   `json:"use_tls"`
}
