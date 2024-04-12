package routes

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/vladwithcode/sibra-cdln/internal/db"
	"github.com/vladwithcode/sibra-cdln/internal/whatsapp"
)

func RegisterAppointmentRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/appointments", CreateNewAppointment)
}

func CreateNewAppointment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Printf("Parse form err: %v\n", err)
		templ, err := template.ParseFiles("web/templates/appointment-form.html")

		if err != nil {
			fmt.Printf("Parse Files err: %v\n", err)
			respondWithError(w, http.StatusInternalServerError, ErrorParams{})
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(500)
		err = templ.ExecuteTemplate(
			w,
			"appointment-inputs",
			map[string]string{
				"ErrorMessage": "El formularío contiene información invalida",
				"ElementId":    "appointment-result",
			},
		)

		if err != nil {
			fmt.Printf("Execute Template err: %v\n", err)
			respondWithError(w, http.StatusInternalServerError, ErrorParams{})
		}
		return
	}

	var (
		strDate = r.Form.Get("date")
		name    = r.Form.Get("name")
		phone   = r.Form.Get("phone")
	)

	date := time.Time{}

	err = date.UnmarshalText([]byte(strDate + "-06:00"))

	if err != nil {
		fmt.Printf("Unmarshal date: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, ErrorParams{})
		return
	}

	phone, err = formatPhone(phone)

	if err != nil {
		fmt.Printf("Parse form err: %v\n", err)
		templ, err := template.ParseFiles("web/templates/appointment-form.html")

		if err != nil {
			fmt.Printf("Parse Files err: %v\n", err)
			respondWithError(w, http.StatusInternalServerError, ErrorParams{})
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(500)
		err = templ.ExecuteTemplate(
			w,
			"appointment-inputs",
			map[string]string{
				"ErrorMessage": "El formularío contiene información invalida",
				"ElementId":    "appointment-result",
			},
		)

		if err != nil {
			fmt.Printf("Execute Template err: %v\n", err)
			respondWithError(w, http.StatusInternalServerError, ErrorParams{})
		}
		return
	}

	id, err := uuid.NewV7()

	if err != nil {
		fmt.Printf("UUID err: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, ErrorParams{})
		return
	}

	appointment := db.Appointment{
		Id:        id.String(),
		Name:      name,
		Phone:     phone,
		Date:      date,
		CreatedAt: time.Now(),
		Attended:  false,
	}

	// Send whatsapp messages
	err = whatsapp.SendNewAppointmentNotification(&appointment)

	if err != nil {
		fmt.Printf("Send message err: %v\n", err)
		templ, err := template.ParseFiles("web/templates/responses/request-error.html")

		if err != nil {
			fmt.Printf("Parse Files err: %v\n", err)
			respondWithError(w, http.StatusInternalServerError, ErrorParams{})
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(500)
		err = templ.Execute(
			w,
			map[string]string{
				"ErrorMessage": "No fue posible generar la cita, intente de nuevo más tarde",
				"ElementId":    "appointment-result",
			},
		)

		if err != nil {
			fmt.Printf("Execute Template err: %v\n", err)
			respondWithError(w, http.StatusInternalServerError, ErrorParams{})
		}

		return
	}

	// Save appointment to DB

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(201)
	templ, err := template.ParseFiles("web/templates/responses/request-success.html")

	if err != nil {
		fmt.Printf("Parse Files err: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, ErrorParams{})
		return
	}

	err = templ.Execute(w, nil)

	if err != nil {
		fmt.Printf("Execute Files err: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, ErrorParams{})
		return
	}
}

func formatPhone(p string) (string, error) {
	stripCountryCodeExp := regexp.MustCompile(`^\+52`)
	replaceExp := regexp.MustCompile(`[ -]`)
	numExp := regexp.MustCompile(`[0-9]{10}`)

	phone := stripCountryCodeExp.ReplaceAll([]byte(p), []byte(""))
	phone = replaceExp.ReplaceAll([]byte(phone), []byte(""))

	if !numExp.Match(phone) {
		return "", errors.New(fmt.Sprintf("The string is not a valid phone number: %v", p))
	}

	return string(phone), nil
}
