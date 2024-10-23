package keenlog

import (
	"fmt"
	"github.com/noksa/gokeenapi/internal/config"
	"github.com/noksa/gokeenapi/pkg/models"
	"github.com/spf13/viper"
)

func Info(msg string) {
	fmt.Println(msg)
}

func Infof(msg string, args ...any) {
	fmt.Printf(msg, args...)
}

func Panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

func Debug(msg string) {
	if viper.GetBool(config.ViperDebug) {
		fmt.Println(msg)
	}
}

func PrintParseResponse(parseResponse []models.ParseResponse) {
	if !viper.GetBool(config.ViperDebug) {
		return
	}
	if len(parseResponse) == 0 {
		return
	}
	Info("Result:")
	for _, parse := range parseResponse {
		for _, status := range parse.Parse.Status {
			Infof("  ▪ %v\n", status.Message)
		}
	}
}
