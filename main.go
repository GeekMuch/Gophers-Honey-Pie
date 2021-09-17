package main

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

/*
	Runs fucntions in order.
*/
func main() {
	log.InitLog(true)
	// checkForInternet()
	config.StartSetupSequence("192.168.8.103")

	// fmt.Println("\n [+] Server running!")
	// http.HandleFunc("/", handler)
	// log.Fatal(http.ListenAndServe(":8080", nil))

}
