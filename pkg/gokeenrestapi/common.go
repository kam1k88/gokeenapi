package gokeenrestapi

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/internal/gokeenspinner"
	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeencache"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

var (
	restyClient *resty.Client
	cookie      string
	Common      keeneticCommon
)

type keeneticCommon struct {
}

func (c *keeneticCommon) Auth() error {
	err := gokeenspinner.WrapWithSpinner("Authorizing in Keenetic", func() error {
		response, err := c.GetApiClient().R().Get("/auth")
		var mErr error
		if response != nil {
			if response.StatusCode() == http.StatusUnauthorized {
				realm := response.Header().Get("x-ndm-realm")
				token := response.Header().Get("x-ndm-challenge")
				setCookieStr := response.Header().Get("set-cookie")
				setCookieStrSplitted := strings.Split(setCookieStr, ";")
				cookie = setCookieStrSplitted[0]
				secondRequest := c.GetApiClient().R()
				//secondRequest.Header.Set("Cookie", cookie)

				md5Hash := md5.New()
				_, err = fmt.Fprintf(md5Hash, "%v:%v:%v", viper.GetString(config.ViperKeeneticLogin), realm, viper.GetString(config.ViperKeeneticPassword))
				if err != nil {
					return err
				}
				md5HashArg := md5Hash.Sum(nil)
				md5HashStr := hex.EncodeToString(md5HashArg)
				sha256Hash := sha256.New()
				_, err = fmt.Fprintf(sha256Hash, "%v%v", token, md5HashStr)
				if err != nil {
					return err
				}
				sha256HashArg := sha256Hash.Sum(nil)
				sha256HashStr := hex.EncodeToString(sha256HashArg)
				secondRequest.SetBody(struct {
					Login    string `json:"login"`
					Password string `json:"password"`
				}{
					Login:    viper.GetString(config.ViperKeeneticLogin),
					Password: sha256HashStr,
				})
				response, err = secondRequest.Post("/auth")
				if err != nil {
					return err
				}
				if response.StatusCode() == http.StatusUnauthorized {
					return errors.New("can't authorize in keenetic. Verify your login and password")
				}
			}
		}
		mErr = multierr.Append(mErr, err)
		return mErr
	})
	if err == nil {
		version, err := c.Version()
		if err != nil {
			return err
		}
		gokeenlog.InfoSubStepf("Router: %v", color.CyanString(version.Model))
		gokeenlog.InfoSubStepf("OS version: %v", color.CyanString(version.Title))
		gokeencache.SetVersion(version)
	}
	return err
}

func (c *keeneticCommon) Version() (gokeenrestapimodels.Version, error) {
	b, err := c.ExecuteGetSubPath("/rci/show/version")
	if err != nil {
		return gokeenrestapimodels.Version{}, err
	}
	var version gokeenrestapimodels.Version
	err = json.Unmarshal(b, &version)
	return version, err
}

func (c *keeneticCommon) ExecutePostParse(parse ...gokeenrestapimodels.ParseRequest) ([]gokeenrestapimodels.ParseResponse, error) {
	parseCopy := parse
	var parseResponse []gokeenrestapimodels.ParseResponse
	var mErr error
	for len(parseCopy) > 0 {
		request := c.GetApiClient().R()
		maxParse := viper.GetInt("keenetic.routesPerRequest")
		if maxParse == 0 {
			maxParse = 50
		} else if maxParse < 20 {
			maxParse = 20
		}
		currentLen := len(parseCopy)
		if currentLen < maxParse {
			maxParse = currentLen
		}
		var parseRequest []gokeenrestapimodels.ParseRequest
		for i := 0; i < maxParse; i++ {
			parseRequest = append(parseRequest, parseCopy[i])
		}
		parseCopy = parseCopy[maxParse:]
		request.SetBody(parseRequest)
		response, err := request.Post("/rci/")
		if response != nil {
			if response.StatusCode() != http.StatusOK {
				mErr = multierr.Append(mErr, fmt.Errorf("wrong status code in response from api: %s", response.Status()))
			}
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
	}
	return parseResponse, mErr
}

func (c *keeneticCommon) ExecuteGetSubPath(path string) ([]byte, error) {
	response, err := c.GetApiClient().R().Get(path)
	if err != nil {
		return nil, err
	}
	if response != nil {
		return response.Body(), nil
	}
	return []byte{}, errors.New("no response from keenetic api")
}

func (c *keeneticCommon) ExecutePostSubPath(path string, body any) ([]byte, error) {
	response, err := c.GetApiClient().R().SetBody(body).Post(path)
	if err != nil {
		return nil, err
	}
	if response != nil {
		return response.Body(), nil
	}
	return []byte{}, errors.New("no response from keenetic api")
}

func (*keeneticCommon) GetApiClient() *resty.Client {
	if restyClient != nil {
		if restyClient.Header.Get("Cookie") == "" && cookie != "" {
			restyClient.Header.Set("Cookie", cookie)
		}
		return restyClient
	}
	restyClient = resty.New()
	restyClient.SetDisableWarn(true)
	restyClient.SetBaseURL(viper.GetString(config.ViperKeeneticUrl))
	return restyClient
}
