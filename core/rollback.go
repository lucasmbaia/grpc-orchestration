package core

import (
  "reflect"

  "golang.org/x/net/context"
  "github.com/lucasmbaia/grpc-orchestration/tasks"
)

func rollback(w	tasks.Workflow, results Results) {
  var (
    keys      []reflect.Value
    flow      string
    task      string
    count     = -1
    workflow  tasks.Workflow
    err	      error
    index     int
    exists    bool
    args      []reflect.Value
    wm	      = new(tasks.WorkflowManager)
    //r	      Results
  )

  keys = reflect.ValueOf(results).MapKeys()

  for i, v := range w.Tasks {
    for _, key := range keys {
      if v.Name == key.String() && i > count {
	flow = v.Rollback.Name
	task = v.Rollback.Step
	count = i
      }
    }
  }

  for _, values := range results {
    for _, arg := range values {
      args = append(args, arg)
    }
  }

  if workflow, err = tasks.GetWorkflow(flow); err != nil {
    return
  }

  if index, exists = tasks.GetTasksWorkflow(workflow.Tasks, task); !exists {
    return
  }

  workflow.Tasks = workflow.Tasks[index:]
  workflow.InputParameters = w.InputParameters
  wm = tasks.RegisterWM("rollback")

  if _, err = runWorkflow(context.Background(), workflow, wm, true, args); err != nil {
    return
  }

  tasks.DeleteWM("rollback")
  return
}
