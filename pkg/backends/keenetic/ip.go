package keenetic

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/kam1k88/gokeenapi/internal/gokeencache"
	"github.com/kam1k88/gokeenapi/internal/gokeenlog"
	"github.com/kam1k88/gokeenapi/internal/gokeenspinner"
	"github.com/kam1k88/gokeenapi/pkg/backends/keenetic/models"
	"go.uber.org/multierr"
)

const (
	regex = `(?i)route ADD (\d+.\d+.\d+.\d+) MASK (\d+.\d+.\d+.\d+)`
)

type keeneticIp struct {
}

var (
	// Ip provides IP-related functionality including routing, DNS, and host management
	Ip keeneticIp
)

// GetAllHotspots retrieves all known hosts (devices) from the router's hotspot database
func (*keeneticIp) GetAllHotspots() (models.RciShowIpHotspot, error) {
	var hotspot models.RciShowIpHotspot
	err := gokeenspinner.WrapWithSpinner("Fetching hotspots", func() error {
		body, err := Common.ExecuteGetSubPath("/rci/show/ip/hotspot")
		if err != nil {
			return err
		}
		return json.Unmarshal(body, &hotspot)
	})
	if err == nil {
		gokeenlog.InfoSubStepf("Found %v hosts", color.BlueString("%v", len(hotspot.Host)))
	}
	return hotspot, err
}

// ShowIpRoute retrieves all routes from the router's routing table
func (*keeneticIp) ShowIpRoute(interfaceId string) ([]models.RciShowIpRoute, error) {
	routes := gokeencache.GetRciShowIpRoute()
	if routes == nil {
		err := gokeenspinner.WrapWithSpinner(fmt.Sprintf("Fetching %v table", color.CyanString("ip routing")), func() error {
			body, err := Common.ExecuteGetSubPath("/rci/show/ip/route")
			if err != nil {
				return err
			}
			return json.Unmarshal(body, &routes)
		})
		if err != nil {
			return routes, err
		}
		gokeencache.SetRciShowIpRoute(routes)
	}
	var realRoutes []models.RciShowIpRoute
	for _, route := range routes {
		if route.Interface == interfaceId || interfaceId == "" {
			realRoutes = append(realRoutes, route)
		}
	}
	//if err == nil {
	//	gokeenlog.InfoSubStepf("Found %v routes", color.BlueString("%v", len(realRoutes)))
	//}
	return realRoutes, nil
}

// DeleteKnownHosts removes devices from the router's known hosts list by MAC address
func (*keeneticIp) DeleteKnownHosts(hostMacs []string) error {
	if len(hostMacs) == 0 {
		gokeenlog.Info("No need to delete known hosts")
		return nil
	}
	var parseSlice []models.ParseRequest
	for _, mac := range hostMacs {
		parse := models.ParseRequest{
			Parse: fmt.Sprintf("no known host \"%v\"", mac),
		}
		parseSlice = append(parseSlice, parse)
	}
	return gokeenspinner.WrapWithSpinner(fmt.Sprintf("Deleting %v known hosts", color.BlueString("%v", len(parseSlice))), func() error {
		parseSlice = Common.EnsureSaveConfigAtEnd(parseSlice)
		_, err := Common.ExecutePostParse(parseSlice...)
		return err
	})
}

// GetAllUserRoutesRciIpRoute retrieves all user-defined static routes for a specific interface
func (*keeneticIp) GetAllUserRoutesRciIpRoute(keeneticInterface string) ([]models.RciIpRoute, error) {
	var routes []models.RciIpRoute
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
	var realRoutes []models.RciIpRoute
	for _, route := range routes {
		route := route
		if route.Interface == keeneticInterface {
			realRoutes = append(realRoutes, route)
		}
	}
	gokeenlog.InfoSubStepf("Found %v static routes for %v interface", color.BlueString("%v", len(realRoutes)), keeneticInterface)
	return realRoutes, err
}

