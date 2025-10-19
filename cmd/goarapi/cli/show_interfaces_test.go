package cli

import (
	"testing"

	"github.com/kam1k88/gokeenapi/pkg/backends/keenetic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ShowInterfacesTestSuite struct {
	CmdTestSuite
}

func TestShowInterfacesTestSuite(t *testing.T) {
	suite.Run(t, new(ShowInterfacesTestSuite))
}

func (s *ShowInterfacesTestSuite) TestNewShowInterfacesCmd() {
	cmd := newShowInterfacesCmd()

	assert.Equal(s.T(), CmdShowInterfaces, cmd.Use)
	assert.Equal(s.T(), AliasesShowInterfaces, cmd.Aliases)
	assert.NotEmpty(s.T(), cmd.Short)
	assert.NotNil(s.T(), cmd.RunE)

	typeFlag := cmd.Flags().Lookup("type")
	assert.NotNil(s.T(), typeFlag)
	assert.Equal(s.T(), "stringSlice", typeFlag.Value.Type())
}

func (s *ShowInterfacesTestSuite) TestShowInterfacesCmd_Execute() {
	cmd := newShowInterfacesCmd()
	output, err := s.CaptureOutput(cmd, []string{})

	assert.NoError(s.T(), err)
	assert.Contains(s.T(), output, "Wireguard0")
	assert.Contains(s.T(), output, "ISP")
}

func (s *ShowInterfacesTestSuite) TestShowInterfacesCmd_WithTypeFilter() {
	cmd := newShowInterfacesCmd()
	_ = cmd.Flags().Set("type", keenetic.InterfaceTypeWireguard)
	output, err := s.CaptureOutput(cmd, []string{})

	assert.NoError(s.T(), err)
	assert.Contains(s.T(), output, "Wireguard0")
}
