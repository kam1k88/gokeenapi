package gokeenrestapi

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/internal/gokeenspinner"
	"github.com/noksa/gokeenapi/pkg/models"
)

type keeneticInterface struct {
}

var Interface keeneticInterface

func (*keeneticInterface) GetInterfaceViaRciShowInterfaces(interfaceId string) (models.RciShowInterface, error) {
	var myInterface models.RciShowInterface
	body, err := ExecuteGetSubPath(fmt.Sprintf("/rci/show/interface/%v", interfaceId))
	if err != nil {
		return myInterface, err
	}
	err = json.Unmarshal(body, &myInterface)
	return myInterface, err
}

func (*keeneticInterface) GetInterfaceViaRciShowScInterfaces(interfaceId string) (models.RciShowScInterface, error) {
	var myInterface models.RciShowScInterface
	body, err := ExecuteGetSubPath(fmt.Sprintf("/rci/show/sc/interface/%v", interfaceId))
	if err != nil {
		return myInterface, err
	}
	err = json.Unmarshal(body, &myInterface)
	return myInterface, err
}

func (*keeneticInterface) GetInterfacesViaRciShowInterfaces(interfaceTypes ...string) (map[string]models.RciShowInterface, error) {
	var interfaces map[string]models.RciShowInterface
	err := gokeenspinner.WrapWithSpinner("Fetching interfaces", func() error {
		body, err := ExecuteGetSubPath("/rci/show/interface")
		if err != nil {
			return err
		}
		return json.Unmarshal(body, &interfaces)
	})
	if err != nil {
		return interfaces, err
	}
	if len(interfaceTypes) == 0 {
		return interfaces, nil
	}
	realInterfaces := map[string]models.RciShowInterface{}
	for k, interfaceDetails := range interfaces {
		k := k
		interfaceDetails := interfaceDetails
		for _, v := range interfaceTypes {
			v := v
			if strings.EqualFold(interfaceDetails.Type, v) {
				realInterfaces[k] = interfaceDetails
			}
		}
	}
	return realInterfaces, nil
}

func (*keeneticInterface) GetInterfacesViaRciShowScInterfaces(ids ...string) (map[string]models.RciShowScInterface, error) {
	var interfaces map[string]models.RciShowScInterface
	err := gokeenspinner.WrapWithSpinner("Fetching interfaces", func() error {
		body, err := ExecuteGetSubPath("/rci/show/sc/interface")
		if err != nil {
			return err
		}
		return json.Unmarshal(body, &interfaces)
	})
	if err != nil {
		return interfaces, err
	}
	if len(ids) == 0 {
		return interfaces, nil
	}
	realInterfaces := map[string]models.RciShowScInterface{}
	for k, interfaceDetails := range interfaces {
		if !slices.Contains(ids, k) {
			continue
		}
		realInterfaces[k] = interfaceDetails
	}
	return realInterfaces, nil
}

func (*keeneticInterface) PrintInfoAboutInterfaces(interfaces map[string]models.RciShowInterface) {
	for k, interfaceDetails := range interfaces {
		gokeenlog.Infof("Interface '%v':", color.BlueString(k))
		gokeenlog.InfoSubStepf("Id: %v", color.CyanString(interfaceDetails.Id))
		gokeenlog.InfoSubStepf("Type: %v", color.CyanString(interfaceDetails.Type))
		if interfaceDetails.Description != "" {
			gokeenlog.InfoSubStepf("Description: %v", color.CyanString(interfaceDetails.Description))
		}
		if interfaceDetails.Address != "" {
			gokeenlog.InfoSubStepf("Address: %v", color.CyanString(interfaceDetails.Address))
		}
		gokeenlog.Infof("")
	}

}

func (*keeneticInterface) WaitUntilInterfaceIsUp(interfaceId string) error {
	err := gokeenspinner.WrapWithSpinner(fmt.Sprintf("Waiting 60s until %v interface is up, connected to peers and working", interfaceId), func() error {
		deadline := time.Now().Add(time.Second * 60)
		for time.Now().Before(deadline) {
			myInterface, err := Interface.GetInterfaceViaRciShowInterfaces(interfaceId)
			if err != nil {
				return err
			}
			if myInterface.Connected == "yes" && myInterface.Link == "up" && myInterface.State == "up" {
				return nil
			}
			time.Sleep(time.Millisecond * 500)
		}
		return fmt.Errorf("looks like interface %v is still not up. Please check The keenetic web-interface", interfaceId)
	})
	return err
}

func (*keeneticInterface) UpInterface(interfaceId string) error {
	var parseSlice []models.ParseRequest
	parseSlice = append(parseSlice, models.ParseRequest{
		Parse: fmt.Sprintf("interface %v up", interfaceId),
	}, models.ParseRequest{
		Parse: "system configuration save",
	})
	var parseResponse []models.ParseResponse
	err := gokeenspinner.WrapWithSpinner(fmt.Sprintf("Bringing %v interface up", color.CyanString(interfaceId)), func() error {
		var executeErr error
		parseResponse, executeErr = ExecutePostParse(parseSlice...)
		return executeErr
	})
	gokeenlog.PrintParseResponse(parseResponse)
	return err
}

func (*keeneticInterface) SetGlobalIpInInterface(interfaceId string, global bool) error {
	var parseSlice []models.ParseRequest
	val := "ip global auto"
	if !global {
		val = "no ip global"
	}
	parseSlice = append(parseSlice, models.ParseRequest{
		Parse: fmt.Sprintf("interface %v %v", interfaceId, val),
	}, models.ParseRequest{
		Parse: "system configuration save",
	})
	var parseResponse []models.ParseResponse
	err := gokeenspinner.WrapWithSpinner(fmt.Sprintf("Changing global IP in %v interface to %v", color.CyanString(interfaceId), color.GreenString("%v", global)), func() error {
		var executeErr error
		parseResponse, executeErr = ExecutePostParse(parseSlice...)
		return executeErr
	})
	gokeenlog.PrintParseResponse(parseResponse)
	return err
}
