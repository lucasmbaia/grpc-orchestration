package main

import (
  "log"
  "testing"
  "github.com/lucasmbaia/grpc-orchestration/server"
)

func TestRunApplication(t *testing.T) {
  var (
    err	  error
    conf  server.ConfigCMD
  )

  if err = conf.Run(); err != nil {
    log.Fatalf("Error to run application: ", err)
  }
}
