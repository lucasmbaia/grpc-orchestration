package main

import (
  "testing"
  "log"
  "github.com/lucasmbaia/grpc-orchestration/tasks"
)

func TestRegisterWorkflows(t *testing.T) {
  var (
    err	  error
    conf  tasks.ConfigWorkflow
  )

  conf.Dir = "/root/workspace/go/src/github.com/lucasmbaia/grpc-orchestration/test/tasks/resources"

  if err = conf.RegisterWorkflows(); err != nil {
    log.Fatalf("Error to register Workflows: ", err)
  }
}

func TestGetWorkflow(t *testing.T) {
  var (
    err	      error
  )

  if _, err = tasks.GetWorkflow(""); err != nil {
    log.Fatalf("Erro to list workflow: ", err)
  }
}

func TestRegisterWorkflowManager(t *testing.T) {
  var (
    w *tasks.WorkflowManager
  )

  w = tasks.RegisterWM("")

  log.Println(w)
}

func TestDeleteWorkflowManagerRegistered(t *testing.T) {
  var (
    key = ""
  )

  tasks.DeleteWM(key)
}
