package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"net/http"
	"time"
)

func SendLog(standardLog model.Log) {
	var bearer = config.AuthenticationToken()

	jsonMarshalledLog, err := json.Marshal(standardLog)
	if err != nil {
		log.Logger.Error().Msgf("Error marshalling JSON: %s", err)
	}

	requestBody := bytes.NewReader(jsonMarshalledLog)

	var apiURL = fmt.Sprintf("http://%s/api/addLog", config.Config.C2)

	request, err := http.NewRequest("POST", apiURL, requestBody)
	if err != nil {
		log.Logger.Info().Msgf("[!]\tError on request.\n[ERROR] -  \n", err)
	}

	request.Header.Add("Authorization", bearer)

	client := http.Client{
		Timeout: time.Second * 3,
	}
	response, err := client.Do(request)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)
		log.Logger.Debug().Msgf("Attempting to resend log")
		time.Sleep(time.Second * 1)
		SendLog(standardLog)
	}

	err = response.Body.Close()
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError closing response body.\n[ERROR] -  \n", err)
	}

	log.Logger.Debug().Msgf("Log successfully sent!")
}
