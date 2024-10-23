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

func (*keeneticInterface) RciShowInterfaces(withType string) (map[string]models.RciShowInterface, error) {
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
	for k, interfaceDetails := range interfaces {
		if withType != "" && strings.ToLower(interfaceDetails.Type) != strings.ToLower(withType) {
			continue
		}
		keenlog.Infof("Interface '%v':\n", k)
		keenlog.Infof("  * Id: %v\n", interfaceDetails.Id)
		keenlog.Infof("  * Type: %v\n", interfaceDetails.Type)
		if interfaceDetails.Description != "" {
			keenlog.Infof("  * Description: %v\n", interfaceDetails.Description)
		}
		if interfaceDetails.Address != "" {
			keenlog.Infof("  * Address: %v\n", interfaceDetails.Address)
		}
		keenlog.Infof("\n")
	}
	return interfaces, nil
}
