package envs

import (
	"log"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.AutomaticEnv() // allow viper automatically read os variables

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("error happened while config initialization: ", err)
	}
}

func Get(name string) string {
	value := viper.Get(name).(string)
	if value == "" {
		log.Fatalf("environment variable %semail-service.go is not set", name)
	}

	return value
}
