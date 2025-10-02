package gokeenrestapi

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/noksa/gokeenapi/internal/gokeencache"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/internal/gokeenspinner"
	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"go.uber.org/multierr"
)

var (
	restyClient          *resty.Client
	cleanedOldCacheFiles bool
	Common               keeneticCommon
)

type keeneticCommon struct {
}

type keeneticCacheFile struct {
	Cookie keeneticCacheCookie `json:"cookie,omitempty"`
	path   string
}
type keeneticCacheCookie struct {
	Value      string    `json:"value"`
	UpdateTime time.Time `json:"update_time"`
}

func (f *keeneticCacheFile) Save() error {
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(f.path, b, 0700)
	return err
}

func (c *keeneticCommon) getKeeneticCacheFile() (keeneticCacheFile, error) {
	var dataDir string
	var err error
	if config.Cfg.DataDir != "" {
		dataDir = path.Clean(config.Cfg.DataDir)
	} else {
		dataDir, err = os.UserHomeDir()
		if err != nil {
			return keeneticCacheFile{}, err
		}
	}
	gokeenDir := path.Join(dataDir, ".gokeenapi")
	err = os.MkdirAll(gokeenDir, os.ModePerm)
	if err != nil {
		return keeneticCacheFile{}, err
	}
	if !cleanedOldCacheFiles {
		err = filepath.WalkDir(gokeenDir, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return err
			}
			if time.Since(info.ModTime()) >= time.Hour*24*7 {
				err = os.Remove(path)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return keeneticCacheFile{}, err
		}
		cleanedOldCacheFiles = true
	}
	bHash := md5.Sum([]byte(fmt.Sprintf("%v-%v-%v", config.Cfg.Keenetic.URL, config.Cfg.Keenetic.Login, config.Cfg.Keenetic.Password)))
	hash := hex.EncodeToString(bHash[:])
	keeeticFile := path.Join(gokeenDir, fmt.Sprintf("%v.json", hash))
	_, statErr := os.Stat(keeeticFile)
	if statErr != nil {
		if !errors.Is(statErr, os.ErrNotExist) {
			return keeneticCacheFile{}, statErr
		}
		err = os.WriteFile(keeeticFile, []byte("{}"), os.ModePerm)
		if err != nil {
			return keeneticCacheFile{}, err
		}
	}
	var keeneticCache keeneticCacheFile
	b, err := os.ReadFile(keeeticFile)
	if err != nil {
		return keeneticCacheFile{}, err
	}
	err = json.Unmarshal(b, &keeneticCache)
	keeneticCache.path = keeeticFile
	return keeneticCache, err
}

func (c *keeneticCommon) getAuthCookie() (string, error) {
	cache, err := c.getKeeneticCacheFile()
	if err != nil {
		return "", err
	}
	return cache.Cookie.Value, nil
}

func (c *keeneticCommon) writeAuthCookie(cookie string) error {
	cache, err := c.getKeeneticCacheFile()
	if err != nil {
		return err
	}
	cache.Cookie.Value = cookie
	cache.Cookie.UpdateTime = time.Now()
	return cache.Save()
}

