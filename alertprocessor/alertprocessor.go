package alertprocessor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/alertmanager/notify/webhook"
)

type AlertProcessor struct {
}

func NewAlertProcessor() (*AlertProcessor, error) {
	return &AlertProcessor{}, nil
}

func (ap *AlertProcessor) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Only HTTP Post is allowed
	if req.Method != http.MethodPost {
		log.Printf("Unsupported HTTP method: %s", req.Method)
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Failed to read request body: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parse alertmanager alert
	alert := &webhook.Message{}
	if err := json.Unmarshal(requestBody, alert); err != nil {
		log.Printf("Failed to parse Alertmanager alert: %s", err)
		log.Printf("%s", string(requestBody))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Process alert
	if err := ap.processAlert(alert); err != nil {
		log.Printf("Failed to process alert: %s", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)

}

func (ap *AlertProcessor) processAlert(alert *webhook.Message) error {
	return nil
}
