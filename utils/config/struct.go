package config

// Databases

type Postgres struct {
	User     string `env:"POSTGRES_USER"`
	Pass     string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DATABASE"`
	IP       string `env:"POSTGRES_IP"`
	Port     string `env:"POSTGRES_PORT"`
	Protocol string `env:"POSTGRES_PROTOCOL"`
}

type Memcached struct {
	Address string `env:"MEMCACHED_ADDRESS"`
}

type Redis struct {
	Address string `env:"REDIS_ADDRESS"`
	Pass    string `env:"REDIS_PASS"`
	DB      int    `env:"REDIS_DB"`
}

// Oauth

type Oauth struct {
	GitPath    string `env:"OAUTH_GITHUB_PATH"`
	GooglePath string `env:"OAUTH_GOOGLE_PATH"`
}

// Host

type Host struct {
	Port string `env:"HOST_PORT"`
	Key  string `env:"HOST_KEY_PATH"`  // Path to TLS key
	Cert string `env:"HOST_CERT_PATH"` // Path to TLS certificate
	//
	Templates string `env:"HOST_TEMPLATES"`
	//
	Static string `env:"HOST_STATIC"`
}

// Cookie

type Cookie struct {
	Key  string `env:"COOKIE_KEY"`
	Name string `env:"COOKIE_NAME"`
	Part string `env:"COOKIE_PART"`
}

//

type SMTP struct {
	Mail     string `env:"SMTP_MAIL"`
	Password string `env:"SMTP_PASSWORD"`

	Hostname string `env:"SMTP_HOSTNAME"`
	Port     int    `env:"SMTP_PORT"`
}

//

type Config struct {
	Postgres  Postgres
	Redis     Redis
	Memcached Memcached

	Oauth Oauth
	SMTP  SMTP

	Host Host

	Cookie Cookie
}

func (c *Config) LoadEnv() error {
	return readEnv(c)
}
func (c *Config) Verify() error {
	return verify(c)
}
