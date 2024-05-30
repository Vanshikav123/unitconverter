package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/convert", convert)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // Optional for static files
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func convert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	valueStr := r.FormValue("value")
	conversionType := r.FormValue("type")

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		http.Error(w, "Invalid value", http.StatusBadRequest)
		return
	}

	var result float64
	var resultStr string
	switch conversionType {
	case "temperature":
		result = (value * 9 / 5) + 32
		resultStr = strconv.FormatFloat(result, 'f', 2, 64) + " °F"
	case "length":
		result = value * 3.28084
		resultStr = strconv.FormatFloat(result, 'f', 2, 64) + " ft"
	case "weight":
		result = value * 2.20462
		resultStr = strconv.FormatFloat(result, 'f', 2, 64) + " lbs"
	default:
		http.Error(w, "Invalid conversion type", http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/result.html"))
	tmpl.Execute(w, resultStr)
}
