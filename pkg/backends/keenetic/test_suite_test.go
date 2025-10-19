package keenetic

import (
	"net/http/httptest"

	"github.com/stretchr/testify/suite"
)

// GokeenrestapiTestSuite provides common test setup for keenetic tests
type GokeenrestapiTestSuite struct {
	suite.Suite
	server *httptest.Server
}

// SetupSuite runs once before all tests in the suite
func (s *GokeenrestapiTestSuite) SetupSuite() {
	s.server = SetupMockServer()
	SetupTestConfig(s.server.URL)
}

// TearDownSuite runs once after all tests in the suite
func (s *GokeenrestapiTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}
