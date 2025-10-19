package gokeenlog

import (
	"fmt"

	"github.com/kam1k88/gokeenapi/pkg/backends/keenetic/models"
	"github.com/kam1k88/gokeenapi/pkg/config"
)

func Info(msg string) {
	fmt.Println(msg)
}

func Infof(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	fmt.Printf("%v\n", s)
}

func InfoSubStepf(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	fmt.Printf("      ▪ %v\n", s)
}

func InfoSubStep(msg string) {
	fmt.Printf("      ▪ %v\n", msg)
}

func PrintParseResponse(parseResponse []models.ParseResponse) {
	if !config.Cfg.Logs.Debug {
		return
	}
	if len(parseResponse) == 0 {
		return
	}
	for _, parse := range parseResponse {
		for _, status := range parse.Parse.Status {
			InfoSubStep(status.Message)
		}
	}
}
