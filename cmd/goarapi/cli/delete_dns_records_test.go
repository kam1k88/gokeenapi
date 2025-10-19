package cli

import (
	"testing"

	"github.com/kam1k88/gokeenapi/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DeleteDnsRecordsTestSuite struct {
	CmdTestSuite
}

func TestDeleteDnsRecordsTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteDnsRecordsTestSuite))
}

func (s *DeleteDnsRecordsTestSuite) TestNewDeleteDnsRecordsCmd() {
	cmd := newDeleteDnsRecordsCmd()

	assert.Equal(s.T(), CmdDeleteDnsRecords, cmd.Use)
	assert.Equal(s.T(), AliasesDeleteDnsRecords, cmd.Aliases)
	assert.NotEmpty(s.T(), cmd.Short)
	assert.NotNil(s.T(), cmd.RunE)

	forceFlag := cmd.Flags().Lookup("force")
	assert.NotNil(s.T(), forceFlag)
	assert.Equal(s.T(), "bool", forceFlag.Value.Type())
}

func (s *DeleteDnsRecordsTestSuite) TestDeleteDnsRecordsCmd_WithForce() {
	config.Cfg.DNS = config.DNS{
		Records: []config.DnsRecord{
			{
				Domain: "test.local",
				IP:     []string{"192.168.1.100"},
			},
		},
	}

	cmd := newDeleteDnsRecordsCmd()
	_ = cmd.Flags().Set("force", "true")

	err := cmd.RunE(cmd, []string{})
	assert.NoError(s.T(), err)
}

func (s *DeleteDnsRecordsTestSuite) TestDeleteDnsRecordsCmd_EmptyRecords() {
	config.Cfg.DNS = config.DNS{Records: []config.DnsRecord{}}

	cmd := newDeleteDnsRecordsCmd()
	_ = cmd.Flags().Set("force", "true")

	err := cmd.RunE(cmd, []string{})
	assert.NoError(s.T(), err)
}
