package core

import (
  "reflect"
  "sync"
  "github.com/lucasmbaia/grpc-orchestration/tasks"
)

type stepTasks struct {
  HasDependency	bool
  DepName	[]string
  Task		string
  Key		string
  FN		reflect.Value
}

type stepInfos struct {
  ReadyTasks	chan stepTasks
  Dependents	map[string][]*sync.WaitGroup
  ReadyToCheck	[]string
  UncheckedDeps	map[string]int
  Mutex		*sync.Mutex
}

func setStepTasks(workflow tasks.Workflow) stepInfos {
  var (
    step  stepInfos
    f	  reflect.Value
    err	  error
  )

  step.ReadyTasks = make(chan stepTasks, len(workflow.Tasks))
  step.Mutex = &sync.Mutex{}

  for _, task := range workflow.Tasks {
    f = reflect.Indirect(reflect.ValueOf(tasks.TR[task.TaskReference].Task))

    if len(task.Dependency) > 0 {
      var wg = &sync.WaitGroup{}
      wg.Add(1)

      for _, name := range task.Dependency {
	step.Dependents[name] = append(step.Dependents[name], wg)
      }

      step.UncheckedDeps[task.Name] = len(task.Dependency)

      go func(w *sync.WaitGroup, task reflect.Value, key, taskReference string, dep []string) {
	w.Wait()

	step.Mutex.Lock()

	if err != nil {
	  step.Mutex.Unlock()
	  return
	}

	step.Mutex.Unlock()

	step.ReadyTasks <- stepTasks{FN: task, HasDependency: true, Key: key, DepName: dep, Task: taskReference}
      }(wg, f, task.Name, task.TaskReference, task.Dependency)
    } else {
      step.ReadyToCheck = append(step.ReadyToCheck, task.Name)
      step.ReadyTasks <- stepTasks{FN: f, Task: task.TaskReference, Key: task.Name}
    }
  }

  return step
}

func closeTasks(s stepInfos) {
  for _, dep := range s.Dependents {
    for _, wg := range dep {
      wg.Done()
    }
  }

  close(s.ReadyTasks)

  return
}
