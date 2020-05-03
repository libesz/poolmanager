package webui

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

type PageData struct {
	ConfigProperties map[string]controller.ConfigProperties
	ConfigValues     map[string]controller.Config
	Function         string
	Debug            string
}

func homeHandler(configStore *configstore.ConfigStore, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	data := PageData{
		ConfigProperties: make(map[string]controller.ConfigProperties),
		ConfigValues:     make(map[string]controller.Config),
		Function:         "default"}
	controllers := configStore.GetKeys()
	for _, controllerName := range controllers {
		data.ConfigProperties[controllerName] = configStore.GetProperties(controllerName)
		data.ConfigValues[controllerName] = configStore.Get(controllerName)
	}
	log.Printf("%+v\n", data)
	if err := parsedTemplates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Println(err.Error())
	}
}

type JsonRequest struct {
	Controller string `json:"controller"`
	Type       string `json:"type"`
	Key        string `json:"key"`
	Value      string `json:"value"`
}

type JsonResponse struct {
	Error     string      `json:"error"`
	OrigValue interface{} `json:"origValue"`
}

func homePostHandler(configStore *configstore.ConfigStore, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data JsonRequest
	err := decoder.Decode(&data)
	if err != nil {
		log.Printf("Webui: ecode error on request: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Webui: requested config change: %+v\n", data)

	w.Header().Set("Content-Type", "application/json")

	config := configStore.Get(data.Controller)
	origValue := ""
	switch data.Type {
	case "range":
		log.Printf("Webui: identified numeric config\n")
		convertedValue, err := strconv.ParseFloat(data.Value, 64)
		if err != nil {
			log.Printf("Webui: Failed to parse requested numeric config change for controller %s key %s value %s: %s\n", data.Controller, data.Key, data.Value, err.Error())
			respondError(w, origValue, err)
			return
		}
		origValue = strconv.FormatFloat(config.Ranges[data.Key], 'E', -1, 64)
		config.Ranges[data.Key] = convertedValue
	case "toggle":
		log.Printf("Webui: identified boolean config\n")
		convertedValue, err := strconv.ParseBool(data.Value)
		if err != nil {
			log.Printf("Webui: Failed to parse requested boolean config change for controller %s key %s value %s: %s\n", data.Controller, data.Key, data.Value, err.Error())
			respondError(w, origValue, err)
			return
		}
		origValue = strconv.FormatBool(config.Toggles[data.Key])
		config.Toggles[data.Key] = convertedValue
	default:
		log.Printf("Webui: unknown type: %s\n", data.Type)
	}
	err = configStore.Set(data.Controller, config)
	if err != nil {
		log.Printf("Webui: Failed to update config for controller %s: %s\n", data.Controller, err.Error())
		respondError(w, origValue, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	u := JsonResponse{Error: "OK"}
	json.NewEncoder(w).Encode(u)
}

func respondError(w http.ResponseWriter, origValue interface{}, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	u := JsonResponse{Error: err.Error(), OrigValue: origValue}
	json.NewEncoder(w).Encode(u)
}
