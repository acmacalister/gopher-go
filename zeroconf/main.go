package main

import (
	"fmt"
	"github.com/acmacalister/container/stack"
	// "github.com/andrewtj/dnssd"
	// "log"
	// "net"
	// "net/http"
)

func main() {
	q := stack.StringStack{}
	q.Push("Red Shirt")
	q.Push("Kirk")
	q.Push("Spock")
	str, _ := q.Pop()
	fmt.Println(str)
	fmt.Println(q.Empty())
	fmt.Println(q.Size())
}

// 	listener, err := net.Listen("tcp", ":5000")
// 	if err != nil {
// 		log.Printf("Listen failed: %s", err)
// 		return
// 	}
// 	port := listener.Addr().(*net.TCPAddr).Port

// 	op := dnssd.NewRegisterOp("", "_raop._tcp", port, RegisterCallbackFunc)

// 	// var hardwareAddr net.HardwareAddr
// 	// interfaces, err := net.Interfaces()
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// }

// 	// for _, inter := range interfaces {
// 	// 	hardwareAddr = inter.HardwareAddr
// 	// 	fmt.Println(inter.HardwareAddr)
// 	// }

// 	hardwareAddr := []byte{0x48, 0x5d, 0x60, 0x7c, 0xee, 0x22}

// 	op.SetTXTPair("deviceid", string(hardwareAddr))
// 	op.SetTXTPair("features", fmt.Sprintf("0x%x", 0x7))
// 	op.SetTXTPair("model", "AppleTV2,1")
// 	err = op.Start()
// 	if err != nil {
// 		log.Printf("Failed to register service: %s", err)
// 		return
// 	}

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Hello, %s", r.RemoteAddr)
// 	})
// 	http.Serve(listener, nil)

// 	// later...
// 	op.Stop()
// }

// func RegisterCallbackFunc(op *dnssd.RegisterOp, err error, add bool, name, serviceType, domain string) {
// 	if err != nil {
// 		// op is now inactive
// 		log.Printf("Service registration failed: %s", err)
// 		return
// 	}
// 	if add {
// 		log.Printf("Service registered as “%s“ in %s", name, domain)
// 	} else {
// 		log.Printf("Service “%s” removed from %s", name, domain)
// 	}
// }
