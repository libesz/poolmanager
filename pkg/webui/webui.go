package webui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/libesz/poolmanager/pkg/configstore"
	"github.com/libesz/poolmanager/pkg/controller"
	"github.com/libesz/poolmanager/pkg/io"
	"github.com/libesz/poolmanager/pkg/webui/content"
	"github.com/urfave/negroni"
)

var parsedTemplates *template.Template

func New(listenOn, password string, configStore *configstore.ConfigStore, inputs []io.Input, outputs []io.Output) WebUI {
	r := mux.NewRouter()

	signingKey := []byte(uuid.New().String())
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			log.Println("WebUI: unauthorized request:", err)
			if !strings.Contains(r.URL.Path, "api/") {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				urlToRedirect := r.Header.Get("X-Script-Name") + "/login"
				http.Redirect(w, r, urlToRedirect, 301)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		},
		Extractor: func(r *http.Request) (string, error) {
			if !strings.Contains(r.URL.Path, "api/") {
				tokenFromCookie, err := r.Cookie("token")
				if err != nil {
					return "", fmt.Errorf("Web client shall use cookie tokens")
				}
				//log.Println("Extractor: returning token with cookie value:", tokenFromCookie.Value)
				return tokenFromCookie.Value, nil
			}
			//log.Println("Extractor: passing to FromAuthHeader")
			return jwtmiddleware.FromAuthHeader(r)
		},
	})

	r.PathPrefix("/").Handler(http.FileServer(content.Content))

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginPostHandler(signingKey, password, w, r)
	}).Methods("POST")

	r.Handle("/api/config", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			configPostHandler(configStore, w, r)
		})),
	)).Methods("POST")

	r.Handle("/api/ping", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { myHandler(w, r) })),
	))

	r.Handle("/api/status", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiStatusHandler(configStore, inputs, outputs, w, r)
		})),
	)).Methods("GET")

	server := &http.Server{
		Handler:      r,
		Addr:         listenOn,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return WebUI{server: server}
}

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	fmt.Fprintf(w, "This is an authenticated request")
	fmt.Fprintf(w, "Claim content:\n")
	for k, v := range user.(*jwt.Token).Claims.(jwt.MapClaims) {
		fmt.Fprintf(w, "%s :\t%#v\n", k, v)
	}
})

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

type ApiStatusResponse struct {
	Inputs  map[string]string `json:"inputs"`
	Outputs map[string]bool   `json:"outputs"`
}

func apiStatusHandler(configStore *configstore.ConfigStore, inputs []io.Input, outputs []io.Output, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := ApiStatusResponse{
		Inputs:  make(map[string]string),
		Outputs: make(map[string]bool),
	}
	for _, item := range inputs {
		if item.Value() == io.InputError {
			data.Inputs[item.Name()] = "N/A"
		} else {
			data.Inputs[item.Name()] = fmt.Sprintf("%.2f %s", item.Value(), item.Degree())
		}
	}

	for _, item := range outputs {
		data.Outputs[item.Name()] = item.Get()
	}

	log.Printf("Webui: rendering apiStatusHandler page with data: %+v\n", data)
	_ = json.NewEncoder(w).Encode(data)
}

type ConfigRequest struct {
	Controller string `json:"controller"`
	Type       string `json:"type"`
	Key        string `json:"key"`
	Value      string `json:"value"`
}

type ConfigResponse struct {
	Error     string      `json:"error"`
	OrigValue interface{} `json:"origValue"`
}

func configPostHandler(configStore *configstore.ConfigStore, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var data ConfigRequest
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
	u := ConfigResponse{Error: ""}
	_ = json.NewEncoder(w).Encode(u)
}

func respondError(w http.ResponseWriter, origValue interface{}, err error) {
	w.WriteHeader(http.StatusConflict)
	u := ConfigResponse{Error: err.Error(), OrigValue: origValue}
	_ = json.NewEncoder(w).Encode(u)
}

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Error string `json:"error"`
	Token string `json:"token"`
}

func generateJWT(signingKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": uuid.New(),
	})

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func loginPostHandler(signingKey []byte, password string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Printf("Webui: requested login\n")
	decoder := json.NewDecoder(r.Body)
	var data LoginRequest
	err := decoder.Decode(&data)
	if err != nil {
		log.Printf("Webui: decode error on login request: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		u := ConfigResponse{Error: err.Error(), OrigValue: nil}
		_ = json.NewEncoder(w).Encode(u)
		return
	}

	if data.Password != password {
		log.Printf("Webui: user unauthorized\n")
		w.WriteHeader(http.StatusUnauthorized)
		u := LoginResponse{Error: "Unauthorized"}
		_ = json.NewEncoder(w).Encode(u)
		return
	}
	log.Printf("Webui: user authorized\n")
	tokenString, _ := generateJWT(signingKey)
	w.WriteHeader(http.StatusAccepted)
	u := LoginResponse{Token: tokenString}
	_ = json.NewEncoder(w).Encode(u)
}
