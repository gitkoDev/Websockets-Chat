package main

import (
	"log"
	"net/http"
)


func main()  {
	r := newRoom()

	http.Handle("/", r)

	go r.run()

	log.Println("server running on port :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}