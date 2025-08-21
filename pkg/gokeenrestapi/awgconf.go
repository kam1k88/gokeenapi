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
	"go.uber.org/multierr"
	"gopkg.in/ini.v1"
)

var AwgConf keeneticAwgconf

type keeneticAwgconf struct{}

// ConfigureOrUpdateInterface checks if ACS parameters should be adjusted in the interfaceId
// Starting from 4.3.6 KeeneticOS automatically adds ACS from AWG conf file when it is imported
func (*keeneticAwgconf) ConfigureOrUpdateInterface(confPath, interfaceId string) error {
	if confPath == "" {
		return fmt.Errorf("conf-file flag is required")
	}
	err := Checks.CheckInterfaceId(interfaceId)
	if err != nil {
		return err
	}
	err = Checks.CheckInterfaceExists(interfaceId)
	if err != nil {
		return err
	}
	var Jcstring, Jminstring, Jmaxstring, S1string, S2string, H1string, H2string, H3string, H4string string
	confPath, err = filepath.Abs(confPath)
	if err != nil {
		return err
	}
	cfg, err := ini.Load(confPath)
	if err != nil {
		return err
	}
	interfaceSection, err := cfg.GetSection("Interface")
	if err != nil {
		return err
	}
	Jc, err := interfaceSection.GetKey("Jc")
	if err != nil {
		return err
	}
	Jcstring = Jc.String()
	Jmin, err := interfaceSection.GetKey("Jmin")
	if err != nil {
		return err
	}
	Jminstring = Jmin.String()
	Jmax, err := interfaceSection.GetKey("Jmax")
	if err != nil {
		return err
	}
	Jmaxstring = Jmax.String()
	S1, err := interfaceSection.GetKey("S1")
	if err != nil {
		return err
	}
	S1string = S1.String()
	S2, err := interfaceSection.GetKey("S2")
	if err != nil {
		return err
	}
	S2string = S2.String()
	H1, err := interfaceSection.GetKey("H1")
	if err != nil {
		return err
	}
	H1string = H1.String()
	H2, err := interfaceSection.GetKey("H2")
	if err != nil {
		return err
	}
	H2string = H2.String()
	H3, err := interfaceSection.GetKey("H3")
	if err != nil {
		return err
	}
	H3string = H3.String()
	H4, err := interfaceSection.GetKey("H4")
	if err != nil {
		return err
	}
	H4string = H4.String()

	interfaceDetails, err := Interface.GetInterfaceViaRciShowScInterfaces(interfaceId)
	if err != nil {
		return err
	}
	shouldApply := false

	asc := interfaceDetails.Wireguard.Asc

	if asc.Jc != Jcstring {
		shouldApply = true
	}
	if asc.Jmin != Jminstring {
		shouldApply = true
	}
	if asc.Jmax != Jmaxstring {
		shouldApply = true
	}
	if asc.S1 != S1string {
		shouldApply = true
	}
	if asc.S2 != S2string {
		shouldApply = true
	}
	if asc.H1 != H1string {
		shouldApply = true
	}
	if asc.H2 != H2string {
		shouldApply = true
	}
	if asc.H3 != H3string {
		shouldApply = true
	}
	if asc.H4 != H4string {
		shouldApply = true
	}
	if !shouldApply {
		return nil
	}

	var parseSlice []models.ParseRequest
	confParse := models.ParseRequest{}
	confParse.Parse = fmt.Sprintf("interface %v wireguard asc %v %v %v %v %v %v %v %v %v",
		interfaceId,
		Jcstring,
		Jminstring,
		Jmaxstring,
		S1string,
		S2string,
		H1string,
		H2string,
		H3string,
		H4string)
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
		for _, status := range createdInterface.Status {
			if status.Status == "error" {
				err = multierr.Append(err, fmt.Errorf("%v - %v - %v - %v", status.Status, status.Code, status.Ident, status.Message))
			}
		}
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
