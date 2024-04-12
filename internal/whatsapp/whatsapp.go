package whatsapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/vladwithcode/sibra-cdln/internal"
	"github.com/vladwithcode/sibra-cdln/internal/db"
)

type TemplateVar map[string]interface{}

type TemplateComponent struct {
	ComponentType string        `json:"type"`
	Parameters    []TemplateVar `json:"parameters"`
}

type TemplateData struct {
	TemplateName string
	BodyVars     []TemplateVar
	HeaderVars   []TemplateVar
}

type templatePayload struct {
	Name       string              `json:"name"`
	Components []TemplateComponent `json:"components"`
	Language   struct {
		Code string `json:"code"`
	} `json:"language"`
}

func postCloudAPIMessage(requestPayload interface{}) error {
	phoneNumberId := os.Getenv("FB_PHONE_NUMBER_ID")
	fbAccessToken := os.Getenv("FB_ACCESS_TOKEN")

	if phoneNumberId == "" {
		panic("FB_PHONE_NUMBER_ID not in ENV")
	}
	if fbAccessToken == "" {
		panic("FB_ACCESS_TOKEN not in ENV")
	}

	reqBody, err := json.Marshal(requestPayload)

	if err != nil {
		return err
	}

	reqUrl := fmt.Sprintf("https://graph.facebook.com/v18.0/%v/messages", phoneNumberId)
	req, err := http.NewRequest("post", reqUrl, bytes.NewReader(reqBody))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", fbAccessToken))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		fmt.Printf("Errored with status: %v\n", resp.Status)
		fmt.Printf("resp.Header: %v\n", resp.Header["Www-Authenticate"])

		return errors.New(fmt.Sprintf("Request failed with code: %v", resp.StatusCode))
	}

	return nil
}

func SendTemplateMessage(phoneNumber string, data TemplateData) error {
	var reqPayload struct {
		MessagingProduct string          `json:"messaging_product"`
		MessageType      string          `json:"type"`
		ToPhone          string          `json:"to"`
		Template         templatePayload `json:"template"`
	}

	var components []TemplateComponent

	if len(data.HeaderVars) != 0 {
		components = append(components, TemplateComponent{
			ComponentType: "header",
			Parameters:    data.HeaderVars,
		})
	}

	if len(data.BodyVars) != 0 {
		components = append(components, TemplateComponent{
			ComponentType: "body",
			Parameters:    data.BodyVars,
		})
	}

	reqPayload.MessageType = "template"
	reqPayload.MessagingProduct = "whatsapp"
	reqPayload.ToPhone = phoneNumber
	reqPayload.Template = templatePayload{
		Name: data.TemplateName,
		Language: struct {
			Code string `json:"code"`
		}{Code: "es"},
		Components: components,
	}

	return postCloudAPIMessage(reqPayload)
}

func SendNewAppointmentNotification(appointmentData *db.Appointment) error {
	sendToPhone := os.Getenv("APPOINTMENTS_PHONE")

	if sendToPhone == "" {
		return errors.New("SendTo phone is not in ENV")
	}

	var (
		headerVars []TemplateVar
		bodyVars   []TemplateVar
	)

	dateStr := internal.FormatDate(appointmentData.Date)
	createdDateStr := internal.FormatDate(appointmentData.CreatedAt)

	headerVars = append(headerVars, TemplateVar{
		"type": "text",
		"text": appointmentData.Name,
	})

	bodyVars = append(bodyVars, TemplateVar{
		"type": "text",
		"text": appointmentData.Name,
	})
	bodyVars = append(bodyVars, TemplateVar{
		"type": "text",
		"text": appointmentData.Phone,
	})
	bodyVars = append(bodyVars, TemplateVar{
		"type": "text",
		"text": dateStr,
	})
	bodyVars = append(bodyVars, TemplateVar{
		"type": "text",
		"text": createdDateStr,
	})

	return SendTemplateMessage(sendToPhone, TemplateData{
		TemplateName: "appointment_notification",
		BodyVars:     bodyVars,
		HeaderVars:   headerVars,
	})
}
