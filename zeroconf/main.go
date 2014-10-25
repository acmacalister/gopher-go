package main

import (
	"fmt"
	"github.com/andrewtj/dnssd"
	"log"
	"net"
	"net/http"
)

func main() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Printf("Listen failed: %s", err)
		return
	}
	port := listener.Addr().(*net.TCPAddr).Port

	op, err := dnssd.StartRegisterOp("", "_airplay._tcp", port, RegisterCallbackFunc)
	if err != nil {
		log.Printf("Failed to register service: %s", err)
		return
	}

	var hardwareAddr net.HardwareAddr
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
	}
	for _, inter := range interfaces {
		hardwareAddr = inter.HardwareAddr
		fmt.Println(inter.HardwareAddr)
	}
	op.SetTXTPair("deviceid", hardwareAddr.String())
	op.SetTXTPair("features", fmt.Sprintf("0x%x", 0x7))
	op.SetTXTPair("model", "AppleTV2,1")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s", r.RemoteAddr)
	})
	http.Serve(listener, nil)

	log.Println("hi")

	// later...
	op.Stop()
}

func RegisterCallbackFunc(op *dnssd.RegisterOp, err error, add bool, name, serviceType, domain string) {
	if err != nil {
		// op is now inactive
		log.Printf("Service registration failed: %s", err)
		return
	}
	if add {
		log.Printf("Service registered as “%s“ in %s", name, domain)
	} else {
		log.Printf("Service “%s” removed from %s", name, domain)
	}
}
