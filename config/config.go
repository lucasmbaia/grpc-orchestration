package config

import (
  "log"
  "github.com/lucasmbaia/go-environment/local"
  "github.com/lucasmbaia/go-environment/etcd"
)

var (
  EnvConfig   Config
  EnvLocal  Service
  parsed    = false
)

type Config struct {
  ConsulURL             string	  `env:"CONSUL_URL" envDefault:"127.0.0.1:8500"`
  ServiceName           string	  `env:"SERVICE_NAME" envDefault:"grpc-orchestration"`
  ServicePort           int	  `env:"SERVICE_PORT" envDefault:"5000"`
  TypeConnection	string	  `env:"TYPE_CONNECTION" envDefault:"tcp"`
  ServiceIntervalCheck  int	  `env:"SERVICE_INTERVAL_CHECK" envDefault:"5"`
  ServiceTimeout        int	  `env:"SERVICE_TIMEOUT" envDefault:"1"`
  ServiceDeregister     int	  `env:"SERVICE_DEREGISTER" envDefault:"10"`
  PortUrlCheck          string	  `env:"PORT_URL_CHECK" envDefault:"8080"`
  EndPointCheck         string	  `env:"ENDPOINT_CHECK" envDefault:"v1/health"`
  CertFile              string	  `env:"CERT_FILE" envDefault:""`
  KeyFile               string	  `env:"KEY_FILE" envDefault:""`
  CAFile                string	  `env:"CA_FILE" envDefault:""`
  ServerNameAuthority   string	  `env:"SERVER_NAME_AUTHORITY" envDefault:""`
  WorkflowsName		[]string  `env:"WORKFLOWS_NAME" envDefault:""`
  GrpcSSL		bool	  `env:"GRPC_SSL" envDefault:""`
}

type Service struct {
  ServiceName         string  `env:"SERVICE_NAME" envDefault:"grpc-orchestration"`
  EtcdURL             string  `env:"ETCD_URL" envDefault:"http://127.0.0.1:2379"`
  LinkerdURL          string  `env:"LINKERD_URL" envDefault:"127.0.0.1:4140"`
  CAFile              string  `env:"CA_FILE" envDefault:""`
  ServerNameAuthority string  `env:"SERVER_NAME_AUTHORITY" envDefault:""`
}

func LoadConfig() {
  var (
    config	etcd.Config
    client	etcd.Client
    err		error
  )

  if !parsed {
    if err = local.Get("", &EnvLocal, true, false); err != nil {
      log.Fatalf("Error to get local env: ", err)
    }

    config = etcd.Config {
      Endpoints:  []string{EnvLocal.EtcdURL},
      TimeOut:	  5,
    }

    if client, err = config.NewClient(); err != nil {
      log.Fatalf("Error to connect etcd: ", err)
    }

    if err = client.Get(EnvLocal.ServiceName, &EnvConfig, true, false); err != nil {
      log.Fatalf("Error to get etcd env: ", err)
    }

    parsed = true
  }
}

func LoadLocalEnv() (Service, error) {
  var (
    err	    error
    service Service
  )

  if err = local.Get("", &service, true, false); err != nil {
    return service, err
  }

  return service, nil
}
