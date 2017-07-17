package main

import (
  "log"
  "reflect"

  "github.com/lucasmbaia/grpc-orchestration/tasks"
  "github.com/lucasmbaia/grpc-orchestration/server"
  "github.com/lucasmbaia/grpc-orchestration/proto"
  "github.com/lucasmbaia/grpc-base/config"
  "github.com/lucasmbaia/grpc-base/base"
  mathoperations "github.com/lucasmbaia/grpc-mathoperations/client"
  fibonacci "github.com/lucasmbaia/grpc-fibonacci/client"
)

func init() {
  config.LoadConfig()
}

func main() {
  var (
    err		    error
    configWorkflow  tasks.ConfigWorkflow
    configCMD	    base.ConfigCMD
    errChan	    = make(chan error, 1)
    fib		    fibonacci.Config
    math	    mathoperations.Config
  )

  configWorkflow = tasks.ConfigWorkflow {
    EtcdURL:  config.EnvLocal.EtcdURL,
    Keys:     config.EnvConfig.WorkflowsName,
  }

  if err = configWorkflow.RegisterWorkflows(); err != nil {
    log.Fatalf("Error to register workflows: ", err)
  }

  if err = tasks.RegisterTasksReferences(tasks.TasksReferencesInfos{
    "task_Fibonacci": {
      Task:	  fib.CalcFibonacci,
      StructTask: reflect.TypeOf(fib),
    },
    "task_Multiplus": {
      Task:	  math.CalcDouble,
      StructTask: reflect.TypeOf(math),
    },
  }); err != nil {
    log.Fatalf("Error to register tasks: ", err)
  }

  go func() {
    configCMD = base.ConfigCMD {
      SSL:		true,
      RegisterConsul:	true,
      ServiceServer:	reflect.Indirect(reflect.ValueOf(orchestration.RegisterOrchestrationServiceServer)),
      HandlerEndpoint:	reflect.Indirect(reflect.ValueOf(orchestration.RegisterOrchestrationServiceHandlerFromEndpoint)),
      ServerConfig:	server.NewOrchestrationServer(),
    }

    errChan <- configCMD.Run()
  }()

  select {
  case e := <-errChan:
    log.Fatalf("Error grpc server: ", e)
  }
}
