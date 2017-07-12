package main

import (
  "testing"
  "log"
  "github.com/lucasmbaia/grpc-orchestration/consul"
)

func TestRegisterConsul(t *testing.T) {
  var (
    err	    error
  )

  if err = consul.RegisterService(); err != nil {
    log.Fatalf("Erro to register service: ", err)
  }
}
