package tasks

import (
  "encoding/json"
  "io/ioutil"
  "errors"
  "path"
  "sync"
  "os"

  "github.com/lucasmbaia/go-environment/etcd"
  //"gopkg.in/go-playground/validator.v9"
)

type Workflow struct {
  Name		  string      `json:",omitempty" validate:"required"`
  Description	  string      `json:",omitempty" validate:"required"`
  Version	  int	      `json:",omitempty" validate:"required"`
  InputParameters string      `json:",omitempty" validate:"required"`
  Tasks		  []TasksFlow `json:",omitempty"`
}

type TasksFlow struct {
  Name		string	  `json:",omitempty" validate:"required"`
  TaskReference	string	  `json:",omitempty" validate:"required"`
  Dependency	[]string  `json:",omitempty" validate:"existsTaskReference"`
  Rollback	struct {
    Name  string  `json:",omitempty"`
    Step  string  `json:",omitempty"`
  }
}

type WorkflowManager struct {
  Signal  chan os.Signal
  Stop	  chan bool
  Restart chan bool
}

type ConfigWorkflow struct {
  EtcdURL string
  Keys	  []string
  Dir	  string
}

type WM map[string]*WorkflowManager

var (
  WManager	      = make(map[string]*WorkflowManager)
  WorkflowsRegistred  = make(map[string]Workflow)
  Mutex		      = &sync.Mutex{}
)

func (c ConfigWorkflow) RegisterWorkflows() error {
  var (
    files     []os.FileInfo
    body      []byte
    workflow  Workflow
    err	      error
  )

  if err = c.validConfigWorkflow(); err != nil {
    return err
  }

  if c.EtcdURL != "" {
    var (
      config	etcd.Config
      client	etcd.Client
    )

    config = etcd.Config {
      Endpoints:  []string{c.EtcdURL},
      TimeOut:	  5,
    }

    if client, err = config.NewClient(); err != nil {
      return err
    }

    for _, key := range c.Keys {
      var workflow Workflow

      if err = client.Get(key, &workflow, false, false); err != nil {
	return err
      }

      WorkflowsRegistred[workflow.Name] = workflow
    }
  } else {
    if c.Dir[len(c.Dir)-1:] != "/" {
      c.Dir = c.Dir + "/"
    }

    if files, err = ioutil.ReadDir(c.Dir); err != nil {
      return err
    }

    for _, file := range files {
      if path.Ext(file.Name()) == ".json" {
	if body, err = ioutil.ReadFile(c.Dir + file.Name()); err != nil {
	  return err
	}

	if workflow, err = unmarshalWorkflow(body); err != nil {
	  return err
	}

	WorkflowsRegistred[workflow.Name] = workflow
      }
    }
  }

  return nil
}

func (c ConfigWorkflow) validConfigWorkflow() error {
  if c.EtcdURL == "" && c.Dir == "" {
    return errors.New("Please inform the records directory or etcd url")
  }

  if c.EtcdURL != "" && len(c.Keys) == 0 {
    return errors.New("Please inform keys the etcd to get workflows")
  }

  return nil
}

func GetWorkflow(name string) (Workflow, error) {
  if _, ok := WorkflowsRegistred[name]; ok {
    return WorkflowsRegistred[name], nil
  }

  return Workflow{}, errors.New("Workflow not registered")
}

func RegisterWM(key string) *WorkflowManager {
  var (
    wm = &WorkflowManager{Signal: make(chan os.Signal), Stop: make(chan bool, 1), Restart: make(chan bool, 1)}
  )

  Mutex.Lock()
  WManager[key] = wm
  Mutex.Unlock()

  return wm
}

func DeleteWM(key string) {
  Mutex.Lock()

  if _, ok := WManager[key]; ok {
    delete(WManager, key)
  }

  Mutex.Unlock()
  return
}

func GetTasksWorkflow(tasks []TasksFlow, step string) (int, bool) {
  for index, task := range tasks {
    if task.Name == step {
      return index, true
    }
  }

  return 0, false
}

func unmarshalWorkflow(body []byte) (Workflow, error) {
  var (
    workflow  Workflow
    err       error
  )

  if err = json.Unmarshal(body, &workflow); err != nil {
    return workflow, err
  }

  return workflow, nil
}

/*func validateWorkflow(body []byte) (Workflow, error) {
  var (
    workflow  Workflow
    err	      error
    validate  = validator.New()
  )

  if workflow, err = unmarshalWorkflow(body); err != nil {
    return workflow, err
  }

  if err = validate.Struct(workflow); err != nil {
    log.Println(err)
  }

  return workflow, nil
}*/
