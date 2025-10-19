package keenetic

import (
	"testing"

	"github.com/kam1k88/gokeenapi/pkg/backends/keenetic/models"
	"github.com/stretchr/testify/suite"
)

type ApiTestSuite struct {
	GokeenrestapiTestSuite
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (s *ApiTestSuite) TestAuth() {
	err := Common.Auth()
	s.NoError(err)
}

func (s *ApiTestSuite) TestVersion() {
	version, err := Common.Version()
	s.NoError(err)
	s.Equal("KN-1010", version.Model)
	s.Equal("KeeneticOS 3.7.5", version.Title)
}

func (s *ApiTestSuite) TestGetInterfacesViaRciShowInterfaces() {
	interfaces, err := Interface.GetInterfacesViaRciShowInterfaces(false)
	s.NoError(err)
	s.Len(interfaces, 2)

	wg, exists := interfaces["Wireguard0"]
	s.True(exists)
	s.Equal(InterfaceTypeWireguard, wg.Type)
	s.Equal("10.0.0.1/24", wg.Address)
}

func (s *ApiTestSuite) TestGetInterfacesWithTypeFilter() {
	interfaces, err := Interface.GetInterfacesViaRciShowInterfaces(false, InterfaceTypeWireguard)
	s.NoError(err)
	s.Len(interfaces, 1)

	wg, exists := interfaces["Wireguard0"]
	s.True(exists)
	s.Equal(InterfaceTypeWireguard, wg.Type)
}

func (s *ApiTestSuite) TestGetInterfaceViaRciShowInterfaces() {
	iface, err := Interface.GetInterfaceViaRciShowInterfaces("Wireguard0")
	s.NoError(err)
	s.Equal("Wireguard0", iface.Id)
	s.Equal(InterfaceTypeWireguard, iface.Type)
}

func (s *ApiTestSuite) TestGetInterfacesViaRciShowScInterfaces() {
	interfaces, err := Interface.GetInterfacesViaRciShowScInterfaces()
	s.NoError(err)
	s.Len(interfaces, 1)

	wg, exists := interfaces["Wireguard0"]
	s.True(exists)
	s.Equal("Test WireGuard interface", wg.Description)
}

func (s *ApiTestSuite) TestExecutePostParse() {
	parseRequests := []models.ParseRequest{
		{Parse: "interface Wireguard0 up"},
		{Parse: "system configuration save"},
	}

	responses, err := Common.ExecutePostParse(parseRequests...)
	s.NoError(err)
	s.Len(responses, 2)

	for i, response := range responses {
		s.NotEmpty(response.Parse.Status, "Response %d has no status", i)
		s.Equal(StatusOK, response.Parse.Status[0].Status, "Response %d status", i)
	}
}

func (s *ApiTestSuite) TestShowRunningConfig() {
	config, err := Common.ShowRunningConfig()
	s.NoError(err)
	s.Len(config.Message, 2)
	s.Equal("test running config line 1", config.Message[0])
}

func (s *ApiTestSuite) TestUpInterface() {
	err := Interface.UpInterface("Wireguard0")
	s.NoError(err)
}

func (s *ApiTestSuite) TestSetGlobalIpInInterface() {
	err := Interface.SetGlobalIpInInterface("Wireguard0", true)
	s.NoError(err)

	err = Interface.SetGlobalIpInInterface("Wireguard0", false)
	s.NoError(err)
}
