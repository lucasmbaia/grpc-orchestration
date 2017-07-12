package main

import (
  "log"

  "github.com/lucasmbaia/grpc-orchestration/config"
  "github.com/lucasmbaia/grpc-orchestration/tasks"
  "github.com/lucasmbaia/grpc-orchestration/server"
)

func init() {
  config.LoadConfig()
}

func main() {
  var (
    err		    error
    configWorkflow  tasks.ConfigWorkflow
    configCMD	    server.ConfigCMD
    errChan	    = make(chan error, 1)
  )

  configWorkflow = tasks.ConfigWorkflow {
    EtcdURL:  config.EnvLocal.EtcdURL,
    Keys:     config.EnvConfig.WorkflowsName,
  }

  if err = configWorkflow.RegisterWorkflows(); err != nil {
    log.Fatalf("Error to register workflows: ", err)
  }

  go func() {
    configCMD = server.ConfigCMD {
      SSL:  true,
    }

    errChan <- configCMD.Run()
  }()

  select {
  case e := <-errChan:
    log.Fatalf("Error grpc server: ", e)
  }
}
