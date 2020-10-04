package web

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"service-users/app/config"
	"service-users/app/customErrors"
	"service-users/app/handlers"
	"service-users/app/storage"
	"service-users/app/utils"
)

type Dispatcher struct {
	Router *mux.Router
	db     *storage.DbConnection
}

func InitDispatcher(db *storage.DbConnection, config *config.Config) Dispatcher {
	router := mux.NewRouter()

	dispatcher := Dispatcher{
		Router: router,
		db:     db,
	}

	handlers.CreateValidator()
	initRoutes(dispatcher)

	return dispatcher
}

// Get wraps the router for GET method
func (d *Dispatcher) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	d.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (d *Dispatcher) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	d.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (d *Dispatcher) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	d.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (d *Dispatcher) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	d.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (d *Dispatcher) Run(host string) {
	log.Fatal(http.ListenAndServe(host, d.Router))
}

func (d *Dispatcher) handleRequest(handlerMethod func(h *handlers.Handler) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := handlers.InitHandler(d.db, w, r)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered in f: %w", r)
				h.ResponseWithError("Internal server error", http.StatusInternalServerError)
			}
		}()

		log.Printf("%s %s", r.Method, r.URL)

		if err := handlerMethod(h); nil != err {
			switch causedErr := errors.Cause(err).(type) {
			case *utils.MalformedRequest:
				h.ResponseWithError(causedErr.Msg, causedErr.Status)
			case *customErrors.TypedError:
				h.ResponseWithError(causedErr.Msg, http.StatusBadRequest)
			case *customErrors.TypedStatusError:
				h.ResponseWithError(causedErr.Msg, causedErr.Status)
			default:
				h.ResponseWithError("Internal server error", http.StatusInternalServerError)
			}

			log.Printf("Error occured: %+v\n", err)
		}
	}
}
