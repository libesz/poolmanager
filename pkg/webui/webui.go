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

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/libesz/poolmanager/pkg/configstore"
	"github.com/libesz/poolmanager/pkg/controller"
	"github.com/libesz/poolmanager/pkg/io"
	"github.com/libesz/poolmanager/pkg/webui/content/static"
	"github.com/libesz/poolmanager/pkg/webui/content/templates"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/urfave/negroni"
)

var parsedTemplates *template.Template
var mySigningKey = []byte("captainjacksparrowsayshi")

func New(listenOn, password string, configStore *configstore.ConfigStore, inputs []io.Input, outputs []io.Output) WebUI {
	r := mux.NewRouter()
	parsedTemplates = template.Must(vfstemplate.ParseGlob(templates.Content, nil, "*.html"))

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(static.Content)))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(configStore, inputs, outputs, w, r)
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homePostHandler(configStore, w, r)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginPostHandler(password, w, r)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginGetHandler(w, r)
	}).Methods("GET")

	r.Handle("/api/ping", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(myHandler),
	))

	r.Handle("/api/status", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiStatusHandler(configStore, inputs, outputs, w, r)
		})),
	))

	server := &http.Server{
		Handler:      r,
		Addr:         listenOn,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return WebUI{server: server, jwt: jwtMiddleware}
}

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	fmt.Fprintf(w, "This is an authenticated request")
	fmt.Fprintf(w, "Claim content:\n")
	for k, v := range user.(*jwt.Token).Claims.(jwt.MapClaims) {
		fmt.Fprintf(w, "%s :\t%#v\n", k, v)
	}
})

type ApiStatusResponse struct {
	Inputs          map[string]float64 `json:"inputs"`
	InputErrorConst float64            `json:"inputerrorconst"`
	Outputs         map[string]bool    `json:"outputs"`
}

func apiStatusHandler(configStore *configstore.ConfigStore, inputs []io.Input, outputs []io.Output, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := ApiStatusResponse{
		Inputs:          make(map[string]float64),
		InputErrorConst: io.InputError,
		Outputs:         make(map[string]bool),
	}
	for _, item := range inputs {
		data.Inputs[item.Name()] = item.Value()
	}

	for _, item := range outputs {
		data.Outputs[item.Name()] = item.Get()
	}

	log.Printf("Webui: rendering apiStatusHandler page with data: %+v\n", data)
	_ = json.NewEncoder(w).Encode(data)
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
	InputErrorConst  float64
	Outputs          []io.Output
	Function         string
	Debug            string
}

func homeHandler(configStore *configstore.ConfigStore, inputs []io.Input, outputs []io.Output, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	function := "default"

	data := PageData{
		ConfigProperties: make(map[string]controller.ConfigProperties),
		ConfigValues:     make(map[string]controller.Config),
		Function:         function,
	}

	controllers := configStore.GetKeys()
	for _, controllerName := range controllers {
		data.ConfigProperties[controllerName] = configStore.GetProperties(controllerName)
		data.ConfigValues[controllerName] = configStore.Get(controllerName)
	}
	data.Inputs = inputs
	data.InputErrorConst = io.InputError
	data.Outputs = outputs

	log.Printf("Webui: rendering page with data: %+v\n", data)
	if err := parsedTemplates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Println(err.Error())
	}
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	function := "login"

	data := PageData{
		Function: function,
	}

	log.Printf("Webui: rendering login page with data: %+v\n", data)
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
	Token     string      `json:"token"`
	OrigValue interface{} `json:"origValue"`
}

func homePostHandler(configStore *configstore.ConfigStore, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var data JsonRequest
	err := decoder.Decode(&data)
	if err != nil {
		log.Printf("Webui: decode error on request: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Webui: requested config change: %+v\n", data)

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

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func loginPostHandler(password string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Printf("Webui: requested login\n")
	decoder := json.NewDecoder(r.Body)
	var data LoginRequest
	err := decoder.Decode(&data)
	if err != nil {
		log.Printf("Webui: decode error on login request: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		u := JsonResponse{Error: err.Error(), OrigValue: nil}
		_ = json.NewEncoder(w).Encode(u)
		return
	}

	if data.Password != password {
		log.Printf("Webui: user unauthorized\n")
		w.WriteHeader(http.StatusUnauthorized)
		u := JsonResponse{Error: "Unauthorized", OrigValue: nil}
		_ = json.NewEncoder(w).Encode(u)
		return
	}
	log.Printf("Webui: user authorized\n")
	tokenString, _ := generateJWT()
	w.WriteHeader(http.StatusAccepted)
	u := JsonResponse{Token: tokenString}
	_ = json.NewEncoder(w).Encode(u)
}
