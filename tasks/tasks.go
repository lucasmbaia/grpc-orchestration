package tasks

import (
  "reflect"
  "fmt"
)

type TasksReferences struct {
  Task	      interface{}
  StructTask  reflect.Type
}

type TasksReferencesInfos map[string]TasksReferences

var (
  TR  = make(TasksReferencesInfos)
)

func RegisterTasksReferences(tr TasksReferencesInfos) error {
  var (
    f reflect.Value
    t reflect.Type
  )

  for key, task := range tr {
    f = reflect.Indirect(reflect.ValueOf(task.Task))
    t = f.Type()

    if f.Kind() != reflect.Func {
      return fmt.Errorf("%T must be a function ", f)
    }

    for i := 0; i < t.NumIn(); i++ {
      if t.In(i).Kind() != reflect.Ptr && t.In(i).Kind() != reflect.Struct {
	return fmt.Errorf("%T Only accepts ptr or struct as parameters ", f)
      }
    }

    TR[key] = task
  }
  return nil
}
