package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	data map[string]string
	lock *sync.RWMutex
)

func init() {
	data = map[string]string{
		"/latest/meta-data/ami-id":                      "ami-2bb65342",
		"/latest/meta-data/ami-launch-index":            "",
		"/latest/meta-data/ami-manifest-path":           "",
		"/latest/meta-data/hostname":                    "i-10a64382",
		"/latest/meta-data/instance-action":             "",
		"/latest/meta-data/instance-id":                 "i-10a64382",
		"/latest/meta-data/instance-type":               "t2.small",
		"/latest/meta-data/kernel-id":                   "",
		"/latest/meta-data/local-hostname":              "ip-10-251-50-38.ec2.internal",
		"/latest/meta-data/local-ipv4":                  "10.251.50.38",
		"/latest/meta-data/mac":                         "02:29:96:8f:6a:2d",
		"/latest/meta-data/placement/availability-zone": "eu-west-1a",
		"/latest/meta-data/public-hostname":             "ec2-203-0-113-25.compute-1.amazonaws.com",
		"/latest/meta-data/public-ipv4":                 "ec2-203-0-113-25.compute-1.amazonaws.com",
		"/latest/meta-data/reservation-id":              "r-fea54097",
		"/latest/meta-data/security-groups":             "",
		// block-device-mapping/
		// network/
		// public-keys/
		// services/
		//"/latest/user-data/": "",
	}
	lock = &sync.RWMutex{}
}

func main() {
	h := router()
	s := &http.Server{
		Addr:         "169.254.169.254:80",
		Handler:      h,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	log.Fatal(s.ListenAndServe())
}

func router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	return mux
}

// see http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getHandler(w, r)
	case "HEAD":
		getHandler(w, r)
	case "PUT":
		setHandler(w, r)
	case "POST":
		setHandler(w, r)
	case "DELETE":
		deleteHandler(w, r)
	default:
		defaultHandler(w, r)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()

	if val, found := data[r.URL.Path]; found {
		fmt.Fprintf(w, val)
		return
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	val, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, found := data[r.URL.Path]
	data[r.URL.Path] = string(val)
	if found {
		http.Error(w, "Update", http.StatusAccepted)
		return
	}
	http.Error(w, "Created", http.StatusCreated)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	if _, found := data[r.URL.Path]; found {
		delete(data, r.URL.Path)
		http.Error(w, "Deleted", http.StatusNoContent)
		return
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
