package cmd

import (
	"io"
	"net/http/httptest"
	"os"

	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/spf13/cobra"
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

// CaptureOutput executes a command and captures its stdout output
func (s *CmdTestSuite) CaptureOutput(cmd *cobra.Command, args []string) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.RunE(cmd, args)

	_ = w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	return string(out), err
}
