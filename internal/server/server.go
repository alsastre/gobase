package server

import (
	"context"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/alsastre/gobase/internal/data"
	"github.com/alsastre/gobase/internal/routes"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// Server represents the structure of the server with all the common config
type Server struct {
	ViperCfg *viper.Viper // Configuration
	Logger   *zap.Logger  // Logger to be used
	Data     *data.Data   // Data with dbSession included
}

// The type key is used to prevent collision in the context
type key int

// Context key used to pass the object of the request in the SomethingMiddleware
const somethingKey key = 0

// Router initializes the routes of the server for the somethings
func (server *Server) Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", server.ListSomething)
	r.Post("/", server.CreateSomething)      // POST /somethings
	r.Get("/search", server.SearchSomething) // GET /somethings/search

	// {SomethingID} will be filled with the ID after the /
	r.Route("/{SomethingID}", func(r chi.Router) {
		r.Use(server.SomethingMiddleware)     // Load the *Something on the request context
		r.Get("/", server.GetSomething)       // GET /somethings/123
		r.Put("/", server.UpdateSomething)    // PUT /somethings/123
		r.Delete("/", server.DeleteSomething) // DELETE /somethings/123
	})

	return r
}

// ListSomething request the data framework for the list of items and returns them
func (server *Server) ListSomething(w http.ResponseWriter, r *http.Request) {
	lb, err := NewSomethingResponseListResponse(server.Data.ListSomethings())
	if err != nil {
		render.Render(w, r, routes.ErrInternalServer(err))
	}
	render.RenderList(w, r, lb)
}

// CreateSomething ...
func (server *Server) CreateSomething(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("CreateSomething: NOT IMPLEMENTED"))
}

// SearchSomething ...
func (server *Server) SearchSomething(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SearchSomethings: NOT IMPLEMENTED"))
}

// GetSomething ...
func (server *Server) GetSomething(w http.ResponseWriter, r *http.Request) {
	// Retrive the object from the context thanks to the SomethingMiddleware
	some := r.Context().Value(somethingKey).(*data.Something)
	// Prepare a response with the object
	resp, _ := NewSomethingResponse(some, nil)
	// Respond with the object rendered
	render.Render(w, r, resp)
}

// UpdateSomething ...
func (server *Server) UpdateSomething(w http.ResponseWriter, r *http.Request) {
	some := r.Context().Value(somethingKey).(data.Something)
	w.Write([]byte("UpdateSomething:" + some.Name + " NOT IMPLEMENTED"))
}

// DeleteSomething ...
func (server *Server) DeleteSomething(w http.ResponseWriter, r *http.Request) {
	some := r.Context().Value(somethingKey).(data.Something)
	w.Write([]byte("DeleteSomething:" + some.Name + " NOT IMPLEMENTED"))
}

// SomethingMiddleware obtains the corresponding Something object from the DB. If not found, we return 404
func (server *Server) SomethingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var mySomething *data.Something
		var err error
		if ID := chi.URLParam(r, "SomethingID"); ID != "" {
			mySomething, err = server.Data.GetSomething(ID)
		}
		if err != nil {
			render.Render(w, r, routes.ErrNotFound(err))
			return
		}

		// Fill in the context with the value
		ctx := context.WithValue(r.Context(), somethingKey, mySomething)

		// Continue with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SomethingResponse is the response payload for the Something data model
// This struct may be useful in case the reponse must include more fields apart from the object itself
type SomethingResponse struct {
	*data.Something
}

// NewSomethingResponse builds a response from a obj Something
func NewSomethingResponse(some *data.Something, err error) (*SomethingResponse, error) {
	if err != nil {
		return nil, err
	}
	resp := &SomethingResponse{Something: some}
	return resp, nil
}

// Render allow us to meet the render interface and in the future it may be usefull to modify/ommit some data while generating a response
// This is used by the chi render
func (rd *SomethingResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// NewSomethingResponseListResponse builds a response from a list of somethings
func NewSomethingResponseListResponse(somes []*data.Something, err error) ([]render.Renderer, error) {
	if err != nil {
		return nil, err
	}

	list := []render.Renderer{}
	for _, some := range somes {
		item, _ := NewSomethingResponse(some, nil)
		list = append(list, item)
	}

	return list, nil
}
