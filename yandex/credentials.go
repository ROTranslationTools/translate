package yandex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	envCredentialsKey = "YANDEX_API_CREDENTIALS"
)

var (
	ErrNilCredentials          = errors.New("nil credentials passed in")
	ErrNoCredentialsDiscovered = fmt.Errorf("no credentials were discovered, please set %q in your environment", envCredentialsKey)
)

type Credentials struct {
	ApiKey string `json:"api_key"`
}

func discoverCredentials() (*Credentials, error) {
	envResolvedFile := os.Getenv(envCredentialsKey)
	if envResolvedFile == "" {
		return nil, ErrNoCredentialsDiscovered
	}

	f, err := os.Open(envResolvedFile)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	_ = f.Close()
	if err != nil {
		return nil, err
	}

	creds := &Credentials{}
	if err := json.Unmarshal(data, creds); err != nil {
		return nil, err
	}
	return creds, nil
}
