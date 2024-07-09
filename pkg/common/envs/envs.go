package envs

import (
	"log"
	"strconv"

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
		log.Fatalf("environment variable %semail-service.go is not set\n", name)
	}

	return value
}

func GetInt(name string) int {
	strValue := Get(name)
	value, err := strconv.Atoi(strValue)
	if err != nil {
		log.Fatalf("converting string value `%s` to int: %v", strValue, err)
	}

	return value
}
