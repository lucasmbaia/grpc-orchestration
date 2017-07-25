package core

import (
  "encoding/json"
  "reflect"
  "strings"
  "github.com/lucasmbaia/grpc-orchestration/tasks"
  "runtime"
  "path/filepath"
)

func RunWorkflow(workflow tasks.Workflow, wm *tasks.WorkflowManager) (Results, error) {
  var (
    results Results
    err	    error
    args    []reflect.Value
  )

  if results, err = runWorkflow(workflow, wm, false, args); err != nil {
    //if workflow.Rollback.Name != "" {
      //go rollback(workflow.Rollback.Name, workflow.Rollback.Step, results)
    //}
  }

    go rollback(workflow, results)
  return results, err
}

func runWorkflow(workflow tasks.Workflow, wm *tasks.WorkflowManager, rollback bool, argsRollback []reflect.Value) (Results, error) {
  var (
    steps	stepInfos
    results	= make(Results, len(workflow.Tasks))
    sA		reflect.Value
    totalTasks	= len(workflow.Tasks)
    body        = []byte(workflow.InputParameters)
    err		error
    errTask	= make(chan error, len(workflow.Tasks))
    rep         = strings.NewReplacer(".", "", "-", "", "fm", "")
    size	= len(workflow.Tasks)
  )

  steps = setStepTasks(workflow, errTask)

  for task := range steps.ReadyTasks {
    go func(st stepTasks) {
      var (
	args	[]reflect.Value
	argsRes	[]reflect.Value
	output	[]reflect.Value
	fName	string
	rName	string
	ref	reflect.Value
	ok	bool
	t	reflect.Type
	argsDep	[]reflect.Value
      )

      defer steps.Mutex.Unlock()
      t = st.FN.Type()

      if t.NumIn() > 0 {
	for i := 0; i < t.NumIn(); i++ {
	  if rollback {
	    for _, dep := range argsRollback {
	      argsDep = append(argsDep, dep)
	    }
	  }

	  for _, depName := range st.DepName {
	    for _, dep := range results[depName] {
	      argsDep = append(argsDep, dep)
	    }
	  }

	  sA, _ = setArgs(argsDep, body, t.In(i))
	  args = append(args, sA)
	}
      }

      if tasks.TR[st.Task].StructTask != nil {
	rName = runtime.FuncForPC(reflect.ValueOf(tasks.TR[st.Task].Task).Pointer()).Name()
	fName = rep.Replace(filepath.Ext(rName))

	ref = reflect.New(tasks.TR[st.Task].StructTask)

	if len(body) > 0 {
	  if err = json.Unmarshal(body, ref.Interface()); err != nil {
	    for i := 0; i < size; i++ {
	      errTask <-err
	    }

	    closeTasks(steps)
	    return
	  }
	}

	output = reflect.ValueOf(ref.Interface()).MethodByName(fName).Call(args)
      } else {
	output = st.FN.Call(args)
      }

      steps.Mutex.Lock()

      if err == nil {
	for r := range output {
	  if t.Out(r).Kind() != reflect.Interface {
	    argsRes = append(argsRes, output[r])
	  } else {
	    if err, ok = output[r].Interface().(error); ok {
	      for i := 0; i < size; i++ {
		errTask <-err
	      }

	      closeTasks(steps)
	      return
	    }
	  }
	}

	results[st.Key] = argsRes

	if wgs, ok := steps.Dependents[st.Key]; ok {
	  for _, wg := range wgs {
	    wg.Done()
	  }

	  delete(steps.Dependents, st.Key)
	}

	totalTasks = totalTasks - 1

	if totalTasks == 0 {
	  close(steps.ReadyTasks)
	  close(errTask)
	}
      }
    }(task)
  }

  return results, err
}

func setArgs(args []reflect.Value, body []byte, tS reflect.Type) (reflect.Value, error) {
  var (
    err	      error
    numFields int
    s	      reflect.Value
  )

  switch tS.Kind() {
  case reflect.Ptr:
    s = reflect.New(tS.Elem())
    numFields = tS.Elem().NumField()

    if err = json.Unmarshal(body, s.Interface()); err != nil {
      return s, err
    }
  case reflect.Struct:
    s = reflect.New(tS).Elem()
    numFields = tS.NumField()

    if err = json.Unmarshal(body, &s); err != nil {
      return s, err
    }
  }

  for _, arg := range args {
    arg = ptr(arg)

    for i := 0; i < arg.NumField(); i++ {
      for j := 0; j < numFields; j++ {
	switch tS.Kind() {
	case reflect.Ptr:
	  if tS.Elem().Field(j).Name == arg.Type().Field(i).Name && tS.Elem().Field(j).Type == arg.Type().Field(i).Type {
	    s.Elem().FieldByName(tS.Elem().Field(j).Name).Set(arg.Field(i))
	  }

	  if tagName(arg.Type().Field(i).Tag, tS.Elem().Field(j).Tag) && tS.Elem().Field(j).Type == arg.Type().Field(i).Type {
	    s.Elem().FieldByName(tS.Elem().Field(j).Name).Set(arg.Field(i))
	  }
	case reflect.Struct:
	  if tS.Field(j).Name == arg.Type().Field(i).Name && tS.Field(j).Type == arg.Type().Field(i).Type {
	    s.FieldByName(tS.Field(j).Name).Set(arg.Field(i))
	  }

	  if tagName(arg.Type().Field(i).Tag, tS.Field(j).Tag) && tS.Field(j).Type == arg.Type().Field(i).Type {
	    s.FieldByName(tS.Field(j).Name).Set(arg.Field(i))
	  }
	}
      }
    }
  }

  return s, nil
}

func ptr(arg reflect.Value) reflect.Value {
  if arg.Kind() == reflect.Ptr {
    return ptr(arg.Elem())
  }

  return arg
}

func tagName(tag1, tag2 reflect.StructTag) bool {
  var (
    valueA  string
    valueB  string
    ok      bool
    rep     = strings.NewReplacer(",", "", "json", "", "omitempty", "")
  )

  if valueA, ok = tag1.Lookup("json"); !ok {
    return false
  }

  if valueB, ok = tag2.Lookup("json"); !ok {
    return false
  }

  if rep.Replace(valueA) == "" || rep.Replace(valueB) == "" {
    return false
  }

  if rep.Replace(valueA) == rep.Replace(valueB) {
    return true
  }

  return false
}
