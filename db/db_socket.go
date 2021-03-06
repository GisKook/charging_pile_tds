package db

import (
	"database/sql"
	"fmt"
	"github.com/giskook/charging_pile_tds/conf"
	"github.com/giskook/charging_pile_tds/pb"
	_ "github.com/lib/pq"
	"log"
	"sync"
	"time"
)

type DbSocket struct {
	Db                   *sql.DB
	CommandResult        []*Report.Command
	mutex_command_result sync.Mutex
	BeginTraned          bool
	ticker               *time.Ticker
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
			Db:            db,
			CommandResult: make([]*Report.Command, 0),
			ticker:        time.NewTicker(time.Duration(db_config.TranInterval) * time.Second),
			BeginTraned:   false,
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

func (db_socket *DbSocket) RecvCommandResult(command *Report.Command) {
	db_socket.mutex_command_result.Lock()
	db_socket.CommandResult = append(db_socket.CommandResult, command)
	db_socket.mutex_command_result.Unlock()
}

func (db_socket *DbSocket) DoWork() {
	defer func() {
		db_socket.Close()
	}()

	for {
		select {
		case <-db_socket.ticker.C:
			go db_socket.ProccessCommandResult()

		}
	}
}
