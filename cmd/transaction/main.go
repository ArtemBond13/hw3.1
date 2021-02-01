package main

import (
	"encoding/xml"
	transaction2 "github.com/ArtemBond13/hw3.1/pkg/transaction"
	"log"
	"time"
)

func main() {
	transaction := []transaction2.Transaction{
		{
			Id: "2123",
			From: "3456",
			To: "0076",
			Amount: 1000_00,
			Created: time.Now().Unix(),
		},
		{
			Id:      "2",
			From:    "0001",
			To:      "0002",
			Amount:  200_00,
			Created: time.Now().Unix(),
		},
	}
	// json.Marshal() возвращает срез байт []byte
	encoded, err := xml.Marshal(transaction)
	if err != nil {
		log.Println(err)
	}
	// добавил общий заголовок XML
	encoded = append([]byte(xml.Header), encoded...)
	log.Println(string(encoded))

	var decoded []transaction2.Transaction
	// Важно: передаём указатель, чтобы функция смогла записать данные
	err = xml.Unmarshal(encoded, &decoded)
	log.Printf("%v\n", decoded)


}
