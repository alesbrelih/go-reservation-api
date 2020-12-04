package config

import "time"

type Enviroment struct {
	Jwt struct {
		Secret            string        `env:"JWT_SECRET,required=true"`
		AccessExpiration  time.Duration `env:"JWT_ACCESS_EXPIRATION,required=true"`
		RefreshExpiration time.Duration `env:"JWT_REFRESH_EXPIRATION,required=true"`
	}
	Application struct {
		PORT string `env:"APPLICATION_PORT,required=true"`
	}
	Database struct {
		URL string `env:"POSTGRES_URL,required=true"`
	}
}
6