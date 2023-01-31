package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/common-nighthawk/go-figure"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	log.Println("Starting")

	date := time.Now()

	log.Println("Creating a client")

	c, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Authenticating")

	err = c.Auth()
	if err != nil {
		log.Fatal(err)
	}

	refresh := func() {
		log.Printf("Loading data for %s", date.String())
		sum := GetOnlineSum(c, date.Format("2006-01-02"))

		f := figure.NewColorFigure(fmt.Sprintf("Online: %.2f", sum), "", "green", true)
		f.Print()
	}

	refresh()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		newDate, err := time.Parse("2006-01-02", scanner.Text())
		if err != nil {
			log.Println(err)
		} else {
			date = newDate
			refresh()
		}
	}
}
