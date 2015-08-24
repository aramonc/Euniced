package main

import "os"
import "os/signal"
import "os/exec"
import "syscall"
import "flag"
import "fmt"
// import "log"
import "github.com/aramonc/config"
import "strings"
import "bufio"

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
	for _,procConf := range conf.Workers {
		bootWorker(procConf);
	}
}

func sigHandler(signals <-chan os.Signal) {
	sig := <-signals
	if sig.String() == "interrupt" || sig.String() == "terminated" || sig.String() == "Killed" {
		done <- true
	}
}

func bootWorker(workerConf config.Worker) {
	cmd := make([]*exec.Cmd, workerConf.Max, workerConf.Max)
	for i := 0; i < workerConf.Max; i++ {
		cmd[i] = exec.Command(workerConf.Command, strings.Join(workerConf.Arguments, " "))
		attachToLog(cmd[i]);
	}
	

	for j := 0; j < workerConf.Max; j++ {
		err := cmd[j].Start()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
			os.Exit(1)
		}
	}
}

func attachToLog(cmd *exec.Cmd) {
	reader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(reader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("Running Logger | %s\n", scanner.Text())
		}
	}()
}