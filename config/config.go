package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

var ViperConfig *viper.Viper

var Hostname string

func init() {
	hostAddress, err := os.Hostname()
	if err != nil {
		fmt.Println("Hostname of the current server is not set.")
		panic(err)
	}

	re := regexp.MustCompile(`^(\\w+).*$`)
	arr := strings.Split(hostAddress, ".")
	Hostname = re.FindString(arr[0])

}

func ReadConfig(environment string) error {
	ViperConfig = viper.New()
	ViperConfig.SetConfigName("VelocityApiConfig")
	ViperConfig.SetConfigType("yaml")

	if environment == "development" {
		ViperConfig.SetConfigFile("env/dev.yml")
	} else if environment == "staging" {
		ViperConfig.SetConfigFile("env/stage.yml")
	} else {
		ViperConfig.SetConfigFile("env/prod.yml")
	}
	err := ViperConfig.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
