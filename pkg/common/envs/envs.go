package envs

import (
	"log"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error happened while config initialization: ", err)
	}
}
