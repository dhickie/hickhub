package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/dhickie/hickhub/config"
	"github.com/dhickie/hickhub/log"
	"github.com/nats-io/go-nats"
)

// The messenger module subscribes to messages from the cloud NATS server,
// calling methods on the API when it receives them.
type module struct {
	APIPort   int
	APIClient *http.Client
	NatsConn  *nats.Conn
	NatsSub   *nats.Subscription
}

type message struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   string `json:"body"`
	Status string `json:"status"`
}

var mod module

// Launch launches the messenger module.
func Launch(appConfig config.Config) {
	// Connect to the NATS server
	server := appConfig.Messaging.Server
	nc, err := nats.Connect(server)
	if err != nil {
		panic("Failed to conntect to messaging server")
	}

	// Create the API client
	mod = module{
		APIPort:   appConfig.API.Port,
		APIClient: http.DefaultClient,
		NatsConn:  nc,
	}

	// Subscribe to the topic
	sub, err := nc.Subscribe(appConfig.Auth.Key, mod.internetSubscriber)
	if err != nil {
		panic("Failed to subscribe to messaging topic")
	}
	mod.NatsSub = sub
}

func (module *module) internetSubscriber(m *nats.Msg) {
	// Decode the message
	msg := new(message)
	err := json.Unmarshal(m.Data, msg)
	if err != nil {
		log.Error(fmt.Sprintf("An error occured processing an internet message: %v", err.Error()))
		return
	}

	// Make the correct request for the API module
	url := fmt.Sprintf("http://localhost:%v/%v", module.APIPort, msg.Path)
	request, err := http.NewRequest(msg.Method, url, bytes.NewBuffer([]byte(msg.Body)))
	if err != nil {
		log.Error(fmt.Sprintf("An error creating the API request: %v", err.Error()))
		return
	}

	// Do the request
	resp, err := module.APIClient.Do(request)
	if err != nil {
		log.Error(fmt.Sprintf("An error occured making the call to the HickHub API: %v", err.Error()))
		return
	}

	// Build the reply message (only need to populate body and status for replies)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		log.Error(fmt.Sprintf("An error occured reading the API response body: %v", err.Error()))
		return
	} else if err == io.EOF {
		body = []byte("")
	}

	reply := message{
		Body:   string(body),
		Status: resp.Status,
	}
	msgJSON, err := json.Marshal(reply)
	if err != nil {
		log.Error(fmt.Sprintf("An error occured trying to mashal a reply message: %v", err.Error()))
		return
	}

	err = mod.NatsConn.Publish(m.Reply, msgJSON)
	if err != nil {
		log.Error(fmt.Sprintf("An error occured trying to publish a reply: %v", err.Error()))
		return
	}
}
