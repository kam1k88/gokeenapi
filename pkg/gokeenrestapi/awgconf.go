package gokeenrestapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/gokeenspinner"
	"github.com/noksa/gokeenapi/pkg/models"
)

var AwgConf keeneticAwgconf

type keeneticAwgconf struct{}

func (*keeneticAwgconf) ConfigureOrUpdateInterface(interfaceId, jc, jmin, jmax, s1, s2, h1, h2, h3, h4 string) error {
	var parseSlice []models.ParseRequest
	confParse := models.ParseRequest{}
	confParse.Parse = fmt.Sprintf("interface %v wireguard asc %v %v %v %v %v %v %v %v %v",
		interfaceId,
		jc,
		jmin,
		jmax,
		s1,
		s2,
		h1,
		h2,
		h3,
		h4)
	parseSlice = append(parseSlice, confParse,
		models.ParseRequest{Parse: "system configuration save"})
	return gokeenspinner.WrapWithSpinner(fmt.Sprintf("Configuring %v interface with ASC parameters", color.CyanString(interfaceId)), func() error {
		_, err := ExecutePostParse(parseSlice...)
		return err
	})
}

func (*keeneticAwgconf) AddInterface(confFile string, name string) (CreatedInterface, error) {
	b, err := os.ReadFile(confFile)
	if err != nil {
		return CreatedInterface{}, err
	}
	if name == "" {
		name = filepath.Base(confFile)
	}

	importData := Import{Import: base64.StdEncoding.EncodeToString(b), Name: "", Filename: name}

	var createdInterface CreatedInterface
	err = gokeenspinner.WrapWithSpinner(fmt.Sprintf("Adding interface from the config file"), func() error {
		response, err := ExecutePostSubPath("/rci/interface/wireguard/import", importData)
		if err != nil {
			return err
		}
		err = json.Unmarshal(response, &createdInterface)
		return err
	})
	return createdInterface, err
}

type Import struct {
	Import   string `json:"import"`
	Name     string `json:"name"`
	Filename string `json:"filename"`
}

type CreatedInterface struct {
	Intersects string `json:"intersects"`
	Created    string `json:"created"`
	Status     []struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Ident   string `json:"ident"`
		Message string `json:"message"`
	} `json:"status"`
}
