// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	//_ "net/http/pprof"
	_ "github.com/mkevac/debugcharts"

	"github.com/panjf2000/ants"
)

var addr = flag.String("addr", ":8080", "http service address")

var poolSize int =1000
func serveHome(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path,err:=os.Getwd()
	if(err!=nil){
		log.Println(err)
	}
	http.ServeFile(w, r, path+"/home.html") // "E:\\temp\\go\\broadcast\\chat\\home.html"
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	pool, _ := ants.NewPool(10000)
	defer ants.Release()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		if pool.Running()-pool.Cap()>100 {
			if poolSize<=12000  {
				poolSize+=1000
				pool.ReSize(poolSize)
			}else {
				w.WriteHeader(503)
				w.Write([]byte("服务器繁忙"))
				return
			}
		}

		serveWs(hub, w, r,pool)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
