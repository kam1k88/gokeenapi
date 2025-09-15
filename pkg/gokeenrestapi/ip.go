package gokeenrestapi

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/internal/gokeenspinner"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"go.uber.org/multierr"
)

const (
	regex = `(?i)route ADD (\d+.\d+.\d+.\d+) MASK (\d+.\d+.\d+.\d+)`
)

type keeneticIp struct {
}

var Ip keeneticIp

func (*keeneticIp) GetAllUserRoutesRciIpRoute(keeneticInterface string) ([]gokeenrestapimodels.RciIpRoute, error) {
	var routes []gokeenrestapimodels.RciIpRoute
	err := gokeenspinner.WrapWithSpinner("Fetching static routes", func() error {
		body, err := Common.ExecuteGetSubPath("/rci/ip/route")
		if err != nil {
			return err
		}
		return json.Unmarshal(body, &routes)
	})
	if err != nil {
		return nil, err
	}
	var realRoutes []gokeenrestapimodels.RciIpRoute
	for _, route := range routes {
		route := route
		if route.Interface == keeneticInterface {
			realRoutes = append(realRoutes, route)
		}
	}
	gokeenlog.InfoSubStepf("Found %v static routes for %v interface", color.BlueString("%v", len(realRoutes)), keeneticInterface)
	return realRoutes, err
}

func (*keeneticIp) DeleteRoutes(routes []gokeenrestapimodels.RciIpRoute, interfaceId string) error {
	if len(routes) == 0 {
		gokeenlog.Info("No need to delete static routes")
		return nil
	}
	var parseSlice []gokeenrestapimodels.ParseRequest
	for _, route := range routes {
		if route.Interface != interfaceId {
			continue
		}
		parse := gokeenrestapimodels.ParseRequest{}
		var ip string
		if route.Host != "" {
			ip = route.Host
		}
		if route.Network != "" {
			ip = fmt.Sprintf("%s %s", route.Network, route.Mask)
		}
		parse.Parse = fmt.Sprintf("no ip route %v %v", ip, interfaceId)
		parseSlice = append(parseSlice, parse)
	}
	return gokeenspinner.WrapWithSpinner(fmt.Sprintf("Deleting %v static routes with %v interface", color.BlueString("%v", len(parseSlice)), interfaceId), func() error {
		_, err := Common.ExecutePostParse(parseSlice...)
		return err
	})
}

func (*keeneticIp) AddDnsRecords(domains []string) error {
	var parseSlice []gokeenrestapimodels.ParseRequest
	for _, domain := range domains {
		parse := gokeenrestapimodels.ParseRequest{}
		parse.Parse = fmt.Sprintf("ip host %v", domain)
		parseSlice = append(parseSlice, parse)
	}
	return gokeenspinner.WrapWithSpinner("Adding dns records", func() error {
		_, err := Common.ExecutePostParse(parseSlice...)
		return err
	})
}

func (*keeneticIp) DeleteDnsRecords(domains []string) error {
	var parseSlice []gokeenrestapimodels.ParseRequest
	for _, domain := range domains {
		parse := gokeenrestapimodels.ParseRequest{}
		parse.Parse = fmt.Sprintf("no ip host %v", domain)
		parseSlice = append(parseSlice, parse)
	}
	return gokeenspinner.WrapWithSpinner("Deleting dns records", func() error {
		_, err := Common.ExecutePostParse(parseSlice...)
		return err
	})
}

func (*keeneticIp) AddRoutesFromBatFile(batFile string, interfaceId string) error {
	matcher := regexp.MustCompile(regex)
	b, err := os.ReadFile(batFile)
	if err != nil {
		return err
	}
	str := string(b)
	var mErr error
	splitted := strings.Split(str, "\n")
	var parseSlice []gokeenrestapimodels.ParseRequest
	for _, line := range splitted {
		if line == "" {
			continue
		}
		sl := matcher.FindStringSubmatch(line)
		if len(sl) != 3 {
			gokeenlog.Infof("Skipping line with invalid format: '%v'", line)
			gokeenlog.InfoSubStepf("It doesn't satisfy regexp: '%v'", regex)
			mErr = multierr.Append(mErr, fmt.Errorf("line has invalid format: '%v'", line))
			continue
		}
		ip := sl[1]
		mask := sl[2]
		parseSlice = append(parseSlice, gokeenrestapimodels.ParseRequest{Parse: fmt.Sprintf("ip route %v %v %v auto", ip, mask, interfaceId)})
	}
	var parseResponse []gokeenrestapimodels.ParseResponse
	mErr = multierr.Append(mErr, gokeenspinner.WrapWithSpinner(fmt.Sprintf("Adding %v static routes from %v file to %v interface", color.BlueString("%v", len(parseSlice)), color.CyanString(batFile), color.BlueString(interfaceId)), func() error {
		var executeErr error
		parseResponse, executeErr = Common.ExecutePostParse(parseSlice...)
		return executeErr
	}))
	gokeenlog.PrintParseResponse(parseResponse)
	return mErr
}

func (*keeneticIp) AddRoutesFromBatUrl(url string, interfaceId string) error {
	matcher := regexp.MustCompile(regex)
	rClient := resty.New()
	rClient.SetDisableWarn(true)
	rClient.SetTimeout(time.Second * 5)
	var err error
	var response *resty.Response
	err = gokeenspinner.WrapWithSpinner(fmt.Sprintf("Fetching %v url", color.CyanString(url)), func() error {
		response, err = rClient.R().Get(url)
		return err
	})
	if err != nil {
		return err
	}
	str := string(response.Body())
	var mErr error
	splitted := strings.Split(str, "\n")
	var parseSlice []gokeenrestapimodels.ParseRequest
	for _, line := range splitted {
		if line == "" {
			continue
		}
		sl := matcher.FindStringSubmatch(line)
		if len(sl) != 3 {
			gokeenlog.Infof("Skipping line with invalid format: '%v'", line)
			gokeenlog.InfoSubStepf("It doesn't satisfy regexp: '%v'", regex)
			mErr = multierr.Append(mErr, fmt.Errorf("line has invalid format: '%v'", line))
			continue
		}
		ip := sl[1]
		mask := sl[2]
		parseSlice = append(parseSlice, gokeenrestapimodels.ParseRequest{Parse: fmt.Sprintf("ip route %v %v %v auto", ip, mask, interfaceId)})
	}
	var parseResponse []gokeenrestapimodels.ParseResponse
	mErr = multierr.Append(mErr, gokeenspinner.WrapWithSpinner(fmt.Sprintf("Adding %v static routes to %v interface", color.BlueString("%v", len(parseSlice)), color.BlueString(interfaceId)), func() error {
		var executeErr error
		parseResponse, executeErr = Common.ExecutePostParse(parseSlice...)
		return executeErr
	}))
	gokeenlog.PrintParseResponse(parseResponse)
	return mErr
}
