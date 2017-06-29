package main

import (
  "log"
  "testing"
  "github.com/lucasmbaia/grpc-orchestration/server"
  "github.com/lucasmbaia/grpc-orchestration/proto"
  "golang.org/x/net/context"
)

func TestGrpcServerCallTask(t *testing.T) {
  var (
    err	  error
    s	  = server.NewOrchestrationServer()
    task  = new(orchestration.Task)
    ctx	  context.Context
  )

  if _, err = s.CallTask(ctx, task); err != nil {
    log.Fatalf("Erro to call task: ", err)
  }
}
