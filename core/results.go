package core

import (
  "reflect"
  "encoding/json"
)

type Results map[string][]reflect.Value

func convertResults(res Results) (map[string]string, error) {
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
