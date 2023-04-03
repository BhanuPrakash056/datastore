package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Error struct {
	Res string `json:"error"`
}
type Value struct {
	Value string `json:"value"`
}
type User struct {
	Command string `json:"command"`
}
type Store struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Expiry int    `json:"expiry"`
}
type Queue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var stores []Store
var queue []Queue

func commHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	con := user.Command
	words := strings.Fields(con)
	if words[0] == "SET" {
		setMethod(words, w, r)
	} else if words[0] == "GET" {
		getMethod(words, w, r)
	} else if words[0] == "QPUSH" {
		push(words, w, r)
	} else if words[0] == "QPOP" {
		pop(words, w, r)
	} else {
		res := Error{Res: "invalid command"}
		json.NewEncoder(w).Encode(res)
		return
	}

}

func main() {
	// r := mux.NewRouter()
	// r.HandleFunc("/comm", getCommand).Methods("POST")
	http.HandleFunc("/comm", commHandler)
	log.Fatal(http.ListenAndServe(":8100", nil))
}

func setMethod(words []string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if len(words) == 3 {
		var temp Store
		temp.Key = words[1]
		temp.Value = words[2]
		stores = append(stores, temp)

	} else if len(words) == 4 {
		if words[3] == "XX" {
			for i := 0; i < len(stores); i++ {
				if stores[i].Key == words[1] {
					stores[i].Value = words[2]
					return
				}
			}
			w.WriteHeader(http.StatusBadRequest)
			res := Error{Res: "Record does not exist"}
			json.NewEncoder(w).Encode(res)
			return
		}
		res := Error{Res: "invalid command"}
		json.NewEncoder(w).Encode(res)
		return
	} else if len(words) == 5 {
		if words[3] != "EX" {
			res := Error{Res: "invalid command"}
			json.NewEncoder(w).Encode(res)
			return
		}
		var temp Store
		temp.Key = words[1]
		var str string = words[4]
		myInt, err := strconv.Atoi(str)
		if err != nil {
			res := Error{Res: "invalid command"}
			json.NewEncoder(w).Encode(res)
			return
		}
		temp.Expiry = myInt
		stores = append(stores, temp)
	}
	if len(words) == 6 {
		if words[3] != "EX" {
			res := Error{Res: "invalid command"}
			json.NewEncoder(w).Encode(res)
			return
		}
		if words[5] == "NX" {
			for i := 0; i < len(stores); i++ {
				if stores[i].Key == words[1] {
					res := Error{Res: "record already exists"}
					json.NewEncoder(w).Encode(res)
					return
				}
			}
			var temp Store
			temp.Key = words[1]
			temp.Value = words[2]
			var str string = words[4]
			myInt, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("This is a err", err)
			}
			temp.Expiry = myInt
			stores = append(stores, temp)
			return
		}
		res := Error{Res: "invalid command"}
		json.NewEncoder(w).Encode(res)
		return
	}
}
func getMethod(words []string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if len(words) > 2 {
		res := Error{Res: "invalid command"}
		json.NewEncoder(w).Encode(res)
		return
	}
	for i := 0; i < len(stores); i++ {
		if stores[i].Key == words[1] {
			res := Value{Value: stores[i].Value}
			json.NewEncoder(w).Encode(res)
			return
		}
		res := Value{Value: ""}
		json.NewEncoder(w).Encode(res)
		return
	}

}

func push(words []string, w http.ResponseWriter, r *http.Request) {
	var temp Queue
	temp.Key = words[1]
	temp.Value = words[2]
	queue = append(queue, temp)
}
func pop(words []string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if len(words) < 2 {
		res := Error{Res: "invalid command"}
		json.NewEncoder(w).Encode(res)
		return
	}
	if len(queue) == 0 {
		res := Error{Res: "Queue is empty"}
		json.NewEncoder(w).Encode(res)
		return
	}
	element := queue[0]
	queue = queue[1:]
	res := Value{Value: element.Value}
	json.NewEncoder(w).Encode(res)
	return

}