func (c *keeneticCommon) Auth() error {
	err := gokeenspinner.WrapWithSpinner(fmt.Sprintf("Authorizing in %v", color.CyanString("Keenetic")), func() error {
		response, err := c.GetApiClient().R().Get("/auth")
		var mErr error
		mErr = multierr.Append(mErr, err)
		if response != nil && response.StatusCode() == http.StatusUnauthorized {
			realm := response.Header().Get("x-ndm-realm")
			token := response.Header().Get("x-ndm-challenge")
			setCookieStr := response.Header().Get("set-cookie")
			setCookieStrSplitted := strings.Split(setCookieStr, ";")
			cookieToSet := setCookieStrSplitted[0]
			err = c.writeAuthCookie(cookieToSet)
			if err != nil {
				mErr = multierr.Append(mErr, err)
				return mErr
			}
			secondRequest := c.GetApiClient().R()
			// cookie should not be set here if we do authorization. It means that old cookie doesn't work anymore and we would get 401 with the old.
			restyClient.Header.Del("Cookie")
			secondRequest.Header.Del("Cookie")
			md5Hash := md5.New()
			_, err = fmt.Fprintf(md5Hash, "%v:%v:%v", config.Cfg.Keenetic.Login, realm, config.Cfg.Keenetic.Password)
			if err != nil {
				mErr = multierr.Append(mErr, err)
				return mErr
			}
			md5HashArg := md5Hash.Sum(nil)
			md5HashStr := hex.EncodeToString(md5HashArg)
			sha256Hash := sha256.New()
			_, err = fmt.Fprintf(sha256Hash, "%v%v", token, md5HashStr)
			if err != nil {
				mErr = multierr.Append(mErr, err)
				return mErr
			}
			sha256HashArg := sha256Hash.Sum(nil)
			sha256HashStr := hex.EncodeToString(sha256HashArg)
			secondRequest.SetBody(struct {
				Login    string `json:"login"`
				Password string `json:"password"`
			}{
				Login:    config.Cfg.Keenetic.Login,
				Password: sha256HashStr,
			})
			response, err = secondRequest.Post("/auth")
			if err != nil {
				mErr = multierr.Append(mErr, err)
				return mErr
			}
			if response.StatusCode() == http.StatusUnauthorized {
				mErr = multierr.Append(mErr, errors.New("can't authorize in keenetic. Verify your login and password"))
				return mErr
			}
		}
		return mErr
	})
	if err != nil {
		return err
	}
	version, err := c.Version()
	if err != nil {
		return err
	}
	gokeenlog.InfoSubStepf("%v: %v", color.BlueString("Router"), color.CyanString(version.Model))
	gokeenlog.InfoSubStepf("%v: %v", color.BlueString("OS version"), color.CyanString(version.Title))
	gokeencache.UpdateRuntimeConfig(func(runtime *config.Runtime) {
		runtime.RouterInfo.Version = version
	})
	return nil
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
	var parseResponses []gokeenrestapimodels.ParseResponse
	var mErr error
	for len(parseCopy) > 0 {
		request := c.GetApiClient().R()
		maxParse := 50
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
			var parseResponse []gokeenrestapimodels.ParseResponse
			decodeErr := json.Unmarshal(response.Body(), &parseResponse)
			mErr = multierr.Append(mErr, decodeErr)
			for i, myParse := range parseResponse {
				if i == 0 {
					parseResponse[i].Parse.DynamicData = string(response.Body())
				}
				for _, status := range myParse.Parse.Status {
					if status.Status == "error" {
						mErr = multierr.Append(mErr, fmt.Errorf("%s - %s - %s - %s", status.Status, status.Code, status.Ident, status.Message))
					}
				}
			}
			parseResponses = append(parseResponses, parseResponse...)
		}
		mErr = multierr.Append(mErr, err)
	}
	return parseResponses, mErr
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

func (c *keeneticCommon) GetApiClient() *resty.Client {
	if restyClient == nil {
		restyClient = resty.New()
		restyClient.SetDisableWarn(true)
	}
	// do it each time in case of GUI version
	restyClient.SetBaseURL(config.Cfg.Keenetic.URL)
	if restyClient.Header.Get("Cookie") == "" {
		cookie, err := c.getAuthCookie()
		if err != nil {
			panic(err)
		}
		if cookie != "" {
			restyClient.Header.Set("Cookie", cookie)
		}
	}
	return restyClient
}

func (c *keeneticCommon) ShowRunningConfig() (gokeenrestapimodels.RunningConfig, error) {
	var runningConfig gokeenrestapimodels.RunningConfig
	err := gokeenspinner.WrapWithSpinner(fmt.Sprintf("Fetching %v", color.CyanString("running-config")), func() error {
		b, err := c.ExecuteGetSubPath("/rci/show/running-config")
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, &runningConfig)
		return err
	})
	return runningConfig, err
}
