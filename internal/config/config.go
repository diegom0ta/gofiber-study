package config

var jwtSecret string
var DBCfg string

func LoadDB() {
	DBCfg = "host=localhost user=postgres password=postgres dbname=server port=5433 sslmode=disable"
}

func LoadJwtSecret() string {
	jwtSecret = "secret"
	return jwtSecret
}
