package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", HomePage)
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}

func renderTemplate(w http.ResponseWriter, templatePath string, data interface{}) {
	tmpl, err := template.ParseFiles("./pages/" + templatePath)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "HomePage.html", nil)
}
