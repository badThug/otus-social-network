package web

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"runtime/debug"
	"service-users/app/customErrors"
	"service-users/app/handlers"
	"service-users/app/utils"
)

type Dispatcher struct {
	Router *mux.Router
}

func InitDispatcher() Dispatcher {
	router := mux.NewRouter()
	router.NotFoundHandler = NotFoundHandler

	dispatcher := Dispatcher{
		Router: router,
	}

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
		h := handlers.InitHandler(w, r)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered in f: %w; Stack trace %s", r, string(debug.Stack()))
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
				log.Printf("Error occured: %+v\n", err)
			}
		}
	}
}

var NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	utils.SendResponseJsonWithStatusCode(w, utils.ResponseErrorMessage("This resources was not found on our server"), http.StatusNotFound)
})
