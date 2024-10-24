package keeneticapi

import (
	"encoding/json"
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
		body, err := ExecuteGetSubPath("/show/interface")
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
		keenlog.Infof("Interface '%v':", k)
		keenlog.InfoSubStepf("Id: %v", interfaceDetails.Id)
		keenlog.InfoSubStepf("Type: %v", interfaceDetails.Type)
		if interfaceDetails.Description != "" {
			keenlog.InfoSubStepf("Description: %v", interfaceDetails.Description)
		}
		if interfaceDetails.Address != "" {
			keenlog.InfoSubStepf("Address: %v", interfaceDetails.Address)
		}
		keenlog.Infof("")
	}

}
