package gokeenrestapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"github.com/stretchr/testify/suite"
)

type AwgTestSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *AwgTestSuite) SetupSuite() {
	s.server = s.setupMockServerForAWG()
	SetupTestConfig(s.server.URL)
}

func (s *AwgTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *AwgTestSuite) setupMockServerForAWG() *httptest.Server {
	mux := http.NewServeMux()

	// Add interface endpoints that CheckInterfaceExists needs
	mux.HandleFunc("/rci/show/interface", func(w http.ResponseWriter, r *http.Request) {
		interfaces := map[string]gokeenrestapimodels.RciShowInterface{
			"Wireguard0": {
				Id:          "Wireguard0",
				Type:        "Wireguard",
				Description: "Test WireGuard interface",
				Address:     "10.0.0.1/24",
				Connected:   "yes",
				Link:        "up",
				State:       "up",
			},
		}
		encodeJSON(w, interfaces)
	})

	// Add single SC interface endpoint
	mux.HandleFunc("/rci/show/sc/interface/", func(w http.ResponseWriter, r *http.Request) {
		interfaceId := r.URL.Path[len("/rci/show/sc/interface/"):]
		if interfaceId == "Wireguard0" {
			iface := gokeenrestapimodels.RciShowScInterface{
				Description: "Test WireGuard interface",
				Wireguard: gokeenrestapimodels.Wireguard{
					Asc: gokeenrestapimodels.Asc{
						Jc:   "40", // Different from config to trigger update
						Jmin: "5",
						Jmax: "95",
						S1:   "10",
						S2:   "20",
						H1:   "1",
						H2:   "2",
						H3:   "3",
						H4:   "4",
					},
				},
			}
			encodeJSON(w, iface)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	// Parse endpoint for AWG operations
	mux.HandleFunc("/rci/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var requests []gokeenrestapimodels.ParseRequest
		if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responses := make([]gokeenrestapimodels.ParseResponse, len(requests))
		for i := range requests {
			responses[i] = gokeenrestapimodels.ParseResponse{
				Parse: gokeenrestapimodels.Parse{
					Status: []gokeenrestapimodels.Status{
						{
							Status:  "ok",
							Code:    "0",
							Message: "AWG configuration applied successfully",
						},
					},
				},
			}
		}
		encodeJSON(w, responses)
	})

	return httptest.NewServer(mux)
}

func (s *AwgTestSuite) createTestWireGuardConfig() string {
	confContent := `[Interface]
PrivateKey = cOFA+3p5IjkzIjkzIjkzIjkzIjkzIjkzIjkzIjkzIjk=
Address = 10.0.0.2/24
DNS = 8.8.8.8
Jc = 50
Jmin = 5
Jmax = 95
S1 = 10
S2 = 20
H1 = 1
H2 = 2
H3 = 3
H4 = 4

[Peer]
PublicKey = gN65BkIKy1eCE9pP1wdc8ROUunkiVXrBvGAKBEKdOQI=
Endpoint = example.com:51820
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25`

	tmpDir := s.T().TempDir()
	confPath := filepath.Join(tmpDir, "test.conf")

	err := os.WriteFile(confPath, []byte(confContent), 0644)
	s.Require().NoError(err, "Failed to create test config file")

	return confPath
}

func TestAwgTestSuite(t *testing.T) {
	suite.Run(t, new(AwgTestSuite))
}

func (s *AwgTestSuite) TestConfigureOrUpdateInterface() {
	confPath := s.createTestWireGuardConfig()

	// Test with existing interface - this should work now
	err := AwgConf.ConfigureOrUpdateInterface(confPath, "Wireguard0")
	s.NoError(err)
}

func (s *AwgTestSuite) TestConfigureOrUpdateInterfaceNonExistent() {
	confPath := s.createTestWireGuardConfig()

	// Test with non-existent interface
	err := AwgConf.ConfigureOrUpdateInterface(confPath, "NonExistentInterface")
	s.Error(err)

	// The error should be about interface not found
	s.Contains(err.Error(), "doesn't have interface")
}

func (s *AwgTestSuite) TestConfigureOrUpdateInterfaceEmptyPath() {
	err := AwgConf.ConfigureOrUpdateInterface("", "Wireguard0")
	s.Error(err)
	s.Equal("conf-file flag is required", err.Error())
}
