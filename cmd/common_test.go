package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommonTestSuite struct {
	suite.Suite
}

func TestCommonTestSuite(t *testing.T) {
	suite.Run(t, new(CommonTestSuite))
}

func (s *CommonTestSuite) TestCheckRequiredFields_Valid() {
	config.Cfg.Keenetic.URL = "http://192.168.1.1"
	config.Cfg.Keenetic.Login = "admin"
	config.Cfg.Keenetic.Password = "password"

	err := checkRequiredFields()
	assert.NoError(s.T(), err)
}

func (s *CommonTestSuite) TestCheckRequiredFields_MissingURL() {
	config.Cfg.Keenetic.URL = ""
	config.Cfg.Keenetic.Login = "admin"
	config.Cfg.Keenetic.Password = "password"

	err := checkRequiredFields()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "keenetic.url")
}

func (s *CommonTestSuite) TestCheckRequiredFields_MissingLogin() {
	config.Cfg.Keenetic.URL = "http://192.168.1.1"
	config.Cfg.Keenetic.Login = ""
	config.Cfg.Keenetic.Password = "password"

	err := checkRequiredFields()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "keenetic.login")
}

func (s *CommonTestSuite) TestCheckRequiredFields_MissingPassword() {
	config.Cfg.Keenetic.URL = "http://192.168.1.1"
	config.Cfg.Keenetic.Login = "admin"
	config.Cfg.Keenetic.Password = ""

	err := checkRequiredFields()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "keenetic.password")
}

func (s *CommonTestSuite) TestRestoreCursor() {
	// Test that RestoreCursor doesn't panic
	assert.NotPanics(s.T(), func() {
		RestoreCursor()
	})
}

func (s *CommonTestSuite) TestConfirmAction_Yes() {
	// Simulate user input "y"
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, _ := os.Pipe()
	os.Stdin = r
	go func() {
		defer func() { _ = w.Close() }()
		_, _ = w.Write([]byte("y\n"))
	}()

	result, err := confirmAction("Test question?")
	assert.NoError(s.T(), err)
	assert.True(s.T(), result)
}

func (s *CommonTestSuite) TestConfirmAction_No() {
	// Simulate user input "n"
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, _ := os.Pipe()
	os.Stdin = r
	go func() {
		defer func() { _ = w.Close() }()
		_, _ = w.Write([]byte("n\n"))
	}()

	result, err := confirmAction("Test question?")
	assert.NoError(s.T(), err)
	assert.False(s.T(), result)
}

func (s *CommonTestSuite) TestConfirmAction_EOF() {
	// Simulate EOF (Ctrl+D)
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, _ := os.Pipe()
	os.Stdin = r
	_ = w.Close() // Close immediately to simulate EOF

	result, err := confirmAction("Test question?")
	assert.Error(s.T(), err)
	assert.False(s.T(), result)
	assert.True(s.T(), strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "canceled"))
}
