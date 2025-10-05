package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DeleteKnownHostsTestSuite struct {
	CmdTestSuite
}

func TestDeleteKnownHostsTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteKnownHostsTestSuite))
}

func (s *DeleteKnownHostsTestSuite) TestNewDeleteKnownHostsCmd() {
	cmd := newDeleteKnownHostsCmd()

	assert.Equal(s.T(), CmdDeleteKnownHosts, cmd.Use)
	assert.Equal(s.T(), AliasesDeleteKnownHosts, cmd.Aliases)
	assert.NotEmpty(s.T(), cmd.Short)
	assert.NotNil(s.T(), cmd.RunE)

	namePatternFlag := cmd.Flags().Lookup("name-pattern")
	assert.NotNil(s.T(), namePatternFlag)

	macPatternFlag := cmd.Flags().Lookup("mac-pattern")
	assert.NotNil(s.T(), macPatternFlag)

	forceFlag := cmd.Flags().Lookup("force")
	assert.NotNil(s.T(), forceFlag)
	assert.Equal(s.T(), "bool", forceFlag.Value.Type())
}

func (s *DeleteKnownHostsTestSuite) TestDeleteKnownHostsCmd_NoPattern() {
	cmd := newDeleteKnownHostsCmd()

	err := cmd.RunE(cmd, []string{})
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "exactly one of --name-pattern or --mac-pattern must be specified")
}

func (s *DeleteKnownHostsTestSuite) TestDeleteKnownHostsCmd_BothPatterns() {
	cmd := newDeleteKnownHostsCmd()
	_ = cmd.Flags().Set("name-pattern", "test")
	_ = cmd.Flags().Set("mac-pattern", "aa:bb")

	err := cmd.RunE(cmd, []string{})
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "exactly one of --name-pattern or --mac-pattern must be specified")
}

func (s *DeleteKnownHostsTestSuite) TestDeleteKnownHostsCmd_InvalidRegex() {
	cmd := newDeleteKnownHostsCmd()
	_ = cmd.Flags().Set("name-pattern", "[invalid")

	err := cmd.RunE(cmd, []string{})
	assert.Error(s.T(), err)
}

func (s *DeleteKnownHostsTestSuite) TestDeleteKnownHostsCmd_NamePattern() {
	cmd := newDeleteKnownHostsCmd()
	_ = cmd.Flags().Set("name-pattern", "nonexistent")
	_ = cmd.Flags().Set("force", "true")

	err := cmd.RunE(cmd, []string{})
	assert.NoError(s.T(), err)
}

func (s *DeleteKnownHostsTestSuite) TestDeleteKnownHostsCmd_MacPattern() {
	cmd := newDeleteKnownHostsCmd()
	_ = cmd.Flags().Set("mac-pattern", "nonexistent")
	_ = cmd.Flags().Set("force", "true")

	err := cmd.RunE(cmd, []string{})
	assert.NoError(s.T(), err)
}