// DeleteRoutes removes static routes from the specified interface
func (*keeneticIp) DeleteRoutes(routes []models.RciIpRoute, interfaceId string) error {
	if len(routes) == 0 {
		gokeenlog.Info("No need to delete static routes")
		return nil
	}
	var parseSlice []models.ParseRequest
	for _, route := range routes {
		if route.Interface != interfaceId {
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
		parse.Parse = fmt.Sprintf("no ip route %v %v", ip, interfaceId)
		parseSlice = append(parseSlice, parse)
	}
	return gokeenspinner.WrapWithSpinner(fmt.Sprintf("Deleting %v static routes with %v interface", color.BlueString("%v", len(parseSlice)), interfaceId), func() error {
		parseSlice = Common.EnsureSaveConfigAtEnd(parseSlice)
		_, err := Common.ExecutePostParse(parseSlice...)
		return err
	})
}

// AddDnsRecords adds static DNS records to the router configuration
func (*keeneticIp) AddDnsRecords(domains []string) error {
	var parseSlice []models.ParseRequest
	for _, domain := range domains {
		parse := models.ParseRequest{}
		parse.Parse = fmt.Sprintf("ip host %v", domain)
		parseSlice = append(parseSlice, parse)
	}
	return gokeenspinner.WrapWithSpinner("Adding dns records", func() error {
		parseSlice = Common.EnsureSaveConfigAtEnd(parseSlice)
		_, err := Common.ExecutePostParse(parseSlice...)
		return err
	})
}

// DeleteDnsRecords removes static DNS records from the router configuration
func (*keeneticIp) DeleteDnsRecords(domains []string) error {
	var parseSlice []models.ParseRequest
	for _, domain := range domains {
		parse := models.ParseRequest{}
		parse.Parse = fmt.Sprintf("no ip host %v", domain)
		parseSlice = append(parseSlice, parse)
	}
	return gokeenspinner.WrapWithSpinner("Deleting dns records", func() error {
		parseSlice = Common.EnsureSaveConfigAtEnd(parseSlice)
		_, err := Common.ExecutePostParse(parseSlice...)
		return err
	})
}

// AddRoutesFromBatFile parses a local .bat file and adds the contained routes to the specified interface
func (*keeneticIp) AddRoutesFromBatFile(batFile string, interfaceId string) error {
	routes, err := Ip.ShowIpRoute(interfaceId)
	if err != nil {
		return err
	}
	matcher := regexp.MustCompile(regex)
	b, err := os.ReadFile(batFile)
	if err != nil {
		return err
	}
	str := string(b)
	var mErr error
	splitted := strings.Split(str, "\n")
	var parseSlice []models.ParseRequest
	for _, line := range splitted {
		if line == "" {
			continue
		}
		sl := matcher.FindStringSubmatch(line)
		if len(sl) != 3 {
			gokeenlog.InfoSubStepf("Skipping line with invalid format: '%v'", line)
			gokeenlog.InfoSubStepf("It doesn't satisfy regexp: '%v'", regex)
			mErr = multierr.Append(mErr, fmt.Errorf("line has invalid format: '%v'", line))
			continue
		}
		ip := sl[1]
		mask := sl[2]
		contains, err := checkInterfaceContainsRoute(ip, mask, interfaceId, routes)
		if err != nil {
			return err
		}
		if contains {
			//gokeenlog.Infof("Skipping line with already existing route: '%v'", line)
			continue
		}
		parseSlice = append(parseSlice, models.ParseRequest{Parse: fmt.Sprintf("ip route %v %v %v auto", ip, mask, interfaceId)})
	}
	if len(parseSlice) == 0 {
		gokeenlog.InfoSubStepf("No need to add new static routes from %v file", color.CyanString("%v", batFile))
		return nil
	}
	gokeencache.SetRciShowIpRoute(nil)
	var parseResponse []models.ParseResponse
	mErr = multierr.Append(mErr, gokeenspinner.WrapWithSpinner(fmt.Sprintf("Adding new %v static routes from %v file to %v interface", color.CyanString("%v", len(parseSlice)), color.CyanString(batFile), color.BlueString(interfaceId)), func() error {
		var executeErr error
		parseSlice = Common.EnsureSaveConfigAtEnd(parseSlice)
		parseResponse, executeErr = Common.ExecutePostParse(parseSlice...)
		return executeErr
	}))
	gokeenlog.PrintParseResponse(parseResponse)
	return mErr
}

// AddRoutesFromBatUrl downloads a .bat file from a URL and adds the contained routes to the specified interface
func (*keeneticIp) AddRoutesFromBatUrl(url string, interfaceId string) error {
	routes, err := Ip.ShowIpRoute(interfaceId)
	if err != nil {
		return err
	}
	matcher := regexp.MustCompile(regex)
	rClient := resty.New()
	rClient.SetDisableWarn(true)
	rClient.SetTimeout(time.Second * 5)
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
	var parseSlice []models.ParseRequest
	for _, line := range splitted {
		if line == "" {
			continue
		}
		sl := matcher.FindStringSubmatch(line)
		if len(sl) != 3 {
			gokeenlog.InfoSubStepf("Skipping line with invalid format: '%v'", line)
			gokeenlog.InfoSubStepf("It doesn't satisfy regexp: '%v'", regex)
			mErr = multierr.Append(mErr, fmt.Errorf("line has invalid format: '%v'", line))
			continue
		}
		ip := sl[1]
		mask := sl[2]
		contains, err := checkInterfaceContainsRoute(ip, mask, interfaceId, routes)
		if err != nil {
			return err
		}
		if contains {
			//gokeenlog.Infof("Skipping line with already existing route: '%v'", line)
			continue
		}
		parseSlice = append(parseSlice, models.ParseRequest{Parse: fmt.Sprintf("ip route %v %v %v auto", ip, mask, interfaceId)})
	}
	if len(parseSlice) == 0 {
		gokeenlog.InfoSubStepf("No need to add new static routes from %v url", color.CyanString("%v", url))
		return nil
	}
	gokeencache.SetRciShowIpRoute(nil)
	var parseResponse []models.ParseResponse
	mErr = multierr.Append(mErr, gokeenspinner.WrapWithSpinner(fmt.Sprintf("Adding new %v static routes to %v interface", color.CyanString("%v", len(parseSlice)), color.BlueString(interfaceId)), func() error {
		var executeErr error
		parseSlice = Common.EnsureSaveConfigAtEnd(parseSlice)
		parseResponse, executeErr = Common.ExecutePostParse(parseSlice...)
		return executeErr
	}))
	gokeenlog.PrintParseResponse(parseResponse)
	return mErr
}

func maskToCIDR(mask string) (int, error) {
	ip := net.ParseIP(mask)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP")
	}

	ipv4Mask := net.IPMask(ip.To4())
	ones, _ := ipv4Mask.Size()
	return ones, nil
}

func checkInterfaceContainsRoute(routeIp, mask, interfaceId string, existingRoutes []models.RciShowIpRoute) (bool, error) {
	cidr, err := maskToCIDR(mask)
	if err != nil {
		return false, err
	}

	_, newNetwork, err := net.ParseCIDR(fmt.Sprintf("%v/%d", routeIp, cidr))
	if err != nil {
		return false, err
	}

	for _, route := range existingRoutes {
		// skip default 0.0.0.0/0
		if strings.EqualFold(route.Destination, "0.0.0.0/0") {
			continue
		}
		if route.Interface != interfaceId {
			continue
		}

		// Check exact match
		destination := fmt.Sprintf("%v/%d", routeIp, cidr)
		if route.Destination == destination {
			return true, nil
		}

		// Check if existing route covers the new route
		_, existingNetwork, err := net.ParseCIDR(route.Destination)
		if err != nil {
			continue
		}
		if existingNetwork.Contains(newNetwork.IP) {
			return true, nil
		}
	}
	return false, nil
}
