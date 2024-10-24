package keeneticapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/noksa/gokeenapi/internal/config"
	"github.com/noksa/gokeenapi/pkg/models"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

var restyClient *resty.Client

func ExecutePostParse(parse ...models.ParseRequest) ([]models.ParseResponse, error) {
	request := GetApiClient().R()
	request.SetBody(parse)
	response, err := request.Post("/")
	var parseResponse []models.ParseResponse
	var mErr error
	if response != nil {
		decodeErr := json.Unmarshal(response.Body(), &parseResponse)
		mErr = multierr.Append(mErr, decodeErr)
		for _, myParse := range parseResponse {
			for _, status := range myParse.Parse.Status {
				if status.Status == "error" {
					mErr = multierr.Append(mErr, fmt.Errorf("%s - %s - %s - %s", status.Status, status.Code, status.Ident, status.Message))
				}
			}
		}
	}
	mErr = multierr.Append(mErr, err)
	return parseResponse, mErr
}

func ExecuteGetSubPath(path string) ([]byte, error) {
	response, err := GetApiClient().R().Get(path)
	if err != nil {
		return nil, err
	}
	if response != nil {
		return response.Body(), nil
	}
	return []byte{}, errors.New("no response from keenetic api")
}

func GetApiClient() *resty.Client {
	if restyClient != nil {
		return restyClient
	}
	restyClient = resty.New()
	restyClient.SetBasicAuth(viper.GetString(config.ViperKeeneticLogin), viper.GetString(config.ViperKeeneticPassword))
	restyClient.SetBaseURL(viper.GetString(config.ViperKeeneticApi))
	return restyClient
}
