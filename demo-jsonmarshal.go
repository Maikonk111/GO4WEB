package main

import (
	"encoding/json"
	"fmt"
)

type employee struct {
	//อ่านเฉพาะพิมพ์ใหญ่
	ID           int
	EmployeeName string
	Tel          string
}

func main() {
	data, _ := json.Marshal(&employee{101, "Mai Sasiprapha", "0981280533"})
	fmt.Println(string(data))
}
