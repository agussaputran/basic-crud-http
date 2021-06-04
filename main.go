package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var db = []Data{}

type Data struct {
	ID        uint32 `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       uint16 `json:"age"`
}

type Response struct {
	Code   uint16      `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func main() {

	// ROUTER
	http.HandleFunc("/get", getApi)
	http.HandleFunc("/post", postApi)
	// http.HandleFunc("/put", putApi)
	// http.HandleFunc("/delete", deleteApi)
	log.Println("Server up and run on localhost:9000...")
	http.ListenAndServe(":9000", nil)
}

func writeResponse(w http.ResponseWriter, r *http.Request, data interface{}, err error) {
	var response Response

	if err != nil {
		response = Response{
			Code:   http.StatusInternalServerError,
			Status: fmt.Sprint(err),
			Data:   data,
		}
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err == nil {
		response = Response{
			Code:   http.StatusOK,
			Status: "Ok!",
			Data:   data,
		}
	}

	respByte, _ := json.Marshal(response)
	w.Write(respByte)

}

func getApi(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, r, db, nil)
}

func postApi(w http.ResponseWriter, r *http.Request) {
	var dataReq Data

	if r.Method == "POST" {

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
		fmt.Fprintln(w, "badrequest")
		writeResponse(w, r, nil, fmt.Errorf("Bad Request"))
	}
}
