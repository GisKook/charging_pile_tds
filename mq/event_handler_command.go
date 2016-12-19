package mq

import (
	"github.com/giskook/charging_pile_tds/db"
	"github.com/giskook/charging_pile_tds/pb"
	"github.com/golang/protobuf/proto"
	"log"
)

func RecvNsq(socket *NsqSocket, message []byte) {
	command := &Report.Command{}
	err := proto.Unmarshal(message, command)
	if err != nil {
		log.Println("unmarshal error")
	} else {
		log.Printf("<IN NSQ> %s %d \n", command.Uuid, command.Tid)
		switch command.Type {
		case Report.Command_CMT_REP_NOTIFY_SET_PRICE:
			db.GetDBScoket().RecvSetPriceResult(command)
			break
		}
	}
}
