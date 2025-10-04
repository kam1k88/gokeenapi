package cmd

import (
	"net/http/httptest"

	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/stretchr/testify/suite"
)

// CmdTestSuite provides common test setup for all command tests
type CmdTestSuite struct {
	suite.Suite
	server *httptest.Server
}

// SetupSuite runs once before all tests in the suite
func (s *CmdTestSuite) SetupSuite() {
	s.server = gokeenrestapi.SetupMockServer()
	gokeenrestapi.SetupTestConfig(s.server.URL)

	err := gokeenrestapi.Common.Auth()
	s.Require().NoError(err)
}

// TearDownSuite runs once after all tests in the suite
func (s *CmdTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}
