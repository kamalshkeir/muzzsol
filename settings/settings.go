package settings

// Config will handle global application config
var Config = &config{}

type config struct {
	Port string `kenv:"PORT|80"`
	Secret string `kenv:"SECRET|"`
	DB struct {
		Type string `kenv:"DB_TYPE|mysql"`
		Name string `kenv:"DB_NAME|"`
		Dsn string `kenv:"DB_DSN|"`
	}
}