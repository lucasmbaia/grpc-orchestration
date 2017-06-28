package main

import (
  "log"
  "testing"
  "github.com/lucasmbaia/grpc-orchestration/tasks"
)

func TestRegisterTasksReference(t *testing.T) {
  var (
    err	error
    tr	tasks.TasksReferencesInfos
  )

  if err = tasks.RegisterTasksReferences(tr); err != nil {
    log.Fatalf("Error to register tasks: ", err)
  }
}
