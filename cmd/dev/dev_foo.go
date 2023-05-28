// Copyright (c) 2023, Geert JM Vanderkelen

package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"time"

	_ "github.com/golistic/pxmysql" // load SQL driver
)

func main() {
	bug28()
}

func bug28() {
	db, err := sql.Open("mysqlpx", "kelvin:green@tcp(127.0.0.1:33060)/?useTLS=true")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	_, err = db.Exec("SET @@SESSION.wait_timeout = 2")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	time.Sleep(3 * time.Second)

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if err := tx.Rollback(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}

func bugMaxConnections() {
	chanSignals := make(chan os.Signal, 1)
	signal.Notify(chanSignals, os.Interrupt)

	db, err := sql.Open("mysqlpx", "kelvin:green@tcp(127.0.0.1:33060)/?useTLS=true")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	start := time.Now()

	for {
		select {
		case <-chanSignals:
			fmt.Println("\rdone    ")
			return
		default:
			fmt.Print(time.Now().Format("15:04:05"))
			time.Sleep(time.Second * 1)
			if time.Now().After(start.Add(1 * time.Second)) {
				fmt.Println("\ncheck    ")
				_, err := db.Query("SELECT 1")
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				start = time.Now()
			} else {
				fmt.Print("\r")
			}
		}
	}
}
