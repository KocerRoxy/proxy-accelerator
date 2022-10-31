package main

import "github.com/gen2brain/beeep"

func notify(title, body string) {
	if len(title) == 0 {
		title = "Accelerator"
	}
	err := beeep.Notify(title, body, "")
	if err != nil {
		panic(err)
	}
}

func alert(title, body string) {
	if len(title) == 0 {
		title = "Accelerator"
	}
	err := beeep.Alert(title, body, "")
	if err != nil {
		panic(err)
	}
}
