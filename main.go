package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/sahilchug/loadbalancer/config"
	viper "github.com/spf13/viper.git"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var listeners config.Listeners

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&listeners)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	defaultHTTPHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s %s\n", req.RemoteAddr, req.Method, req.URL)
		io.WriteString(w, "Hello, world!\n")
	}

	var wg sync.WaitGroup
	for _, x := range listeners.Listeners {
		fmt.Println(x.Protocol)
		fmt.Println(x.Port)
		if x.Protocol == "http" {
			wg.Add(1)
			go func(x config.Listener) {
				httpServer := http.NewServeMux()
				fmt.Println("starting server")
				httpServer.HandleFunc("/", defaultHTTPHandler)
				addr := fmt.Sprintf("127.0.0.1:%d", x.Port)
				fmt.Printf("Serving Server %s", addr)
				log.Fatal(http.ListenAndServe(addr, httpServer))
			}(x)
		}
	}
	wg.Wait()
}
