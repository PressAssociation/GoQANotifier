package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
    "github.com/fogleman/gg"
    "bytes"
    "strconv"
)

func main() {
    setupServices()

    router := mux.NewRouter()
    router.HandleFunc("/ping", GetAll).Methods("GET")
    router.HandleFunc("/ping/{env}", GetEnv).Methods("GET")
    router.HandleFunc("/ping/{env}/{id}.png", GetImage).Methods("GET")
    router.HandleFunc("/ping/{env}/{id}", GetPing).Methods("GET")
    router.HandleFunc("/ping/{env}/{id}", CreatePing).Methods("POST")
    log.Fatal(http.ListenAndServe(":8000", router))
}

func GetAll(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(getServices())
}

func GetEnv(w http.ResponseWriter, r *http.Request) {
    envservices := make([]Service, 0)

    params := mux.Vars(r)
    for _, item := range getServices() {
        if item.Env == params["env"] {
          envservices = append(envservices, item)
        }
    }

    json.NewEncoder(w).Encode(envservices)
}


func GetPing(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    json.NewEncoder(w).Encode(getService(params["id"], params["env"]));
}

func GetImage(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    item := *getService(params["id"], params["env"])
    if item != (Service{}) {
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

func CreatePing(w http.ResponseWriter, r *http.Request) {
    type jsonbody struct {
        Passed bool `json:"passed"`
    }

    params := mux.Vars(r)
    decoder := json.NewDecoder(r.Body) 
    var t jsonbody
    err := decoder.Decode(&t)
    if err != nil {
        //Need a body
        fmt.Println("Error decoding json.")
        return
    }

    ping := createPing(params["id"], t.Passed)
    attr := getService(params["id"], params["env"]);

    fmt.Print(*ping);
    fmt.Print(" updating " + params["id"] + " " + params["env"] + ". ");

    if *attr != (Service{}) {
      fmt.Println("Previously existing service...")
    } else {
      fmt.Println("New service...")
      attr = createService(params["id"], params["env"])
    }

    updateService(attr, ping, t.Passed)

    json.NewEncoder(w).Encode(attr)
}
