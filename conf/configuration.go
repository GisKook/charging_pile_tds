package conf

import (
	"encoding/json"
	"os"
)

type ConsumerConf struct {
	Addr     string
	Topic    string
	Channels []string
}

type NsqConfiguration struct {
	Consumer *ConsumerConf
}

type DBConfigure struct {
	Host         string
	Port         string
	User         string
	Passwd       string
	DbName       string
	TranInterval int
}

type Configuration struct {
	Uuid string
	Nsq  *NsqConfiguration
	DB   *DBConfigure
}

var G_conf *Configuration

func ReadConfig(confpath string) (*Configuration, error) {
	file, _ := os.Open(confpath)
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	G_conf = &config

	return &config, err
}

func GetConf() *Configuration {
	return G_conf
}
