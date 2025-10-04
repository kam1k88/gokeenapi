package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigureAwgTestSuite struct {
	CmdTestSuite
}

func TestConfigureAwgTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigureAwgTestSuite))
}

func (s *ConfigureAwgTestSuite) TestNewUpdateAwgCmd() {
	cmd := newUpdateAwgCmd()

	assert.Equal(s.T(), CmdUpdateAwg, cmd.Use)
	assert.Equal(s.T(), AliasesUpdateAwg, cmd.Aliases)
	assert.NotEmpty(s.T(), cmd.Short)
	assert.NotNil(s.T(), cmd.RunE)

	confFileFlag := cmd.Flags().Lookup("conf-file")
	assert.NotNil(s.T(), confFileFlag)

	interfaceIdFlag := cmd.Flags().Lookup("interface-id")
	assert.NotNil(s.T(), interfaceIdFlag)
}

func (s *ConfigureAwgTestSuite) TestUpdateAwgCmd_MissingConfFile() {
	cmd := newUpdateAwgCmd()
	cmd.Flags().Set("interface-id", "Wireguard0")

	err := cmd.RunE(cmd, []string{})
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "--conf-file flag is required")
}

func (s *ConfigureAwgTestSuite) TestUpdateAwgCmd_MissingInterfaceId() {
	cmd := newUpdateAwgCmd()
	cmd.Flags().Set("conf-file", "/tmp/test.conf")

	err := cmd.RunE(cmd, []string{})
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "--interface-id flag is required")
}
