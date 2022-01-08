package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Entering the function healthz ...\n")
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			fmt.Printf("%s=%s\n", k, v[0])
			w.Header().Add(k, v[0])
		}
	}
	version := os.Getenv("VERSION")
	if version != "" {
		w.Header().Add("VERSION", version)
	} else {
		w.Header().Add("VERSION", "None")
	}
	clientIp, err := GetClientIP(r)
	if err == nil {
		fmt.Printf("CLIENT_IP=%s\n", clientIp)
	}
}

func GetClientIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	return "", errors.New("no valid ip found")
}

func healthz(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(200)
}