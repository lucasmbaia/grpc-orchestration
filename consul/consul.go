package consul

import (
  "os"
  "strings"
  "strconv"
  "net"

  _consul "github.com/hashicorp/consul/api"
  "github.com/lucasmbaia/grpc-orchestration/config"
)

const (
  URL_CHECK = "http://{url_check}:{port_check}/{endpoint_check}"
)

type ConsulConfig struct {
  ConsulURL             string
  ServiceName           string
  ServicePort           int
  ServiceIntervalCheck  int
  ServiceTimeout        int
  ServiceDeregister     int
  PortUrlCheck          string
  EndPointCheck         string
}

type consulConfig struct {
  ConsulURL	string
  ID            string
  Name          string
  Address       string
  UrlCheck      string
  Port          int
  IntervalCheck int
  Timeout       int
  Deregister    int
}

func RegisterService() error {
  var (
    err	      error
    ips	      []string
    urlCheck  *strings.Replacer
    hostname  string
    service   consulConfig
  )

  if ips, err = getIPs(); err != nil {
    return err
  }

  if hostname, err = os.Hostname(); err != nil {
    return err
  }

  urlCheck = strings.NewReplacer("{url_check}", ips[0], "{port_check}", config.EnvConfig.PortUrlCheck, "{endpoint_check}", config.EnvConfig.EndPointCheck)

  service = consulConfig{
    ConsulURL:	    config.EnvConfig.ConsulURL,
    ID:		    hostname,
    Name:	    config.EnvConfig.ServiceName,
    Address:	    ips[0],
    Port:	    config.EnvConfig.ServicePort,
    UrlCheck:	    urlCheck.Replace(URL_CHECK),
    IntervalCheck:  config.EnvConfig.ServiceIntervalCheck,
    Timeout:	    config.EnvConfig.ServiceTimeout,
    Deregister:	    config.EnvConfig.ServiceDeregister,
  }

  return register(service)
}

func initConsul(host string) (*_consul.Client, error) {
  var (
    config  = _consul.DefaultConfig()
  )

  config.Address = host

  return _consul.NewClient(config)
}

func register(c consulConfig) error {
  var (
    client  *_consul.Client
    err	    error
    service _consul.AgentServiceRegistration
  )

  if client, err = initConsul(c.ConsulURL); err != nil {
    return err
  }

  service = _consul.AgentServiceRegistration{
    ID:	      c.ID,
    Name:     c.Name,
    Port:     c.Port,
    Address:  c.Address,
    Check:    &_consul.AgentServiceCheck{
      HTTP:			      c.UrlCheck,
      Interval:			      strconv.Itoa(c.IntervalCheck) + "s",
      Timeout:			      strconv.Itoa(c.Timeout) + "s",
      DeregisterCriticalServiceAfter: strconv.Itoa(c.Deregister) + "s",
    },
  }

  return client.Agent().ServiceRegister(&service)
}

func getIPs() ([]string, error) {
  var (
    addrs []net.Addr
    err   error
    ips   []string
  )

  if addrs, err = net.InterfaceAddrs(); err != nil {
    return ips, err
  }

  for _, addr := range addrs {
    if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
      if ip.IP.To4() != nil {
	ips = append(ips, ip.IP.String())
      }
    }
  }

  return ips, nil
}
