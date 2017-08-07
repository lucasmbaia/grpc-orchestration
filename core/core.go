package core

import (
  "encoding/json"
  "reflect"
  "strings"
  "runtime"
  "path/filepath"

  "github.com/lucasmbaia/grpc-orchestration/tasks"
  "github.com/opentracing/opentracing-go"
  zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

func RunWorkflow(workflow tasks.Workflow, wm *tasks.WorkflowManager) (Results, error) {
  var (
    results Results
    err	    error
    args    []reflect.Value
  )

  if results, err = runWorkflow(workflow, wm, false, args); err != nil {
    go rollback(workflow, results)
    //if workflow.Rollback.Name != "" {
      //go rollback(workflow.Rollback.Name, workflow.Rollback.Step, results)
    //}
  }

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
    collector	zipkin.Collector
    tracer	opentracing.Tracer
    parent	opentracing.Span
  )

  if collector, err = zipkin.NewHTTPCollector("http://172.16.95.113:9411/api/v1/spans"); err != nil {
    return results, err
  }

  if tracer, err = zipkin.NewTracer(zipkin.NewRecorder(collector, false, "127.0.0.1:0", "orchestration")); err != nil {
    return results, err
  }

  parent = tracer.StartSpan("Parent")
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
	child	opentracing.Span
      )

      defer steps.Mutex.Unlock()

      rName = runtime.FuncForPC(reflect.ValueOf(tasks.TR[st.Task].Task).Pointer()).Name()
      fName = rep.Replace(filepath.Ext(rName))

      parent.LogEvent("Start goroutine to exec function " + fName)
      child = tracer.StartSpan("Function name: " + fName, opentracing.ChildOf(parent.Context()))
      child.SetTag("FN", st.FN)

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
	ref = reflect.New(tasks.TR[st.Task].StructTask)

	if len(body) > 0 {
	  if err = json.Unmarshal(body, ref.Interface()); err != nil {
	    for i := 0; i < size; i++ {
	      errTask <-err
	    }

	    closeTasks(steps)
	    child.Finish()
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
	      child.Finish()
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

      child.Finish()
    }(task)
  }

  parent.LogEvent("Finished Task")
  parent.Finish()

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
