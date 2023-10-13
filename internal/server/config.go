package server

import (
  "os"
  "encoding/json"
)

type host struct {
  Address string `json:"host"`
  Port string `json:"port"`
  SubscribeSubject string `json:"subscribe_subject"`
}

type base struct {
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type config struct {
  Host host `json:"server"`
  DB base `json:"database"`
}

func CreateConfig(cfgPath string) (*config, error) {
data, err := os.ReadFile(cfgPath)
  if err != nil {
    return nil, err
  }
  config := config{}
  err = json.Unmarshal(data, &config)
  if err != nil {
    return nil, err
  }
  return &config, nil
}
