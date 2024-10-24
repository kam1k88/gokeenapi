package keeneticapi

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/keenlog"
	"github.com/noksa/gokeenapi/internal/keenspinner"
	"github.com/noksa/gokeenapi/pkg/models"
	"strings"
)

type keeneticInterface struct {
}

var Interface keeneticInterface

func (*keeneticInterface) GetInterfacesViaRciShowInterfaces(interfaceTypes ...string) (map[string]models.RciShowInterface, error) {
	var interfaces map[string]models.RciShowInterface
	err := keenspinner.WrapWithSpinner("Fetching interfaces", func() error {
		body, err := ExecuteGetSubPath("/rci/show/interface")
		if err != nil {
			return err
		}
		return json.Unmarshal(body, &interfaces)
	})
	if err != nil {
		return interfaces, err
	}
	var interfaceTypesLower []string
	for _, v := range interfaceTypes {
		v := v
		interfaceTypesLower = append(interfaceTypesLower, strings.ToLower(v))
	}
	return interfaces, nil
}

func (*keeneticInterface) PrintInfoAboutInterfaces(interfaces map[string]models.RciShowInterface) {
	for k, interfaceDetails := range interfaces {
		keenlog.Infof("Interface '%v':", color.BlueString(k))
		keenlog.InfoSubStepf("Id: %v", color.CyanString(interfaceDetails.Id))
		keenlog.InfoSubStepf("Type: %v", color.CyanString(interfaceDetails.Type))
		if interfaceDetails.Description != "" {
			keenlog.InfoSubStepf("Description: %v", color.CyanString(interfaceDetails.Description))
		}
		if interfaceDetails.Address != "" {
			keenlog.InfoSubStepf("Address: %v", color.CyanString(interfaceDetails.Address))
		}
		keenlog.Infof("")
	}

}
