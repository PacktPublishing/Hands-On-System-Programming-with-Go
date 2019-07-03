package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	tpl, err := template.ParseGlob("ch09/template/parse/*.html")
	if err != nil {
		log.Fatal("Error:", err)
	}
	data := map[string]string{
		"name":       "Jin Kazama",
		"style":      "Karate",
		"appearance": "Tekken 3",
	}
	if err := tpl.Execute(os.Stdout, data); err != nil {
		log.Fatal("Error:", err)
	}
}
