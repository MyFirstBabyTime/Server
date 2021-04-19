package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// aligoAgent is struct that agent API about aligo including sending message, get message list, etc ...
type aligoAgent struct {
	apiKey, id, sender string
}

func AligoAgent(apiKey, id, sender string) *aligoAgent {
	return &aligoAgent{
		apiKey: apiKey,
		id:     id,
		sender: sender,
	}
}

// SendSMSToOne method send SMS message to one receiver
func (aa *aligoAgent) SendSMSToOne(receiver, content string) (err error) {
	return aa.sendMsgToReceivers([]string{receiver}, "", content, "SMS")
}

func (aa *aligoAgent) sendMsgToReceivers(receivers []string, title, content, _type string) (err error) {
	req, err := http.NewRequest("POST", "https://apis.aligo.in/send/", nil)
	if err != nil {
		err = errors.New(fmt.Sprintf("some error occurs while creating request, err: %v", err))
		return
	}

	q := req.URL.Query()
	q.Add("key", aa.apiKey)
	q.Add("user_id", aa.id)
	q.Add("sender", aa.sender)
	if _type != "" {
		q.Add("msg_type", _type)
	}
	if title != "" {
		q.Add("title", title)
	}
	q.Add("receiver", strings.Join(receivers, ","))
	q.Add("msg", content)
	req.URL.RawQuery = q.Encode()

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("some error occurs while sending request, err: %v", err))
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("aligo API dosen't return 200, status code: %d", resp.StatusCode))
		return
	}

	respBody := struct {
		Code int    `json:"result_code"`
		Msg  string `json:"message"`
	}{}
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	if respBody.Code != 0 {
		err = errors.New(fmt.Sprintf("aligo API return unexpected result code, result code: %d, message: %s", respBody.Code, respBody.Msg))
		return
	}
	return
}
