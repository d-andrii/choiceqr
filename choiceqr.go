package main

import (
	"log"
	"os"

	"gioui.org/app"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	log.Println("Starting")

	go func() {
		log.Println("Opening window")
		w := app.NewWindow(app.Size(300, 300), app.Title("Choice Online"))
		a, err := NewApp(w)
		if err != nil {
			log.Fatal(err)
		}
		if err = a.Loop(); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	app.Main()
}
