package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
			"ID":101,
			"name":"Python",
			"price":2590,
			"instructor":"BorntoDev"
		},
		{
			"ID":102,
			"name":"JavaScript",
			"price":0,
			"instructor":"BorntoDev"
		},
		{
			"ID":103,
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
	}
	return highestID + 1
}

func findID(ID int) (*Course, int) {
	for i, course := range CourseList {
		if course.ID == ID {
			return &course, i
		}
	}
	return nil, -1
}

// hand
func courseHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "course/")
	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	course, listItemIndex := findID(ID)
	if course == nil {
		http.Error(w, fmt.Sprintf("no course with id %d", ID), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		courseJSON, err := json.Marshal(course)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-type", "application/json")
		w.Write(courseJSON)

	case http.MethodPut:
		var updateCourse Course
		byteBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(byteBody, &updateCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updateCourse.ID != ID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		course = &updateCourse
		CourseList[listItemIndex] = *course
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func coursseHandler(w http.ResponseWriter, r *http.Request) {
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
			return
		}

		newCourse.ID = getNextID()
		CourseList = append(CourseList, newCourse)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func main() {
	http.HandleFunc("/course/", courseHandler)
	http.HandleFunc("/course", coursseHandler)
	http.ListenAndServe(":5000", nil)
}
