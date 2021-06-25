package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/glbter/go-genesis-ses-2021/config"
	"github.com/glbter/go-genesis-ses-2021/model"
	"github.com/glbter/go-genesis-ses-2021/util"
)

// btcRate
var ExRate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	respData := SendRequest(config.GetConfig().NbuLink + config.GetConfig().NbuParam)

	var respObjNbu []model.NbuData
	json.Unmarshal([]byte(respData), &respObjNbu)
	uahToUsd := respObjNbu[0].Rate

	respData = SendRequest(config.GetConfig().BlockchainLink + config.GetConfig().BlockchainParam)

	var respObjBlockchain model.BlockchainResponse
	json.Unmarshal([]byte(respData), &respObjBlockchain)
	blockchainToUsd := respObjBlockchain.Last_trade_price

	converted := float64(blockchainToUsd) * uahToUsd
	prettyResult := util.Round(converted, 0.01)

	w.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(prettyResult)
	w.Write([]byte(payload))
})

func SendRequest(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return nil
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return respData
}
