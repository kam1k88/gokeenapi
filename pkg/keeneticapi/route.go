package keeneticapi

import (
	"encoding/json"
	"fmt"
	"github.com/noksa/gokeenapi/internal/config"
	"github.com/noksa/gokeenapi/internal/keenlog"
	"github.com/noksa/gokeenapi/internal/keenspinner"
	"github.com/noksa/gokeenapi/pkg/models"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type keeneticRoute struct {
}

var Route keeneticRoute

func (*keeneticRoute) GetAllUserRoutesRciIpRoute(keeneticInterface string) ([]models.RciIpRoute, error) {
	var routes []models.RciIpRoute
	err := keenspinner.WrapWithSpinner("Fetching static routes", func() error {
		body, err := ExecuteGetSubPath("/ip/route")
		if err != nil {
			return err
		}
		return json.Unmarshal(body, &routes)
	})
	if err != nil {
		return nil, err
	}
	var realRoutes []models.RciIpRoute
	for _, route := range routes {
		route := route
		if route.Interface == keeneticInterface {
			realRoutes = append(realRoutes, route)
		}
	}
	keenlog.Infof("Got %v static routes", len(realRoutes))
	return realRoutes, err
}

func (*keeneticRoute) DeleteRoutes(routes []models.RciIpRoute) error {
	if len(routes) == 0 {
		keenlog.Info("No need to delete static routes")
		return nil
	}
	var parseSlice []models.ParseRequest
	keeneticInterface := viper.GetString(config.ViperKeeneticInterface)
	for _, route := range routes {
		if route.Interface != keeneticInterface {
			continue
		}
		parse := models.ParseRequest{}
		var ip string
		if route.Host != "" {
			ip = route.Host
		}
		if route.Network != "" {
			ip = fmt.Sprintf("%s %s", route.Network, route.Mask)
		}
		parse.Parse = fmt.Sprintf("no ip route %v %v", ip, keeneticInterface)
		parseSlice = append(parseSlice, parse)
	}
	return keenspinner.WrapWithSpinner(fmt.Sprintf("Deleting %v static routes with %v interface", len(parseSlice), keeneticInterface), func() error {
		_, err := ExecutePostParse(parseSlice...)
		return err
	})
}

func (*keeneticRoute) AddRoutesFromBatFile(batFile string) error {
	absBatFile, err := filepath.Abs(batFile)
	if err != nil {
		return err
	}
	matcher := regexp.MustCompile(`route ADD (\d+.\d+.\d+.\d+) MASK (\d+.\d+.\d+.\d+)`)
	b, err := os.ReadFile(absBatFile)
	if err != nil {
		return err
	}
	str := string(b)
	splitted := strings.Split(str, "\n")
	var parseSlice []models.ParseRequest
	for _, line := range splitted {
		if line == "" {
			continue
		}
		sl := matcher.FindStringSubmatch(line)
		if len(sl) != 3 {
			continue
		}
		ip := sl[1]
		mask := sl[2]
		parseSlice = append(parseSlice, models.ParseRequest{Parse: fmt.Sprintf("ip route %v %v %v auto", ip, mask, viper.GetString(config.ViperKeeneticInterface))})
	}
	var parseResponse []models.ParseResponse
	err = keenspinner.WrapWithSpinner(fmt.Sprintf("Adding/Updating %v static routes to %v interface", len(parseSlice), viper.GetString(config.ViperKeeneticInterface)), func() error {
		var executeErr error
		parseResponse, executeErr = ExecutePostParse(parseSlice...)
		return executeErr
	})
	keenlog.PrintParseResponse(parseResponse)
	return err
}
