package main

import (
    "time"
    "os"
    "encoding/csv"
    "log"
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
const FILENAME string = "services.csv"

func setupServices() {
    services = make([]Service, 0)
    readFromfile()
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

    return &services[len(services)-1]
}

func updateService(service *Service, ping *Ping, hasPassedTest bool) {
    if(hasPassedTest){
      service.LastPass = ping;
    } else {
      service.LastFail = ping;
    }

    writeToFile()
}


func writeToFile() {
    file, err := os.Create(FILENAME)
    if err != nil {
      log.Print("Cannot create file", err)
      return
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, service := range getServices() {
        lastPass := ""
        lastFail := ""
        if (service.LastPass != nil) {
          lastPass = service.LastPass.Time.Format(time.RFC3339)
        }
        if (service.LastFail != nil) {
          lastFail = service.LastFail.Time.Format(time.RFC3339)
        }

        csvline := []string{service.Env, service.ID, lastPass, lastFail}

        err := writer.Write(csvline)
        if err != nil {
          log.Print("Cannot write to file", err)
          return
        }
    }
}

func readFromfile() {
    // Open CSV file
    f, err := os.Open(FILENAME)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    // Read File into a Variable
    lines, err := csv.NewReader(f).ReadAll()
    if err != nil {
        panic(err)
    }

    // Loop through lines & turn into object
    for _, line := range lines {

        var passPing *Ping
        if(line[2] != "") {
          passPing = createPing(line[1], true)
          passPing.Time, _ = time.Parse(time.RFC3339, line[2])
        }

        var failPing *Ping
        if(line[3] != "") {
          failPing = createPing(line[1], false)
          failPing.Time, _ = time.Parse(time.RFC3339, line[3])
        }

        service := createService(line[1], line[0])
        service.LastPass = passPing
        service.LastFail = failPing
    }

    log.Print("Finished Reading In Data")
}
