package frontend

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/venatiodecorus/ml-deploy/src/utils"
)

var templates map[string]*template.Template

func Index(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Name string
	}

	err := templates["home.html"].ExecuteTemplate(w, "base", Data{Name: "Tony"})
	if err != nil {
		log.Printf("failed to execute template: %s", err)
	}

}

func Docker(w http.ResponseWriter, r *http.Request) {

	// Define the data to pass to the template
	// type Data struct {
	// 	Name string
	// }

	err := templates["docker.html"].ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Printf("failed to execute template: %s", err)
	}

}

func DockerListHandler(w http.ResponseWriter, r *http.Request) {
	imageList,err := utils.DockerList()
	if err != nil {
		log.Printf("failed to get image list: %s", err)
	}

	err = templates["imageList.html"].ExecuteTemplate(w, "imageList", imageList)
	if err != nil {
		log.Printf("failed to execute template: %s", err)
	}
}

func RegisterRoutes() {
    // if templates == nil {
    //     templates = make(map[string]*template.Template)
    // }

    // layouts, err := filepath.Glob("./src/frontend/templates/*.html")
    // if err != nil {
    //     log.Fatal(err)
    // }

    // for _, layout := range layouts {
	// 	// Apparently need to include any sub templates required by base.html here, for example navbar.html
    //     templates[filepath.Base(layout)] = template.Must(template.ParseFiles(layout, "./src/frontend/templates/base.html", "./src/frontend/templates/navbar.html"))
    // }

	if templates == nil {
        templates = make(map[string]*template.Template)
    }

    err := filepath.Walk("./src/frontend/templates", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        if filepath.Ext(path) != ".html" {
            return nil
        }

        templates[filepath.Base(path)] = template.Must(template.ParseFiles(path, "./src/frontend/templates/base.html", "./src/frontend/templates/navbar.html"))
        return nil
    })

    if err != nil {
        log.Fatal(err)
    }

	http.HandleFunc("/frontend", Index)
	http.HandleFunc("/docker", Docker)

	http.HandleFunc("/docker/list", DockerListHandler)
}