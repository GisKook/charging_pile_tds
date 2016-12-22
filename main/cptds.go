package main

import (
	"fmt"
	"github.com/giskook/charging_pile_tds/conf"
	"github.com/giskook/charging_pile_tds/db"
	"github.com/giskook/charging_pile_tds/mq"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// read configuration
	configuration, err := conf.ReadConfig("./conf.json")

	checkError(err)
	// create a mq socket
	mq_socket := mq.NewNsqSocket(configuration.Nsq)
	go mq_socket.Start()
	// create a db socket
	db_socket, e := db.NewDbSocket(configuration.DB)
	checkError(e)
	db_socket.DoWork()

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
	db_socket.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
