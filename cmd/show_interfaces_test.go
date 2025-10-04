package cmd

import (
	"io"
	"os"
	"testing"

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
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := newShowInterfacesCmd()
	err := cmd.RunE(cmd, []string{})

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.NoError(s.T(), err)
	assert.Contains(s.T(), output, "Wireguard0")
	assert.Contains(s.T(), output, "ISP")
}

func (s *ShowInterfacesTestSuite) TestShowInterfacesCmd_WithTypeFilter() {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := newShowInterfacesCmd()
	cmd.Flags().Set("type", "Wireguard")
	err := cmd.RunE(cmd, []string{})

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.NoError(s.T(), err)
	assert.Contains(s.T(), output, "Wireguard0")
}
