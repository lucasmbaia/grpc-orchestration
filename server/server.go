package server

import (
  "github.com/lucasmbaia/grpc-orchestration/proto"
  "github.com/lucasmbaia/grpc-orchestration/tasks"
  "github.com/lucasmbaia/grpc-orchestration/core"

  "golang.org/x/net/context"
  empty "github.com/golang/protobuf/ptypes/empty"
)

type OrchestrationServer struct {}

func NewOrchestrationServer() OrchestrationServer {
  return OrchestrationServer{}
}

func (o OrchestrationServer) CallTask(ctx context.Context, t *orchestration.Task) (*orchestration.Result, error) {
  var (
    err		error
    workflow	tasks.Workflow
    wm		= new(tasks.WorkflowManager)
    resultCore	core.Results
    resultTask	*orchestration.Result
  )

  if workflow, err = tasks.GetWorkflow(t.Name); err != nil {
    return resultTask, err
  }

  wm = tasks.RegisterWM(t.Tracking)
  workflow.InputParameters = t.Parameters

  if resultCore, err = core.RunWorkflow(workflow, wm); err != nil {
    return resultTask, err
  }

  if resultTask.Response, err = core.ConvertResults(resultCore); err != nil {
    return resultTask, err
  }

  return resultTask, nil
}

func (o OrchestrationServer) Health(ctx context.Context, emp *empty.Empty) (*empty.Empty, error) {
  select {
  case <-ctx.Done():
    return nil, ctx.Err()
  default:
    return new(empty.Empty), nil
  }
}
