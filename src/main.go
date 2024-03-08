package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MQTTSubscriber é uma estrutura que representa um assinante MQTT.
type MQTTSubscriber struct {
    client MQTT.Client
}

// MessageReceiver é uma interface que define um método para receber mensagens MQTT.
type MessageReceiver interface {
    ReceiveMessage(client MQTT.Client, msg MQTT.Message)
}

func openFile(path string) *os.File {
    file, err := os.Open(path)
    if err != nil {
        log.Fatalf("Erro ao abrir o arquivo: %s", err)
    }
    return file
}

func readFile(file *os.File) []byte {
    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatalf("Erro ao ler o arquivo: %s", err)
    }
    return bytes
}


func publishObject(newObject map[string]interface{}, singletonClient *MQTTSubscriber) string {
    jsonData, err := json.Marshal(newObject)
    if err != nil {
        fmt.Println("Error marshalling JSON:", err)
        return ""
    }
    token := singletonClient.client.Publish("topic/publisher", 0, false, jsonData)
    token.Wait()
    fmt.Println("Publicado:", string(jsonData))
    return string(jsonData)
}

var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
    fmt.Println("Connected")
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
    fmt.Printf("Connection lost: %v", err)
}

// NewMQTTSubscriber cria e retorna um novo assinante MQTT.
func NewMQTTSubscriber() *MQTTSubscriber {
    opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
    opts.OnConnect = connectHandler
    opts.OnConnectionLost = connectLostHandler

    opts.SetClientID("go_subscriber")

    client := MQTT.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatalf("Error connecting to MQTT broker: %s", token.Error())
    }

    return &MQTTSubscriber{client: client}
}

func (s *MQTTSubscriber) ReceiveMessage(client MQTT.Client, msg MQTT.Message) {
    var result map[string]interface{}
    json.Unmarshal([]byte(msg.Payload()), &result)

    var id = result["id"].(string)
    var temperatura = result["temperatura"].(float64)
    if result["tipo"] =="freezer"{

        if temperatura >= -15 {
            fmt.Printf("Lj " + id[3:4]+ ": Freezer "+ id[6:7] + "  | %f°C | Alerta: [ALERTA: Temperatura ALTA]\n", temperatura)
        }
        if temperatura <= -25 {
            fmt.Printf("Lj " + id[3:4]+ ": Freezer "+ id[6:7] + "  | %f°C | Alerta: [ALERTA: Temperatura BAIXA ]\n", temperatura)
        }
        if temperatura >= -25 && temperatura <= -15 {
            fmt.Printf("Lj " + id[3:4]+ ": Freezer "+ id[6:7] + "  | %f°C\n", temperatura)
        }
    }

    if result["tipo"] =="geladeira"{
        if temperatura >= 10 {
            fmt.Printf("Lj " + id[3:4]+ ": Geladeira "+ id[6:7] + "  | %f°C | Alerta: [ALERTA: Temperatura ALTA]\n", temperatura)
        }
        if temperatura <= 2 {
            fmt.Printf("Lj " + id[3:4]+ ": Geladeira "+ id[6:7] + "  | %f°C | Alerta: [ALERTA: Temperatura BAIXA]\n", temperatura)
        
        }
        if temperatura >= 2 && temperatura <= 10 {
            fmt.Printf("Lj " + id[3:4]+ ": Geladeira "+ id[6:7] + "  | %f°C\n", temperatura)
        }
    }
}

func main() {
    subscriber := NewMQTTSubscriber()

    subscriber.client.Subscribe("topic/publisher", 1, func(client MQTT.Client, msg MQTT.Message) {
        subscriber.ReceiveMessage(client, msg)
    })

	var file = readFile(openFile("data.json"))

	result := []map[string]interface{}{}
	var err = json.Unmarshal(file, &result)
	if err != nil {
		log.Fatalf("Erro ao decodificar o JSON: %s", err)
	}
    for _, item := range result {
        publishObject(item, subscriber)
        time.Sleep(1 * time.Second)
    }
	
	
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
    <-sigCh
    fmt.Println("Encerrando o programa.")
    subscriber.client.Disconnect(250)
}


