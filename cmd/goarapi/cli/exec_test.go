package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExecTestSuite struct {
	CmdTestSuite
}

func TestExecTestSuite(t *testing.T) {
	suite.Run(t, new(ExecTestSuite))
}

func (s *ExecTestSuite) TestNewExecCmd() {
	cmd := newExecCmd()

	assert.Equal(s.T(), CmdExec, cmd.Use)
	assert.Equal(s.T(), AliasesExec, cmd.Aliases)
	assert.NotEmpty(s.T(), cmd.Short)
	assert.NotNil(s.T(), cmd.RunE)
}

func (s *ExecTestSuite) TestExecCmd_Execute() {
	cmd := newExecCmd()

	err := cmd.RunE(cmd, []string{"system", "configuration", "save"})
	assert.NoError(s.T(), err)
}

func (s *ExecTestSuite) TestExecCmd_NoArgs() {
	cmd := newExecCmd()

	err := cmd.RunE(cmd, []string{})
	assert.NoError(s.T(), err)
}
