package db

import (
	"database/sql"
	"fmt"
	"github.com/giskook/charging_pile_tds/conf"
	"github.com/giskook/charging_pile_tds/pb"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type DbSocket struct {
	Db             *sql.DB
	SetPriceResult []*Report.Command
	ticker         *time.Ticker
}

var G_DbSocket *DbSocket = nil

func NewDbSocket(db_config *conf.DBConfigure) (*DbSocket, error) {
	if G_DbSocket == nil {
		conn_string := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", db_config.User, db_config.Passwd, db_config.Host, db_config.Port, db_config.DbName)

		log.Println(conn_string)
		db, err := sql.Open(db_config.User, conn_string)

		if err != nil {
			return nil, err
		}
		log.Println("db open success")

		G_DbSocket = &DbSocket{
			Db:             db,
			SetPriceResult: make([]*Report.Command, 0),
			ticker:         time.NewTicker(time.Duration(db_config.TranInterval) * time.Second),
		}

	}

	return G_DbSocket, nil
}

func GetDBScoket() *DbSocket {
	return G_DbSocket
}

func (db_socket *DbSocket) Close() {
	db_socket.ticker.Stop()
	db_socket.Db.Close()
}

func (db_socket *DbSocket) RecvSetPriceResult(command *Report.Command) {
	db_socket.SetPriceResult = append(db_socket.SetPriceResult, command)
}

func (db_socket *DbSocket) DoWork() {
	defer func() {
		db_socket.Close()
	}()

	for {
		select {
		case <-db_socket.ticker.C:
			go db_socket.ProccessSetPriceResult()
		}
	}
}
