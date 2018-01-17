package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/xml"
	"net/http"
	"github.com/geekfghuang/snowflake"
	"encoding/json"
)

type Conf struct {
	WorkerId int64 `xml:"workerId"`
	Port string `xml:"port"`
	Path string `xml:"path"`
}

// Code 200:success, -1:error
type HttpResponse struct {
	Code int
	Msg string
	Id int64
}

var worker *snowflake.Worker

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	resp := new(HttpResponse)
	id, err := worker.NextId()
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
	} else {
		resp.Code = 200
		resp.Msg = "OK"
	}
	resp.Id = id

	reply, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("error json marshal: %v\n", err)
		os.Exit(1)
	}
	w.Header().Set("Content-Type","application/json")
	w.Write(reply)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("worker.xml is expected, usage:./uid-http ../conf/worker.xml")
		os.Exit(1)
	}
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("error read worker.xml: %v\n", err)
		os.Exit(1)
	}
	var conf Conf
	err = xml.Unmarshal(content, &conf)
	if err != nil {
		fmt.Printf("error parse worker.xml: %v\n", err)
		os.Exit(1)
	}

	worker, err = snowflake.NewWorker(conf.WorkerId)
	if err != nil {
		fmt.Printf("error build worker: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc(conf.Path, serveHTTP)
	err = http.ListenAndServe(":" + conf.Port, nil)
	if err != nil {
		fmt.Printf("error start service: %v\n", err)
		os.Exit(1)
	}
}