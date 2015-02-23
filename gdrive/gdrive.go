package gdrive

import (
	"github.com/prasmussen/google-api-go-client/drive/v2"
	"github.com/prasmussen/gdrive/auth"
	"github.com/prasmussen/gdrive/config"
	"github.com/prasmussen/gdrive/util"
	"net/http"
	"path/filepath"
)

// File paths and names
var (
	AppPath     = filepath.Join(util.Homedir(), ".gdrive")
	ConfigFname = "config.json"
	TokenFname  = "token.json"
	//ConfigPath = filepath.Join(ConfigDir, "config.json")
	//TokenPath = filepath.Join(ConfigDir, "token.json")
)

type Drive struct {
	*drive.Service
	client *http.Client
}

// Returns the raw http client which has the oauth transport
func (self *Drive) Client() *http.Client {
	return self.client
}

func New(customAppPath string, advancedMode bool, promptUser bool) (*Drive, error) {
	if customAppPath != "" {
		AppPath = customAppPath
	}

	// Build paths to config files
	configPath := filepath.Join(AppPath, ConfigFname)
	tokenPath := filepath.Join(AppPath, TokenFname)

	config := config.Load(configPath, advancedMode)
	client, err := auth.GetOauth2Client(config.ClientId, config.ClientSecret, tokenPath, promptUser)
	if err != nil {
		return nil, err
	}

	drive, err := drive.New(client)
	if err != nil {
		return nil, err
	}

	// Return a new authorized Drive client.
	return &Drive{drive, client}, nil
}
