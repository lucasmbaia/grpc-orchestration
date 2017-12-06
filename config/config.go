package config

import (
  "log"
  "os"

  "github.com/lucasmbaia/grpc-orchestration/consul"
)

var (
  EnvLocal  Config
  parsed    = false
)

type Config struct {
  IPs	    []string
  Hostname  string
}

func LoadConfig() {
  var (
    err		error
  )

  if !parsed {
    if EnvLocal.IPs, err = consul.GetIPs(); err != nil {
      log.Fatalf("Error to get ips: ", err)
    }

    if EnvLocal.Hostname, err = os.Hostname(); err != nil {
      log.Fatalf("Error to get hostname: ", err)
    }

    parsed = true
  }
}
