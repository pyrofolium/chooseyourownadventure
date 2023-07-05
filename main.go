package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func (s Arc) createHandler(templateFolder string, templateName string) http.HandlerFunc {
	//templateObject, err := template.New("arc").ParseFiles("templates/arc.gohtml", "templates/content.gohtml", "templates/footer.gohtml", "templates/header.gohtml", "templates/options.gohtml")
	templateObject, err := template.New("arc").ParseFiles("templates/test.gohtml")

	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		err := templateObject.Execute(os.Stdout, s)

		err = templateObject.Execute(writer, s)
		if err != nil {
			http.Error(writer, "error creating template", 500)
		}

	}

}

func createRoutes(arcs map[string]Arc) *http.ServeMux {
	mux := http.NewServeMux()
	for arcName, arc := range arcs {
		mux.HandleFunc(fmt.Sprintf("/%s", arcName),
			arc.createHandler("templates", "arc"))
	}
	return mux
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"Arc"`
}

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

func main() {
	file, err := os.Open("gopher.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	jsonBytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	arcs := make(map[string]Arc)
	err = json.Unmarshal(jsonBytes, &arcs)

	mux := createRoutes(arcs)
	http.ListenAndServe(":8080", mux)

}
