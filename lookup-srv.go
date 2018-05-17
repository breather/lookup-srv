// An utility to lookup SRV records
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var (
	shellScript = flag.String("script", "", "Script to run with result.")
	srvName = flag.String("srv", "", "Service responsible for the DNS records of.")
	dspPort = flag.Bool("port", false, "Enable displaying port number.")
)

func lookup(svcName string, withPort bool) ([]string, error) {
	var endpoints []string
	_, srvRecords, err := net.LookupSRV("", "", svcName)
	if err != nil {
		return endpoints, err
	}
	for _, srvRecord := range srvRecords {
		// The SRV records ends in a "." for the root domain
		host := srvRecord.Target[:len(srvRecord.Target)-1]
		if withPort == true {
			ep := fmt.Sprintf("%v:%v", host, srvRecord.Port)
			endpoints = append(endpoints, ep)
		} else {
			ep := fmt.Sprintf("%v", host)
			endpoints = append(endpoints, ep)
		}
	}
	return endpoints, nil
}

func shellOut(input, script string) {
	out, err := exec.Command("sh", "-c", fmt.Sprintf("echo '%v' | %v", input, script)).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute %v: %v, err: %v", script, string(out), err)
	}
	fmt.Print(string(out))
}

func main() {
	flag.Parse()

	service := *srvName
	if service == "" {
		service = os.Getenv("SRV_NAME")
	}
	if service == "" {
		log.Fatalf("Incomplete args, require -srv or SRV_NAME environment variable")
	}

	withPort := *dspPort
	peers, err := lookup(service, withPort)
	if err != nil {
		log.Fatalf("%v", err)
	}
	sort.Strings(peers)

	script := *shellScript
	if script == "" {
		fmt.Printf("%v\n", strings.Join(peers, "\n"))
	} else {
		shellOut(strings.Join(peers, "\n"), script)
	}
}
