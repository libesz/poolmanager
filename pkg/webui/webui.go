package webui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/libesz/poolmanager/pkg/configstore"
	"github.com/libesz/poolmanager/pkg/controller"
	"github.com/libesz/poolmanager/pkg/io"
	"github.com/libesz/poolmanager/pkg/webui/content/static"
	"github.com/libesz/poolmanager/pkg/webui/content/templates"
	"github.com/shurcooL/httpfs/html/vfstemplate"
)

var parsedTemplates *template.Template

func New(listenOn, password string, configStore *configstore.ConfigStore, inputs []io.Input, outputs []io.Output) WebUI {
	s := sessions.NewCookieStore([]byte("temp"))
	r := mux.NewRouter()
	parsedTemplates = template.Must(vfstemplate.ParseGlob(templates.Content, nil, "*.html"))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(static.Content)))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(s, configStore, inputs, outputs, w, r)
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homePostHandler(configStore, w, r)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginPostHandler(password, s, w, r)
	}).Methods("POST")

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logoutGetHandler(s, w, r)
	}).Methods("GET")

	server := &http.Server{
		Handler:      r,
		Addr:         listenOn,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return WebUI{server: server, sessions: s}
}

func (w *WebUI) Run(stopChan chan struct{}) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := w.server.ListenAndServe()
		log.Printf("Webui: %s\n", err.Error())
		wg.Done()
	}()
	<-stopChan
	w.server.Close()
	wg.Wait()
}

type PageData struct {
	ConfigProperties map[string]controller.ConfigProperties
	ConfigValues     map[string]controller.Config
	Inputs           []io.Input
	Outputs          []io.Output
	Function         string
	Debug            string
}

func homeHandler(s *sessions.CookieStore, configStore *configstore.ConfigStore, inputs []io.Input, outputs []io.Output, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	session, _ := s.Get(r, "session")
	loggedIn, ok := session.Values["logged-in"].(bool)
	if !ok {
		loggedIn = false
	}

	function := "default"
	if !loggedIn {
		function = "login"
	}

	data := PageData{
		ConfigProperties: make(map[string]controller.ConfigProperties),
		ConfigValues:     make(map[string]controller.Config),
		Function:         function,
	}

	if loggedIn {
		controllers := configStore.GetKeys()
		for _, controllerName := range controllers {
			data.ConfigProperties[controllerName] = configStore.GetProperties(controllerName)
			data.ConfigValues[controllerName] = configStore.Get(controllerName)
		}
		data.Inputs = inputs
		data.Outputs = outputs
	}

	log.Printf("Webui: rendering page with data: %+v\n", data)
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
		log.Printf("Webui: decode error on request: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Webui: requested config change: %+v\n", data)

	w.Header().Set("Content-Type", "application/json")

	config := configStore.Get(data.Controller)
	origValueString := ""
	switch data.Type {
	case "range":
		log.Printf("Webui: identified numeric config\n")
		convertedValue, err := strconv.ParseFloat(data.Value, 64)
		if err != nil {
			log.Printf("Webui: Failed to parse requested numeric config change for controller %s key %s value %s: %s\n", data.Controller, data.Key, data.Value, err.Error())
			respondError(w, origValueString, err)
			return
		}
		origValue, ok := config.Ranges[data.Key]
		if !ok {
			log.Printf("Webui: Non-existing numeric config key for controller %s key %s value %s\n", data.Controller, data.Key, data.Value)
			respondError(w, origValueString, fmt.Errorf("Non-existing config key"))
			return
		}
		origValueString = strconv.FormatFloat(origValue, 'E', -1, 64)
		config.Ranges[data.Key] = convertedValue
	case "toggle":
		log.Printf("Webui: identified boolean config\n")
		convertedValue, err := strconv.ParseBool(data.Value)
		if err != nil {
			log.Printf("Webui: Failed to parse requested boolean config change for controller %s key %s value %s: %s\n", data.Controller, data.Key, data.Value, err.Error())
			respondError(w, origValueString, err)
			return
		}
		origValue, ok := config.Toggles[data.Key]
		if !ok {
			log.Printf("Webui: Non-existing boolean config key for controller %s key %s value %s\n", data.Controller, data.Key, data.Value)
			respondError(w, origValueString, fmt.Errorf("Non-existing config key"))
			return
		}
		origValueString = strconv.FormatBool(origValue)
		config.Toggles[data.Key] = convertedValue
	default:
		log.Printf("Webui: unknown type: %s\n", data.Type)
	}
	err = configStore.Set(data.Controller, config, true)
	if err != nil {
		log.Printf("Webui: Failed to update config for controller %s: %s\n", data.Controller, err.Error())
		respondError(w, origValueString, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	u := JsonResponse{Error: "OK"}
	_ = json.NewEncoder(w).Encode(u)
}

func respondError(w http.ResponseWriter, origValue interface{}, err error) {
	w.WriteHeader(http.StatusConflict)
	u := JsonResponse{Error: err.Error(), OrigValue: origValue}
	_ = json.NewEncoder(w).Encode(u)
}

type LoginRequest struct {
	Password string `json:"password"`
}

func loginPostHandler(password string, s *sessions.CookieStore, w http.ResponseWriter, r *http.Request) {
	log.Printf("Webui: requested login\n")
	decoder := json.NewDecoder(r.Body)
	var data LoginRequest
	err := decoder.Decode(&data)
	if err != nil {
		log.Printf("Webui: decode error on request: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data.Password != password {
		log.Printf("Webui: user unauthorized\n")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Printf("Webui: user authorized\n")
	session, _ := s.Get(r, "session")
	session.Values["logged-in"] = true
	_ = session.Save(r, w)
	w.WriteHeader(http.StatusAccepted)
}

func logoutGetHandler(s *sessions.CookieStore, w http.ResponseWriter, r *http.Request) {
	log.Printf("Webui: requested logout\n")
	session, _ := s.Get(r, "session")
	if loggedIn, ok := session.Values["logged-in"].(bool); !ok || !loggedIn {
		log.Printf("Webui: not logged in while requesting logout\n")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	session.Values["logged-in"] = false
	_ = session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
