package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/valyala/fasthttp"
)

const (
	charset   = "abcdefghijklmnopqrstuvwxyz0123456789"
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"
	Delay     = 35 * time.Second // session time for 1 connection
)

func extractValue(response string, key string) string { // thanks claude.ai
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

func createClientOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://www.popdeng.click/")
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
		"origin":     {"https://www.popdeng.click"},
		"user-agent": {UserAgent},
	})
	opts.SetUsername(extractValue(body, "user"))
	opts.SetPassword(extractValue(body, "token"))
	opts.SetClientID(extractValue(body, "clientId"))
	opts.SetCleanSession(true)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	})
	opts.OnConnect = func(client mqtt.Client) {
		fmt.Println("Connected")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v\n", err)
	}
	return opts
}

func Connect() {
	for {
		opts := createClientOptions()
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
		// client.Subscribe("popdeng/leaderboard", 1, nil)
		// client.Subscribe("popdeng/clicks", 1, nil)
		quitT := time.Now().Add(Delay)
		for time.Now().Before(quitT) {
			// token := client.Publish("popdeng/leaderboard", 0, false, `[{"country_code":"WW","clicks":0,"cps":0},{"country_code":"KP","clicks":6969696969,"cps":6969696969},{"country_code":"IS","clicks":6969696969,"cps":6969696969},{"country_code":"WW","clicks":0,"cps":0},{"country_code":"HK","clicks":69,"cps":1},{"country_code":"AT","clicks":69,"cps":2},{"country_code":"CN","clicks":69,"cps":3},{"country_code":"KW","clicks":69,"cps":4},{"country_code":"ZZ","clicks":0,"cps":0},{"country_code":"BE","clicks":69,"cps":5},{"country_code":"YE","clicks":69,"cps":6},{"country_code":"ZZ","clicks":0,"cps":0},{"country_code":"LR","clicks":69,"cps":7},{"country_code":"AQ","clicks":69,"cps":8},{"country_code":"RO","clicks":69,"cps":9},{"country_code":"IS","clicks":69,"cps":10},{"country_code":"NO","clicks":69,"cps":11},{"country_code":"AO","clicks":69,"cps":12},{"country_code":"ZZ","clicks":0,"cps":0},{"country_code":"TH","clicks":1000,"cps":1},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{}]`)
			token := client.Publish("popdeng/clicks", 0, false, `RHKLKWVMTGVHG8CCC!.8h4n`)
			if token.Wait() && token.Error() != nil {
				fmt.Println("send click error", token.Error())
				break
			}
			time.Sleep(time.Millisecond)
		}
		client.Disconnect(5)
	}
}

func main() {
	for i := 0; i < 50; i++ {
		go Connect()
		time.Sleep(time.Millisecond * 25)
	}
	select {}
}
