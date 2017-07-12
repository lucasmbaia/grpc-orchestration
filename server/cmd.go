package server

import (
  "strconv"
  "net"

  "github.com/lucasmbaia/grpc-orchestration/config"
  "github.com/lucasmbaia/grpc-orchestration/proto"
  "github.com/lucasmbaia/grpc-orchestration/consul"
  "google.golang.org/grpc/credentials"
  "google.golang.org/grpc/reflection"
  "google.golang.org/grpc"
)

type ConfigCMD struct {
  SSL		  bool
  RegisterConsul  bool
}

func (c ConfigCMD) Run() error {
  var (
    listen		net.Listener
    err			error
    errChan		= make(chan error, 1)
    creds		credentials.TransportCredentials
    opts		[]grpc.ServerOption
    s			*grpc.Server
    orchestrationServer = NewOrchestrationServer()
  )

  go func() {
    if listen, err = net.Listen(config.EnvConfig.TypeConnection, ":" + strconv.Itoa(config.EnvConfig.ServicePort)); err != nil {
      errChan <- err
      return
    }

    if c.SSL {
      if creds, err = credentials.NewServerTLSFromFile(config.EnvConfig.CertFile, config.EnvConfig.KeyFile); err != nil {
	errChan <- err
	return
      }

      opts = []grpc.ServerOption{
	grpc.Creds(creds),
      }
    }

    if c.RegisterConsul {
      if err = consul.RegisterService(); err != nil {
	errChan <- err
	return
      }
    }

    s = grpc.NewServer(opts...)
    orchestration.RegisterOrchestrationServiceServer(s, orchestrationServer)
    reflection.Register(s)

    errChan <-s.Serve(listen)

    errChan <-gateway(c.SSL)
  }()

  select {
  case e := <-errChan:
    return e
  }
}
