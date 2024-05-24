package frontend

import (
	"html/template"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// Render the index template
	t, err := template.ParseFiles("./src/frontend/templates/index.html")
	if err != nil {
		log.Printf("failed to parse template: %s", err)
		return
	}

	// Define the data to pass to the template
	type Data struct {
		Name string
	}

	err = t.Execute(w, Data{Name: "Tony"})
	if err != nil {
		log.Printf("failed to execute template: %s", err)
	}

}