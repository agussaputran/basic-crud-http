package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var db = []data{}

type data struct {
	ID        uint32 `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       uint16 `json:"age"`
}

type response struct {
	Code   uint16      `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func main() {

	// ROUTER
	http.HandleFunc("/get", getApi)
	http.HandleFunc("/post", postApi)
	http.HandleFunc("/put", putApi)
	http.HandleFunc("/delete", deleteApi)
	log.Println("Server up and run on localhost:9000...")
	http.ListenAndServe(":9000", nil)
}

func writeResponse(w http.ResponseWriter, r *http.Request, data interface{}, err error) {
	var res response

	if err != nil {
		res = response{
			Code:   http.StatusInternalServerError,
			Status: fmt.Sprint(err),
			Data:   data,
		}
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err == nil {
		res = response{
			Code:   http.StatusOK,
			Status: "Ok!",
			Data:   data,
		}
	}

	respByte, _ := json.Marshal(res)
	w.Write(respByte)

}

func getApi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		for _, v := range db {
			fmt.Println(v)
		}
		writeResponse(w, r, db, nil)
	} else {
		w.WriteHeader(405)
		fmt.Fprintln(w, "Method Not Allowed")
	}
}

func postApi(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		var dataReq data

		dataByte, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(dataByte, &dataReq)
		if err != nil {
			log.Println("Error on decode data", err)
			writeResponse(w, r, nil, err)
			return
		}
		dataReq.ID = uint32(len(db) + 1)
		db = append(db, dataReq)
		// fmt.Println(data)

		writeResponse(w, r, dataReq, err)
	} else {
		w.WriteHeader(405)
		fmt.Fprintln(w, "Method Not Allowed")
	}
}

func putApi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" || r.Method == "PATCH" {
		var dataReq data
		dataByte, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(dataByte, &dataReq)
		if err != nil {
			log.Println("Error on decode data", err)
			writeResponse(w, r, nil, err)
			return
		}

		q := r.URL.Query()

		for i, v := range db {
			idReq, _ := strconv.Atoi(q.Get("id"))
			if v.ID == uint32(idReq) {
				db[i].FirstName = dataReq.FirstName
				db[i].LastName = dataReq.LastName
				db[i].Age = dataReq.Age
			}
		}

	} else {
		w.WriteHeader(405)
		fmt.Fprintln(w, "Method Not Allowed")
	}
}

func deleteApi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		q := r.URL.Query()
		idReq, _ := strconv.Atoi(q.Get("id"))
		for i, v := range db {
			if v.ID == uint32(idReq) {
				// deleting data with index selected
				db = append(db[:i], db[i+1:]...)
			}
		}
		fmt.Println()
		for _, v := range db {
			fmt.Println(v)
		}
	} else {
		w.WriteHeader(405)
		fmt.Fprintln(w, "Method Not Allowed")
	}
}
