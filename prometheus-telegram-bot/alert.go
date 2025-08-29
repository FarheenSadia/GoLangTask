package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Alert struct {
	Status string `json:"status"`
	Alerts []struct {
		Labels      map[string]string `json:"labels"`
		Annotations map[string]string `json:"annotations"`
		StartsAt    string            `json:"startsAt"`
		EndsAt      string            `json:"endsAt"`
	} `json:"alerts"`
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	var alert Alert
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, a := range alert.Alerts {
		msg := fmt.Sprintf(
			"*Alert:* %s\n Instance: %s\n Severity: %s\n Summary: %s\nStatus: %s",
			a.Labels["alertname"],
			a.Labels["instance"],
			a.Labels["severity"],
			a.Annotations["summary"],
			alert.Status,
		)

		sendToTelegram(msg)
	}
	w.WriteHeader(http.StatusOK)
}
