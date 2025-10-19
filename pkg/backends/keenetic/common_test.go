package keenetic

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kam1k88/gokeenapi/pkg/backends/keenetic/models"
	"github.com/stretchr/testify/suite"
)

type CheckRouterModeTestSuite struct {
	suite.Suite
}

func TestCheckRouterModeTestSuite(t *testing.T) {
	suite.Run(t, new(CheckRouterModeTestSuite))
}

func (s *CheckRouterModeTestSuite) TestCheckRouterMode_RouterMode() {
	// Create separate server for extender mode test
	mux := http.NewServeMux()
	mux.HandleFunc("/rci/show/system/mode", func(w http.ResponseWriter, r *http.Request) {
		mode := models.SystemMode{
			Active:   "router",
			Selected: "router",
		}
		encodeJSON(w, mode)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	SetupTestConfig(server.URL)
	active, selected, err := Common.CheckRouterMode()
	s.NoError(err)
	s.Equal("router", active)
	s.Equal("router", selected)
}

func (s *CheckRouterModeTestSuite) TestCheckRouterMode_ExtenderMode() {
	// Create separate server for extender mode test
	mux := http.NewServeMux()
	mux.HandleFunc("/rci/show/system/mode", func(w http.ResponseWriter, r *http.Request) {
		mode := models.SystemMode{
			Active:   "extender",
			Selected: "extender",
		}
		encodeJSON(w, mode)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	SetupTestConfig(server.URL)

	active, selected, err := Common.CheckRouterMode()
	s.Error(err)
	s.Equal("extender", active)
	s.Equal("extender", selected)
	s.Contains(err.Error(), "router is not in router mode")
}

func (s *CheckRouterModeTestSuite) TestCheckRouterMode_MixedMode() {
	// Create separate server for mixed mode test
	mux := http.NewServeMux()
	mux.HandleFunc("/rci/show/system/mode", func(w http.ResponseWriter, r *http.Request) {
		mode := models.SystemMode{
			Active:   "router",
			Selected: "extender",
		}
		encodeJSON(w, mode)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	SetupTestConfig(server.URL)

	active, selected, err := Common.CheckRouterMode()
	s.Error(err)
	s.Equal("router", active)
	s.Equal("extender", selected)
	s.Contains(err.Error(), "router is not in router mode")
}
