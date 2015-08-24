package config

import "gopkg.in/yaml.v2"
import "io/ioutil"

type RabbitConfig struct {
	Host string
	Port string
	Username string
	Password string
	Exchange string
}

type Worker struct {
	Command string
	Arguments []string
	Max int
	Min int
	Logpath string
}

type Config struct {
	Workers []Worker
	Rabbit RabbitConfig
}

func Parse(path string, c chan<- Config) {
	var config Config
	source, err := ioutil.ReadFile(path)
	eCheck(err)

	err = yaml.Unmarshal(source, &config)
	eCheck(err)

	c <- config
}

func eCheck(e error) {
	if e != nil {
		panic(e)
	}
}
