package cmd

import (
	"testing"

	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DeleteRoutesTestSuite struct {
	CmdTestSuite
}

func TestDeleteRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteRoutesTestSuite))
}

func (s *DeleteRoutesTestSuite) TestNewDeleteRoutesCmd() {
	cmd := newDeleteRoutesCmd()

	assert.Equal(s.T(), CmdDeleteRoutes, cmd.Use)
	assert.Equal(s.T(), AliasesDeleteRoutes, cmd.Aliases)
	assert.NotEmpty(s.T(), cmd.Short)
	assert.NotNil(s.T(), cmd.RunE)

	interfaceIdFlag := cmd.Flags().Lookup("interface-id")
	assert.NotNil(s.T(), interfaceIdFlag)

	forceFlag := cmd.Flags().Lookup("force")
	assert.NotNil(s.T(), forceFlag)
	assert.Equal(s.T(), "bool", forceFlag.Value.Type())
}

func (s *DeleteRoutesTestSuite) TestDeleteRoutesCmd_WithInterfaceId() {
	config.Cfg.Routes = []config.Route{
		{InterfaceID: "Wireguard0"},
	}

	cmd := newDeleteRoutesCmd()
	cmd.Flags().Set("interface-id", "Wireguard0")
	cmd.Flags().Set("force", "true")

	err := cmd.RunE(cmd, []string{})
	assert.NoError(s.T(), err)
}

func (s *DeleteRoutesTestSuite) TestDeleteRoutesCmd_WithForceFlag() {
	config.Cfg.Routes = []config.Route{
		{InterfaceID: "Wireguard0"},
	}

	cmd := newDeleteRoutesCmd()
	cmd.Flags().Set("force", "true")

	err := cmd.RunE(cmd, []string{})
	assert.NoError(s.T(), err)
}
