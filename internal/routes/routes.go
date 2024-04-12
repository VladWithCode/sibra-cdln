package routes

import (
	"fmt"
	"html/template"
	"net/http"
)

func NewRouter() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", RenderIndex)

	// Serve static
	// *Might change it to serve static files through nginx
	fs := http.FileServer(http.Dir("web/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Appointment Routes
	RegisterAppointmentRoutes(router)

	return router
}

func RenderIndex(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/index.html",
		"web/templates/contact-form.html",
		"web/templates/appointment-form.html",
	)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		w.WriteHeader(500)
		w.Write([]byte("Ocurrio un error en el servidor"))
		return
	}

	err = templ.Execute(w, nil)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		respondWithError(w, 500, ErrorParams{})
	}
}

type ErrorParams struct {
	ErrorMessage string
	ElementId    string
	ElementClass string
}

func respondWithError(w http.ResponseWriter, code int, params ErrorParams) {
	w.Header().Set("Content-Type", "text/plain")
	if code == http.StatusInternalServerError {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ocurrio un error inesperado en el servidor"))
		return
	}

	templ, err := template.ParseFiles("web/templates/responses/request-error.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ocurrio un error inesperado en el servidor"))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	err = templ.Execute(w, params)

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ocurrio un error inesperado en el servidor"))
		return
	}
}
