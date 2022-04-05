package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/sindreslungaard/assert"
)

type JSON map[string]interface{}

func Res(w http.ResponseWriter, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		Warn("Failed to json marshal API response", payload)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func DecodeRequestBody(r *http.Request, dst interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		return err
	}
	return nil
}

func ApiListenAndServe(port int) {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	r.Get("/api/deployments", GetDeploymentsHandler)
	r.Post("/api/deployments", PostDeploymentsHandler)

	Info(fmt.Sprintf("Listening on port %v..", port))

	Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))

}

func GetDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	config := ReadConfig()
	deployments := []Deployment{}

	for _, d := range config.Deployments {
		deployments = append(deployments, d)
	}

	Res(w, JSON{"deployments": deployments})
}

func PostDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	var d Deployment

	if err := DecodeRequestBody(r, &d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid, err1 := assert.Is(d.UID).NotEmpty().String()
	_, err2 := assert.Is(d.Image).NotEmpty().String()

	err := assert.First(err1, err2)

	if err != nil || d.Ports == nil || d.Env == nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	config := ReadConfig()

	_, ok := config.Deployments[uid]

	// deployment uid already taken
	if ok {
		Res(w, JSON{"error": fmt.Sprintf("Deployment name \"%s\" is already in use", uid)})
		return
	}

	UpdateConfig(func(c Config) Config {
		c.Deployments[uid] = d

		return c
	})

	RunContainerFromDeployment(Docker(), d)

	Res(w, d)
}
