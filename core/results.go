package core

import (
  "strings"
  "reflect"
  "encoding/json"
)

type Results map[string][]reflect.Value

func ConvertResults(res Results) (map[string]string, error) {
  var (
    body  []byte
    err	  error
    r	  = make(map[string]string, len(res))
  )

  for key, value := range res {
    for _, s := range value {
      if body, err = json.Marshal(s.Interface()); err != nil {
	return r, err
      }

      r[key] = string(body)
    }
  }

  return r, nil
}

func ConvertMapToString(m map[string]string) string {
  var (
    body  []byte
    err	  error
  )

  if body, err = json.Marshal(m); err != nil {
    return err.Error()
  }

  return string(body)
}

func convertArgsToString(r []reflect.Value) string {
  var (
    body  []byte
    err	  error
    args  []string
  )

  for _, value := range r {
    if body, err = json.Marshal(value.Interface()); err != nil {
      return err.Error()
    }

    args = append(args, string(body))
  }

  return strings.Join(args, ",")
}
