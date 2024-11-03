package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/valyala/fasthttp"
	"os"
)

const (
	charset   = "abcdefghijklmnopqrstuvwxyz0123456789"
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"
	Delay     = 2 * time.Minute // session time for 1 connection
)

func extractValue(response string, key string) string {
	startKey := `\"` + key + `\":\"`
	start := strings.Index(response, startKey) + len(startKey)
	if start < len(startKey) {
		return ""
	}
	end := strings.Index(response[start:], `\"`) + start
	if end < start {
		return ""
	}
	return response[start:end]
}

func createClientOptions(reconnectFunc func()) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://popdeng.click/")
	req.Header.SetUserAgent(UserAgent)
	resp := fasthttp.AcquireResponse()
	if fasthttp.Do(req, resp) != nil {
		fmt.Println("Get info fail")
		return nil
	}
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	body := string(resp.Body())

	opts.AddBroker(extractValue(body, "host"))
	opts.SetHTTPHeaders(http.Header{
		"origin":     {"https://popdeng.click"},
		"user-agent": {UserAgent},
	})
	opts.SetUsername(extractValue(body, "user"))
	opts.SetPassword(extractValue(body, "token"))
	opts.SetClientID(extractValue(body, "clientId"))
	opts.SetCleanSession(false)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	})
	opts.OnConnect = func(client mqtt.Client) {
		fmt.Println("Connected")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v\n", err)
		// Reconnect immediately
		go reconnectFunc()
	}
	return opts
}

func Connect() {
	for {
		// Define a function to call Connect recursively on disconnect
		reconnect := func() {
			fmt.Println("Attempting to reconnect ✅✅✅")
			Connect() // Recursive call to attempt reconnection
		}

		opts := createClientOptions(reconnect)
		if opts == nil {
			time.Sleep(time.Second * 2)
			continue
		}
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Println("Connect Fail", token.Error())
			time.Sleep(time.Second)
			continue
		}

		client.Subscribe("popdeng/clicks", 1, nil)
		quitT := time.Now().Add(Delay)
		for time.Now().Before(quitT) {
			for i := 0; i < 3; i++ { // Publish 1000 times
				token := client.Publish("popdeng/clicks", 0, false, `OZKLKWVMTGVHG8CCC!.8vrr`)
				if token.Wait() && token.Error() != nil {
					fmt.Println("send click error", token.Error())
					break
				}
			}
			time.Sleep(time.Millisecond * 25) // Wait 1 second before the next 1000 publishes
		}

		client.Disconnect(5)
		break
	}
}

func main() {
	// Run the MQTT connections
	for i := 0; i < 5; i++ {
		go Connect()
		time.Sleep(time.Millisecond * 25)
	}

	// Set up HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running")
	})

	go func() {
		fmt.Printf("Server is listening on port %s\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	select {} // Keep the main goroutine running
}
