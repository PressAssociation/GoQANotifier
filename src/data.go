package main

import (
    "time"
)

type Ping struct {
    ID    string      `json:"service,omitempty"`
    Time  time.Time   `json:"last,omitempty"`
    Passed bool    `json:"passed,omitempty"`
}

type Service struct {
    ID    string      `json:"service,omitempty"`
    Env   string      `json:"env,omitempty"`
    LastPass *Ping      `json:"lastpass,omitempty"`
    LastFail *Ping      `json:"lastfail,omitempty"`
}

var services []Service

func setupServices() {
    services = make([]Service, 0)
}

func getServices() []Service {
  return services
}

func getService(id string, env string) *Service {

    for i := 0; i < len(getServices()); i++ {
      attr := &getServices()[i]
      if attr.ID == id && attr.Env == env {
        return attr
      }
    }

    return new(Service)
}

func createPing(id string, passed bool) *Ping {
    var ping Ping
    ping.ID = id
    ping.Time = time.Now()
    ping.Passed = passed
    return &ping;
}

func createService(id string, env string) *Service {

    service := Service{ID: id, Env: env}
    services = append(services, service)

    return &service
}

func updateService(service *Service, ping *Ping, hasPassedTest bool) {
    if(hasPassedTest){
      service.LastPass = ping;
    } else {
      service.LastFail = ping;
    }
}
