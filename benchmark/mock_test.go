package benchmark

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	MQTT "github.com/liaoliaopro/paho.mqtt.golang"
	"go.uber.org/zap"

	"github.com/fangwendong/surgemq/acl"
	"github.com/fangwendong/surgemq/service"
)

func Test1(t *testing.T) {
	var f MQTT.MessageHandler = func(MQTT.Client, MQTT.Message) {
		fmt.Println("rece")
	}

	connOpts := &MQTT.ClientOptions{
		ClientID:             "ds-live",
		CleanSession:         true,
		Username:             "fwd",
		Password:             "zzz",
		MaxReconnectInterval: 1 * time.Second,
		KeepAlive:            30 * time.Second,
		AutoReconnect:        true,
		PingTimeout:          10 * time.Second,
		ConnectTimeout:       30 * time.Second,
		TLSConfig:            tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert},
		OnConnectionLost:     func(c MQTT.Client, err error) { fmt.Println("mqtt disconnected.", zap.Error(err)) },
	}
	connOpts.AddBroker("tcp://127.0.0.1:8080")

	mc := MQTT.NewClient(connOpts)
	if token := mc.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("mqtt connection failed.", zap.Error(token.Error()))
		return
	}

	defer mc.Disconnect(12)
	if token := mc.Subscribe("ds", 0, f); token.Wait() && token.Error() != nil {
		fmt.Println("publish failed.", zap.Error(token.Error()))

	} else {
		t := token.(*MQTT.SubscribeToken)
		fmt.Println(t.Result())
	}

	if token := mc.Subscribe("dsx", 0, f); token.Wait() && token.Error() != nil {
		fmt.Println("publish failed.", zap.Error(token.Error()))

	} else {
		t := token.(*MQTT.SubscribeToken)
		fmt.Println(t.Result())
	}

	if token := mc.Subscribe("ds", 0, f); token.Wait() && token.Error() != nil {
		fmt.Println("publish failed.", zap.Error(token.Error()))

	} else {
		t := token.(*MQTT.SubscribeToken)
		fmt.Println(t.Result())
	}

	if token := mc.Unsubscribe("ds"); token.Wait() && token.Error() != nil {
		fmt.Println("publish failed.", zap.Error(token.Error()))

	}
	if token := mc.Subscribe("dsx", 0, f); token.Wait() && token.Error() != nil {
		fmt.Println("publish failed.", zap.Error(token.Error()))

	} else {
		t := token.(*MQTT.SubscribeToken)
		fmt.Println(t.Result())
	}

}

func Test2(t *testing.T) {
	mqttServer := &service.Server{
		KeepAlive:        300,           // seconds
		ConnectTimeout:   2,             // seconds
		SessionsProvider: "mem",         // keeps sessions in memory
		Authenticator:    "mockSuccess", // always succeed
		TopicsProvider:   "mem",         // keeps topic subscriptions in memory
		AclProvider:      acl.TopicNumAuthType,
		GetAuthFunc: func(userName, topic string) interface{} {
			return 1
		},
	}

	if err := mqttServer.ListenAndServe("tcp://127.0.0.1:8080"); err != nil {
		fmt.Println("mqtt error", zap.Error(err))
	}
}

func TestTggw(t *testing.T) {
	mqttServer := &service.Server{
		KeepAlive:        300,           // seconds
		ConnectTimeout:   2,             // seconds
		SessionsProvider: "mem",         // keeps sessions in memory
		Authenticator:    "mockSuccess", // always succeed
		TopicsProvider:   "mem",         // keeps topic subscriptions in memory
		AclProvider:      acl.TopicSetAuthType,
		GetAuthFunc: func(userName, topic string) interface{} {

			return userName == topic
		},
	}

	if err := mqttServer.ListenAndServe("tcp://127.0.0.1:8080"); err != nil {
		fmt.Println("mqtt error", zap.Error(err))
	}
}

func Test3(t *testing.T) {
	mqttServer := &service.Server{
		KeepAlive:        300,           // seconds
		ConnectTimeout:   2,             // seconds
		SessionsProvider: "mem",         // keeps sessions in memory
		Authenticator:    "mockSuccess", // always succeed
		TopicsProvider:   "mem",         // keeps topic subscriptions in memory
		AclProvider:      acl.TopicSetAuthType,
		GetAuthFunc: func(userName, topic string) interface{} {
			return true
		},
	}

	if err := mqttServer.ListenAndServe("tcp://127.0.0.1:8081"); err != nil {
		fmt.Println("mqtt error", zap.Error(err))
	}
}

func Test4(t *testing.T) {
	var f MQTT.MessageHandler = func(MQTT.Client, MQTT.Message) {
		fmt.Println("rece")
	}

	connOpts := &MQTT.ClientOptions{
		ClientID:             "ds-live",
		CleanSession:         true,
		Username:             "fwd",
		Password:             "zzz",
		MaxReconnectInterval: 1 * time.Second,
		KeepAlive:            30 * time.Second,
		AutoReconnect:        true,
		PingTimeout:          10 * time.Second,
		ConnectTimeout:       30 * time.Second,
		TLSConfig:            tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert},
		OnConnectionLost:     func(c MQTT.Client, err error) { fmt.Println("mqtt disconnected.", zap.Error(err)) },
	}
	connOpts.AddBroker("tcp://127.0.0.1:8081")

	mc := MQTT.NewClient(connOpts)
	if token := mc.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("mqtt connection failed.", zap.Error(token.Error()))
		return
	}
	defer mc.Disconnect(12)
	if token := mc.Subscribe("ds", 0, f); token.Wait() && token.Error() != nil {
		fmt.Println("publish failed.", zap.Error(token.Error()))

	} else {
		t := token.(*MQTT.SubscribeToken)
		fmt.Println(t.Result())
	}

}
