package main

import (
  "log"
  "testing"
  "github.com/lucasmbaia/grpc-orchestration/client"
  "github.com/lucasmbaia/grpc-orchestration/proto"
)

func TestCallTaskClient(t *testing.T) {
  var (
    config	= client.Config{SSL: true}
    err		error
    task	*orchestration.Task
    parameters	= `{"Value":10}`
    result	*orchestration.Result
    done	= make(chan bool, 1)
  )

  task = &orchestration.Task{Tracking: "lucas-teste", Name: "fibonacci_task", Parameters: parameters}

  for i := 0; i < 100; i++ {
    go func() {
      for {
	if result, err = config.CallTask(task); err != nil {
	  //log.Fatalf("Erro to call task: ", err)
	  log.Println("Erro to call task: ", err)
	} else {
	  log.Println(result)
	}
      }
    }()
  }

  <-done
}
