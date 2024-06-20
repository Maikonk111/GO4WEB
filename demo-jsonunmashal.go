package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type employees struct {
	//อ่านเฉพาะพิมพ์ใหญ่
	ID           int
	EmployeeName string
	Tel          string
}

func main() {
	e := employees{}
	err := json.Unmarshal([]byte(`{"ID":101,"Employee":"Sasiprapha","Tel":"098121112"}`), &e)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(e.EmployeeName)
}
