package main

import (
  "log"
  "testing"
  "github.com/lucasmbaia/grpc-orchestration/core"
  "github.com/lucasmbaia/grpc-orchestration/tasks"
)

type TestRW struct {
  Value int `json:",omitempty"`
}

func Teste(t *TestRW) *TestRW {
  return &TestRW{Value: t.Value * t.Value}
}

func TestRunWorkflow(t *testing.T) {
  var (
    err		error
    w		tasks.Workflow
    wm		*tasks.WorkflowManager
    conf	= tasks.ConfigWorkflow{Dir: "/root/workspace/go/src/github.com/lucasmbaia/grpc-orchestration/test/tasks/resources"}
    parameters  = `{"Value":10}`
  )

  if err = conf.RegisterWorkflows(); err != nil {
    log.Fatalf("Error to register Workflows: ", err)
  }

  if err = tasks.RegisterTasksReferences(
    tasks.TasksReferencesInfos {
      "task_Teste": {
	Task: Teste,
      },
    },
  ); err != nil {
    log.Fatalf("Error to register tasks: ", err)
  }

  if w, err = tasks.GetWorkflow("teste"); err != nil {
    log.Fatalf("Erro to list workflow: ", err)
  }

  w.InputParameters = parameters
  if _, err = core.RunWorkflow(w, wm); err != nil {
    log.Fatalf("Erro to exec: ", err)
  }
}
