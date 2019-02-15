package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "fmt"
    "github.com/fogleman/gg"
    "bytes"
    "strconv"
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

func main() {
    services = make([]Service, 0)

    router := mux.NewRouter()
    router.HandleFunc("/ping", GetAll).Methods("GET")
    router.HandleFunc("/ping/{env}", GetEnv).Methods("GET")
    router.HandleFunc("/ping/{env}/{id}.png", GetImage).Methods("GET")
    router.HandleFunc("/ping/{env}/{id}", GetPing).Methods("GET")
    router.HandleFunc("/ping/{env}/{id}", CreatePing).Methods("POST")
    log.Fatal(http.ListenAndServe(":8000", router))
}

func GetAll(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(services)
}

func GetEnv(w http.ResponseWriter, r *http.Request) {
    envservices := make([]Service, 0)

    params := mux.Vars(r)
    for _, item := range services {
        if item.Env == params["env"] {
          envservices = append(envservices, item)
        }
    }

    json.NewEncoder(w).Encode(envservices)
}


func GetPing(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range services {
        if item.ID == params["id"] && item.Env == params["env"] {
           json.NewEncoder(w).Encode(item)
        }
    }
}

func GetImage(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range services {
        if item.ID == params["id"] && item.Env == params["env"] {
             const X = 320 
             const Y = 30

             dc := gg.NewContext(X, Y)
             dc.SetRGB(1, 1, 1)
             dc.Clear()
             dc.SetRGB(0, 0, 0)

             if (item.LastPass != nil) {
             lastPass, _ := json.Marshal(item.LastPass.Time)
             dc.DrawStringAnchored(string(item.Env) + " Last Pass: " + string(lastPass), 2, Y/4, 0, 0.5)
             }

             if (item.LastFail != nil) {
             lastFail, _ := json.Marshal(item.LastFail.Time)
             dc.DrawStringAnchored(string(item.Env) + " Last Fail: " + string(lastFail), 2, Y/4 * 3, 0, 0.5)
             }

             buffer := new(bytes.Buffer)
             if err := dc.EncodePNG(buffer); err != nil {
                 log.Println("unable to encode image.")
             }

             w.Header().Set("Content-Type", "image/jpeg")
             w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
             if _, err := w.Write(buffer.Bytes()); err != nil {
                 log.Println("unable to write image.")
             }
        }
    }
}


type jsonbody struct {
    Passed bool `json:"passed"`
}

func CreatePing(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    decoder := json.NewDecoder(r.Body) 
    var t jsonbody
    err := decoder.Decode(&t)
    if err != nil {
        //Need a body
        fmt.Println("Error decoding json.")
        return
    }

    var ping Ping
    ping.ID = params["id"]
    ping.Time = time.Now()
    ping.Passed = t.Passed
    fmt.Println(ping);

    for i := 0; i < len(services); i++ {
      attr := &services[i]
      if attr.ID == params["id"] && attr.Env == params["env"] {
        fmt.Println("Previously existing service...")
        if(t.Passed){
          attr.LastPass = &ping;
        } else {
          attr.LastFail = &ping;
        }
        json.NewEncoder(w).Encode(attr)
        return
      }
    }

    fmt.Println("New service...")
    service := Service{ID: params["id"], Env: params["env"]}
    if(t.Passed){
      service.LastPass = &ping;
    } else {
      service.LastFail = &ping;
    }
    services = append(services, service)
    json.NewEncoder(w).Encode(service)
}
