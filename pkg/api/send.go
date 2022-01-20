package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
)

func SendLog(standardLog model.Log) error {
	var bearer = config.AuthenticationToken()

	jsonMarshalledLog, err := json.Marshal(standardLog)
	if err != nil {
		log.Logger.Error().Msgf("Error marshalling JSON: %s", err)
		return err
	}

	requestBody := bytes.NewReader(jsonMarshalledLog)

	C2Host := config.Config.C2
	C2Protocol := config.Config.C2Protocol
	C2Port := config.Config.Port
	var apiURL = fmt.Sprintf("%s://%s:%d/api/logs/addLog", C2Protocol, C2Host, C2Port)

	request, err := http.NewRequest("POST", apiURL, requestBody)
	if err != nil {
		log.Logger.Info().Msgf("[!]\tError on request.\n[ERROR] -  \n", err)
		return err
	}

	request.Header.Add("Authorization", bearer)

	client := http.Client{
		Timeout: time.Second * 3,
	}
	response, err := client.Do(request)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)
		log.Logger.Debug().Msgf("Attempting to resend log")
		return err
	}

	err = response.Body.Close()
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError closing response body.\n[ERROR] -  \n", err)
	}

	log.Logger.Debug().Msgf("Log successfully sent!")
	return nil
}
