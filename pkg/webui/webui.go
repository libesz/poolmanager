package webui

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/libesz/poolmanager/pkg/scheduler"
	"github.com/libesz/poolmanager/pkg/webui/content/static"
	"github.com/libesz/poolmanager/pkg/webui/content/templates"
	"github.com/shurcooL/httpfs/html/vfstemplate"
)

var parsedTemplates *template.Template
var store = sessions.NewCookieStore([]byte("temp"))

func Run(configStore *scheduler.ConfigStore) {
	r := mux.NewRouter()
	parsedTemplates = template.Must(vfstemplate.ParseGlob(templates.Content, nil, "*.html"))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(static.Content)))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(configStore, w, r)
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homePostHandler(configStore, w, r)
	}).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

type ConfigItemWithType struct {
	Key          string
	DetectedType string
	Value        interface{}
}

type PageData struct {
	AllConfig map[string][]ConfigItemWithType
	Function  string
	Debug     string
}

func homeHandler(configStore *scheduler.ConfigStore, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	allConfig := make(map[string][]ConfigItemWithType)
	controllers := configStore.GetKeys()
	for _, controller := range controllers {
		configForController := configStore.Get(controller)
		for key, value := range configForController {
			item := ConfigItemWithType{Key: key, Value: value}
			switch value.(type) {
			case int:
				item.DetectedType = "int"
			case float64:
				item.DetectedType = "float64"
			case bool:
				item.DetectedType = "bool"
			case string:
				item.DetectedType = "string"
			}
			allConfig[controller] = append(allConfig[controller], item)
		}
	}
	data := PageData{AllConfig: allConfig, Function: "default"}
	log.Printf("%+v\n", data)
	if err := parsedTemplates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Println(err.Error())
	}
}

func homePostHandler(configStore *scheduler.ConfigStore, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	r.ParseForm()
	data := PageData{Function: "debug", Debug: r.PostForm.Encode()}
	log.Printf("%+v\n", data)
	if err := parsedTemplates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Println(err.Error())
	}
}
