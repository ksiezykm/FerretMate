package config
type DatabaseConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Config struct {
	Databases []DatabaseConfig `json:"databases"`
}
