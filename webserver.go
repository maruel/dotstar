// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
)

type WebServer struct {
	painter *Painter
}

func StartWebServer(painter *Painter, port int) *WebServer {
	ws := &WebServer{painter: painter}
	mux := http.NewServeMux()
	mux.HandleFunc("/", ws.rootHandler)
	mux.HandleFunc("/switch", ws.switchHandler)
	mux.HandleFunc("/favicon.ico", ws.faviconHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), loggingHandler{mux})
	return ws
}

func (s *WebServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write(read("root.html")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *WebServer) switchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	if n := r.PostFormValue("new"); len(n) != 0 {
		if p := Patterns[n]; p != nil {
			s.painter.SetPattern(p)
		}
	}
}

func (s *WebServer) faviconHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "Cache-Control:public, max-age=2592000") // 30d
	w.Write(read("dotstar.png"))
}

// Private details.

type loggingHandler struct {
	handler http.Handler
}

type loggingResponseWriter struct {
	http.ResponseWriter
	length int
	status int
}

func (l *loggingResponseWriter) Write(data []byte) (size int, err error) {
	size, err = l.ResponseWriter.Write(data)
	l.length += size
	return
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.ResponseWriter.WriteHeader(status)
	l.status = status
}

// ServeHTTP logs each HTTP request if -verbose is passed.
func (l loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrw := &loggingResponseWriter{ResponseWriter: w}
	l.handler.ServeHTTP(lrw, r)
	log.Printf("%s - %3d %6db %4s %s\n", r.RemoteAddr, lrw.status, lrw.length, r.Method, r.RequestURI)
}
