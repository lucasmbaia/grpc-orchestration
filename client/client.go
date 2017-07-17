package client

import (
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "github.com/lucasmbaia/grpc-orchestration/proto"
  "github.com/lucasmbaia/grpc-base/config"
  "google.golang.org/grpc/credentials"
)

type Config struct {
  SSL bool
}

func init() {
  config.LoadConfig()
}

func (c Config) CallTask(t *orchestration.Task) (*orchestration.Result, error) {
  var (
    conn    *grpc.ClientConn
    sc	    orchestration.OrchestrationServiceClient
    err	    error
    result  *orchestration.Result
  )

  if conn, err = c.connect(); err != nil {
    return result, err
  }

  defer conn.Close()

  sc = orchestration.NewOrchestrationServiceClient(conn)

  if result, err = sc.CallTask(context.Background(), t); err != nil {
    return result, err
  }

  return result, nil
}

func (c Config) connect() (*grpc.ClientConn, error) {
  var (
    opts  []grpc.DialOption
    creds credentials.TransportCredentials
    err	  error
  )

  if config.EnvConfig.GrpcSSL {
    if creds, err = credentials.NewClientTLSFromFile(config.EnvConfig.CAFile, config.EnvConfig.ServerNameAuthority); err != nil {
      return new(grpc.ClientConn), err
    }

    opts = []grpc.DialOption{
      grpc.WithTransportCredentials(creds),
    }
  } else {
    opts = []grpc.DialOption{
      grpc.WithInsecure(),
    }
  }

  return grpc.Dial(config.EnvLocal.LinkerdURL, opts...)
}
