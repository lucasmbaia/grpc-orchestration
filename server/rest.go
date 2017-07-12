package server

import (
  "flag"
  "net/http"
  "strconv"

  "google.golang.org/grpc/credentials"
  "github.com/lucasmbaia/grpc-orchestration/config"
  "github.com/lucasmbaia/grpc-orchestration/proto"
  "golang.org/x/net/context"
  "github.com/grpc-ecosystem/grpc-gateway/runtime"
  "google.golang.org/grpc"
)

func registerGateway(ctx context.Context, ssl bool, opts ...runtime.ServeMuxOption) (*runtime.ServeMux, error) {
  var (
    mux       = runtime.NewServeMux(opts...)
    err       error
    endpoint  = flag.String("endpoint", "localhost" + ":" + strconv.Itoa(config.EnvConfig.ServicePort), "endpoints")
    creds     credentials.TransportCredentials
    dialOpts  []grpc.DialOption
  )

  if ssl {
    if creds, err = credentials.NewClientTLSFromFile(config.EnvConfig.CAFile, config.EnvConfig.ServerNameAuthority); err != nil {
      return nil, err
    }

    dialOpts = []grpc.DialOption{
      grpc.WithTransportCredentials(creds),
    }
  } else {
    dialOpts = []grpc.DialOption{
      grpc.WithInsecure(),
    }
  }

  if err = orchestration.RegisterOrchestrationServiceHandlerFromEndpoint(ctx, mux, *endpoint, dialOpts); err != nil {
    return nil, err
  }

  return mux, nil
}

func setHandler(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    h.ServeHTTP(w, r)
  })
}

func gateway(ssl bool) error {
  var (
    mux     *http.ServeMux
    gw      *runtime.ServeMux
    ctx     = context.Background()
    cancel  context.CancelFunc
    opts    []runtime.ServeMuxOption
    err     error
  )

  mux = http.NewServeMux()

  ctx, cancel = context.WithCancel(ctx)
  defer cancel()

  if gw, err = registerGateway(ctx, ssl, opts...); err != nil {
    return err
  }

  mux.Handle("/", gw)

  return http.ListenAndServe(":" + config.EnvConfig.PortUrlCheck, setHandler(mux))
}
