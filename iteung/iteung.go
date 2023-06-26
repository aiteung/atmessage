package iteung

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aiteung/atmessage"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func Send(iteungIPaddress string, im *atmessage.IteungMessage, toJID types.JID, waclient *whatsmeow.Client) {
	url := "http://" + iteungIPaddress + "/iteung/chatbot"
	method := "POST"
	msgs, err := json.Marshal(&im)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(msgs))
	payload := strings.NewReader(string(msgs))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	var sendtosender atmessage.IteungRespon
	Data := []byte(string(body))
	json.Unmarshal(Data, &sendtosender)
	if sendtosender.Message != "" {
		atmessage.SendMessage(sendtosender.Message, toJID, waclient)
	} else {
		fmt.Println("=======Python Backend Iteung Web : Respon Empty Message=======")
	}

}
