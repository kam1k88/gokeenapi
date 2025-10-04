package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AddAwgTestSuite struct {
	CmdTestSuite
}

func TestAddAwgTestSuite(t *testing.T) {
	suite.Run(t, new(AddAwgTestSuite))
}

func (s *AddAwgTestSuite) TestNewAddAwgCmd() {
	cmd := newAddAwgCmd()

	assert.Equal(s.T(), CmdAddAwg, cmd.Use)
	assert.Equal(s.T(), AliasesAddAwg, cmd.Aliases)
	assert.NotEmpty(s.T(), cmd.Short)
	assert.NotNil(s.T(), cmd.RunE)

	confFileFlag := cmd.Flags().Lookup("conf-file")
	assert.NotNil(s.T(), confFileFlag)

	nameFlag := cmd.Flags().Lookup("name")
	assert.NotNil(s.T(), nameFlag)
}

func (s *AddAwgTestSuite) TestAddAwgCmd_MissingConfFile() {
	cmd := newAddAwgCmd()

	err := cmd.RunE(cmd, []string{})
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "conf-file flag is required")
}

func (s *AddAwgTestSuite) createTestWireGuardConfig() string {
	confContent := `[Interface]
PrivateKey = cOFA+3p5IjkzIjkzIjkzIjkzIjkzIjkzIjkzIjkzIjk=
Address = 10.0.0.2/24
DNS = 8.8.8.8

[Peer]
PublicKey = gN65BkIKy1eCE9pP1wdc8ROUunkiVXrBvGAKBEKdOQI=
Endpoint = example.com:51820
AllowedIPs = 0.0.0.0/0`

	tmpDir := s.T().TempDir()
	confPath := filepath.Join(tmpDir, "test.conf")

	err := os.WriteFile(confPath, []byte(confContent), 0644)
	s.Require().NoError(err)

	return confPath
}

func (s *AddAwgTestSuite) TestAddAwgCmd_InvalidPath() {
	cmd := newAddAwgCmd()
	cmd.Flags().Set("conf-file", "/nonexistent/path.conf")

	err := cmd.RunE(cmd, []string{})
	assert.Error(s.T(), err)
}
