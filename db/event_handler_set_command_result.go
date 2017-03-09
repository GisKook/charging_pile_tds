package db

import (
	"log"
)

const sql_update_command_result = "UPDATE t_command_send_log SET response_result=$1 ,response_time=now() , response_content=$2 where cpid=$3 and serial_number=$4"

func (db_socket *DbSocket) ProccessCommandResult() {
	log.Println("proccess set price result")
	db_socket.mutex_command_result.Lock()
	if len(db_socket.CommandResult) == 0 {
		db_socket.mutex_command_result.Unlock()
		return
	}
	db_socket.mutex_command_result.Unlock()

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

	db_socket.mutex_command_result.Lock()
	var content string = ""
	for _, set_command_result := range db_socket.CommandResult {
		cpid := set_command_result.Tid
		serial := set_command_result.SerialNumber
		result := uint8(set_command_result.Paras[0].Npara)
		if len(set_command_result.Paras) > 1 {
			content = set_command_result.Paras[1].Strpara
		}
		log.Println(set_command_result.Paras)
		log.Println(result)
		log.Println(content)
		log.Println(cpid)
		log.Println(serial)
		_, err = stmt.Exec(result, content, cpid, serial)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	} else {
		for i, _ := range db_socket.CommandResult {
			db_socket.CommandResult[i] = nil
		}
		db_socket.CommandResult = db_socket.CommandResult[:0]
	}
	db_socket.mutex_command_result.Unlock()
}
