package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/dhickie/hickhub/config"
	"github.com/dhickie/hickhub/log"
	"github.com/nats-io/go-nats"
)

const subjectPath = "/user/messaging/subject"

var errCantGetSubject = errors.New("Unable to get messaging subject from API")

// The messenger module subscribes to messages from the cloud NATS server,
// calling methods on the API when it receives them.
type module struct {
	APIPort   int
	APIClient *http.Client
	NatsConn  *nats.Conn
	NatsSub   *nats.Subscription
}

type hickHubMessage struct {
	ID   int    `json:"id"`
	Data []byte `json:"data"` // message object encoded as binary
}

type message struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   string `json:"body"`
	Status string `json:"status"`
}

type subjectResponse struct {
	Subject string `json:"subject"`
}

var mod module

// Launch launches the messenger module.
func Launch(appConfig config.Config) {
	// Connect to the NATS server
	log.Info("Launching messenger module")
	log.Info("Connecting to messenger server")
	server := appConfig.Messaging.MessagingServer
	authToken := appConfig.Messaging.AuthKey
	nc, err := nats.Connect(server, nats.Token(authToken))
	if err != nil {
		panic("Failed to conntect to messaging server")
	}

	// Query the HickHub API to get this hub's messaging subject
	log.Info("Getting messaging subject from API")
	subj, err := getMessagingSubject(appConfig.Messaging.APIServer, authToken)
	if err != nil {
		panic("Failed to get messaging subject")
	}

	// Create the API client
	mod = module{
		APIPort:   appConfig.API.Port,
		APIClient: http.DefaultClient,
		NatsConn:  nc,
	}

	// Subscribe to the topic
	log.Info("Subscribing to messenger subject")
	sub, err := nc.Subscribe(subj, mod.internetSubscriber)
	if err != nil {
		panic("Failed to subscribe to messaging topic")
	}
	mod.NatsSub = sub
}

func getMessagingSubject(apiURL, authToken string) (string, error) {
	client := new(http.Client)
	path := apiURL + subjectPath
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errCantGetSubject
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	subjectResponse := new(subjectResponse)
	err = json.Unmarshal(body, subjectResponse)
	if err != nil {
		return "", err
	}

	return subjectResponse.Subject, nil
}

func (module *module) internetSubscriber(m *nats.Msg) {
	// Decode the hickhub message
	hhMsg := new(hickHubMessage)
	err := json.Unmarshal(m.Data, hhMsg)
	if err != nil {
		log.Error(fmt.Sprintf("An error occured processing an internet message: %v", err.Error()))
		return
	}

	// Decode the actual request
	msg := new(message)
	err = json.Unmarshal(hhMsg.Data, msg)
	if err != nil {
		log.Error(fmt.Sprintf("An error occured processing the internet request: %v", err.Error()))
		return
	}

	// Make the correct request for the API module
	url := fmt.Sprintf("http://localhost:%v/api/%v", module.APIPort, msg.Path)
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
		log.Error(fmt.Sprintf("An error occured trying to mashal a reply JSON: %v", err.Error()))
		return
	}

	hhReply := hickHubMessage{
		ID:   hhMsg.ID,
		Data: msgJSON,
	}
	replyJSON, err := json.Marshal(hhReply)
	if err != nil {
		log.Error(fmt.Sprintf("An error occured trying to marshal the reply message: %v", err.Error()))
		return
	}

	err = mod.NatsConn.Publish(m.Reply, replyJSON)
	if err != nil {
		log.Error(fmt.Sprintf("An error occured trying to publish a reply: %v", err.Error()))
		return
	}
}
