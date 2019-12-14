package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"giftano-crud-golang/routers"
)

func main() {
	r := routers.InitRouter()

	fmt.Println("Server listen at: " + os.Getenv("APP_PORT"))

	log.Fatal(http.ListenAndServe(":" + os.Getenv("APP_PORT"), r ))
}