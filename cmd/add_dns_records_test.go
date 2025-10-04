package cmd

import (
	"testing"

	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AddDnsRecordsTestSuite struct {
	CmdTestSuite
}

func TestAddDnsRecordsTestSuite(t *testing.T) {
	suite.Run(t, new(AddDnsRecordsTestSuite))
}

func (s *AddDnsRecordsTestSuite) TestNewAddDnsRecordsCmd() {
	cmd := newAddDnsRecordsCmd()

	assert.Equal(s.T(), CmdAddDnsRecords, cmd.Use)
	assert.Equal(s.T(), AliasesAddDnsRecords, cmd.Aliases)
	assert.NotEmpty(s.T(), cmd.Short)
	assert.NotNil(s.T(), cmd.RunE)
}

func (s *AddDnsRecordsTestSuite) TestAddDnsRecordsCmd_Execute() {
	// Set up test config with DNS records
	config.Cfg.DNS = config.DNS{
		Records: []config.DnsRecord{
			{
				Domain: "test.local",
				IP:     []string{"192.168.1.100"},
			},
		},
	}

	cmd := newAddDnsRecordsCmd()
	err := cmd.RunE(cmd, []string{})

	assert.NoError(s.T(), err)
}

func (s *AddDnsRecordsTestSuite) TestAddDnsRecordsCmd_EmptyRecords() {
	// Empty DNS records config
	config.Cfg.DNS = config.DNS{Records: []config.DnsRecord{}}

	cmd := newAddDnsRecordsCmd()
	err := cmd.RunE(cmd, []string{})

	assert.NoError(s.T(), err)
}
