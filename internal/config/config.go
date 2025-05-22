package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int
	}

	Divar struct {
		APIKey            string `mapstructure:"api_key"`
		BaseURL           string `mapstructure:"base_url"`
		MaxAPICallsPerDay int    `mapstructure:"max_api_calls_per_day"`
	}

	Recommendation struct {
		MinScoreThreshold  float64 `mapstructure:"min_score_threshold"`
		TopN               int     `mapstructure:"top_n"`
		FinalN             int     `mapstructure:"final_n"`
		ProductionYearHigh int     `mapstructure:"production_year_high"`
		ProductionYearLow  int     `mapstructure:"production_year_low"`
		UsageCoefficient   float32 `mapstructure:"usage_coefficient"`
	}

	Database struct {
		Type       string
		Host       string
		Port       int
		User       string
		Password   string
		DBName     string `mapstructure:"dbname"`
		Persistent bool
	}
}

var AppConfig Config

func LoadConfig(path string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error in reading config: %v", err)
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("error in unmarshall config: %v", err)
	}
}
