package server

import (
  "github.com/lucasmbaia/grpc-orchestration/proto"
  "github.com/lucasmbaia/grpc-orchestration/tasks"
  "github.com/lucasmbaia/grpc-orchestration/core"

  "log"
  "golang.org/x/net/context"
  empty "github.com/golang/protobuf/ptypes/empty"
  opentracing "github.com/opentracing/opentracing-go"
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
    result	map[string]string
    span	opentracing.Span
  )

  if workflow, err = tasks.GetWorkflow(t.Name); err != nil {
    return resultTask, err
  }

  wm = tasks.RegisterWM(t.Tracking)
  workflow.InputParameters = t.Parameters

  span = opentracing.SpanFromContext(ctx)
  span.SetTag("tracking", t.Tracking)
  span.SetTag("parameters", t.Parameters)
  span.SetTag("method", t.Name)

  if resultCore, err = core.RunWorkflow(ctx, workflow, wm); err != nil {
    tasks.DeleteWM(t.Tracking)
    return resultTask, err
  }

  if result, err = core.ConvertResults(resultCore); err != nil {
    tasks.DeleteWM(t.Tracking)
    return resultTask, err
  }

  //span.LogKV("result", core.ConvertMapToString(result))

  tasks.DeleteWM(t.Tracking)
  return &orchestration.Result{Response: result}, nil
}

func (o OrchestrationServer) Health(ctx context.Context, emp *empty.Empty) (*empty.Empty, error) {
  log.Println("CHECK")
  select {
  case <-ctx.Done():
    return nil, ctx.Err()
  default:
    return new(empty.Empty), nil
  }
}
