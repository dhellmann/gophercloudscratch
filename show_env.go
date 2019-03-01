package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	val := os.Getenv("IRONIC_URL")
	if val == "" {
		val = "ipmi://192.168.122.1:6233"
	}
	fmt.Printf("got '%s'\n", val)

	u, err := url.Parse(val)
	if err != nil {
		panic(err)
	}
	fmt.Printf("scheme: %s\n", u.Scheme)
	fmt.Printf("host: %s\n", u.Host)
	fmt.Printf("host name: %s\n", u.Hostname())
	fmt.Printf("port: %s\n", u.Port())
	fmt.Printf("path: %v\n", u.Path)
	fmt.Printf("query: %v\n", u.Query())
	fmt.Printf("url: %s\n", u)
}
