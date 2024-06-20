package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Course struct {
	ID         int     `json:"ID"`
	Name       string  `json:"name"`
	Price      float32 `json:"price"`
	Instructor string  `json:"instructor"`
}

var CourseList []Course //ตัวแปรเก็บรายการคอร์ส

func init() {
	CourseJSON := `[
		{
			"ID":1,
			"name":"Python",
			"price":2590,
			"instructor":"BorntoDev"
		},
		{
			"ID":2,
			"name":"JavaScript",
			"price":0,
			"instructor":"BorntoDev"
		},
		{
			"ID":3,
			"name":"SQL",
			"price":0,
			"instructor":"BorntoDev"
		}
	]`
	err := json.Unmarshal([]byte(CourseJSON), &CourseList)

	if err != nil {
		log.Fatal(err)
	}
}

// ตัวรัน ID เอง
func getNextID() int {
	highestID := -1
	for _, course := range CourseList {
		if highestID < course.ID {
			highestID = course.ID
		}
		return highestID + 1
	}
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	courseJSON, err := json.Marshal(CourseList)
	switch r.Method {
	case http.MethodGet: //ถ้าเป็น GET: จะส่งข้อมูล CourseList ในรูปแบบ JSON กลับไปยัง client
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) //เช็คข้อผิดพลาดทาง sever
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)
	case http.MethodPost: //ถ้าเป็น POST: จะอ่านข้อมูลคอร์สใหม่จาก body ของ request และทำการเพิ่มข้อมูลนั้นเข้าไปใน CourseList
		var newCourse Course
		Bodybyte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(Bodybyte, &newCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//ตัวรัน ID เอง
		if newCourse.ID != 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		newCourse.ID = getNextID()
		CourseList = append(CourseList, newCourse)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func main() {
	http.HandleFunc("/course", courseHandler)
	http.ListenAndServe(":5000", nil)
}
