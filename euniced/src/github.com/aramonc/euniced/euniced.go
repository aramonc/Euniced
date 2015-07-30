package main

import "os"
import "os/signal"
import "syscall"
import "flag"
import "fmt"
import "github.com/aramonc/config"

// Define channels
var done = make(chan bool)
var conf = make(chan config.Config, 1)
var sigs = make(chan os.Signal, 1)

func main() {
	confPath := flag.String("conf", "/etc/euniced.conf.yml", "Path to the YAML configuration file to use")
	flag.Parse()

	if _, err := os.Stat(*confPath); os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go sigHandler(sigs)

	go config.Parse(*confPath, conf)
	go initialize(conf)

	<-done
}

func initialize(c <-chan config.Config) {
	conf := <-c
	fmt.Println(conf.Rabbit.Host)
}

func sigHandler(signals <-chan os.Signal) {
	sig := <-signals
	if sig.String() == "interrupt" || sig.String() == "terminated" || sig.String() == "Killed" {
		done <- true
	}
}
