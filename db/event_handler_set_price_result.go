package db

import (
	"log"
)

const sql_update_command_result = "UPDATE t_command_send_log SET response_result=$1 ,response_time=now()  where cpid=$2 and serila_number=$3"

func (db_socket *DbSocket) ProccessSetPriceResult() {
	log.Println("proccess set price result")

	tx, err := db_socket.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	stmt, e := tx.Prepare(sql_update_command_result)
	if e != nil {
		log.Fatal(e)
	}

	defer stmt.Close()
	for _, set_price_result := range db_socket.SetPriceResult {
		cpid := set_price_result.Tid
		serial := set_price_result.SerialNumber
		result := uint8(set_price_result.Paras[0].Npara)
		_, err = stmt.Exec(result, cpid, serial)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	} else {
		for i, _ := range db_socket.SetPriceResult {
			db_socket.SetPriceResult[i] = nil
		}
		db_socket.SetPriceResult = db_socket.SetPriceResult[:0]
	}
}
