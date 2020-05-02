package webui

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/libesz/poolmanager/pkg/configstore"
	"github.com/libesz/poolmanager/pkg/controller"
	"github.com/libesz/poolmanager/pkg/webui/content/static"
	"github.com/libesz/poolmanager/pkg/webui/content/templates"
	"github.com/shurcooL/httpfs/html/vfstemplate"
)

var parsedTemplates *template.Template
var store = sessions.NewCookieStore([]byte("temp"))

func Run(configStore *configstore.ConfigStore) {
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
	Slider       controller.ConfigRange
	Toggle       controller.ConfigToggle
	Value        interface{}
}

type PageData struct {
	AllConfig map[string][]ConfigItemWithType
	Function  string
	Debug     string
}

func homeHandler(configStore *configstore.ConfigStore, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	allConfig := make(map[string][]ConfigItemWithType)
	controllers := configStore.GetKeys()
	for _, controllerName := range controllers {
		configPropertiesForController := configStore.GetProperties(controllerName)
		configValuesForController := configStore.Get(controllerName)
		for key, property := range configPropertiesForController {
			item := ConfigItemWithType{Key: key, Value: configValuesForController[key]}
			switch property.(type) {
			case controller.ConfigRange:
				item.DetectedType = "slider"
				item.Slider = property.(controller.ConfigRange)
			case controller.ConfigToggle:
				item.DetectedType = "toggle"
				item.Toggle = property.(controller.ConfigToggle)
			}
			allConfig[controllerName] = append(allConfig[controllerName], item)
		}
	}
	data := PageData{AllConfig: allConfig, Function: "default"}
	log.Printf("%+v\n", data)
	if err := parsedTemplates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Println(err.Error())
	}
}

type JsonRequest struct {
	Controller string `json:"controller"`
	Key        string `json:"key"`
	Value      string `json:"value"`
}

type JsonResponse struct {
	Error string `json:"error"`
}

func homePostHandler(configStore *configstore.ConfigStore, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data JsonRequest
	err := decoder.Decode(&data)
	if err != nil {
		log.Printf("Decode error on request: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("%+v\n", data)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	u := JsonResponse{Error: "bla"}
	json.NewEncoder(w).Encode(u)
}
