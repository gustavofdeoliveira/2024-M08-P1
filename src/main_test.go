package main

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"testing"
)

func TestOpenFileSuccess(t *testing.T) {
	fmt.Println("TestOpenFileSuccess")
	// Setup: Cria um arquivo temporário para teste.
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Erro ao criar arquivo temporário: %s", err)
	}
	tmpfilePath := tmpfile.Name()

	// Cleanup: Garante que o arquivo temporário seja removido após o teste.
	defer os.Remove(tmpfilePath)
	tmpfile.Close()

	// Teste: Tenta abrir o arquivo temporário.
	file := openFile(tmpfilePath)
	if file == nil {
		t.Errorf("openFile retornou nil para um arquivo existente")
	}

	// Não esqueça de fechar o arquivo aberto pela função openFile.
	file.Close()
}

func TestReadFileSuccess(t *testing.T) {
	fmt.Println("TestReadFileSuccess")
	// Setup: Cria um arquivo temporário para teste.
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Erro ao criar arquivo temporário: %s", err)
	}
	tmpfilePath := tmpfile.Name()

	// Cleanup: Garante que o arquivo temporário seja removido após o teste.
	defer os.Remove(tmpfilePath)
	tmpfile.Close()

	// Teste: Tenta abrir o arquivo temporário.
	file := openFile(tmpfilePath)
	if file == nil {
		t.Errorf("openFile retornou nil para um arquivo existente")
	}
	bytes := readFile(file)

	if bytes == nil {
		t.Errorf("readFile retornou nil para um arquivo existente")
	}
	file.Close()

}

func TestCreateAndPublisObject(t *testing.T) {
	fmt.Println("TestCreateAndPublisObject")
	var file = readFile(openFile("data.json"))
	var result []map[string]interface{}
	json.Unmarshal(file, &result)
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetClientID("go_subscriber")

	subscriber := NewMQTTSubscriber()
	for _, obj := range result {
		publishObject(obj, subscriber)
		if !subscriber.client.IsConnected() {
			t.Errorf("Erro de conexão")
		}
	}
	subscriber.client.Disconnect(250)
}

func TestPublicAndRecevedMessageQos(t *testing.T) {
	fmt.Println("TestPublicAndRecevedMessageQos")
	var file = openFile("data.json")
	var bytes = readFile(file)

	var result []map[string]interface{}
	var err = json.Unmarshal(bytes, &result)
	if err != nil {
		t.Fatalf("Erro ao decodificar o JSON: %s", err)
	}

	var subscriber = NewMQTTSubscriber()
	var messageReceiver = &MQTTSubscriber{}

	messageChannel := make(chan string)
	qosChannel := make(chan string)

	subscriber.client.Subscribe("topic/publisher", 1, func(client MQTT.Client, msg MQTT.Message){
		messageReceiver.ReceiveMessage(client, msg)
		messageChannel <- string(msg.Payload())
		qosChannel <- string(msg.Qos())
	})

	for _, obj := range result {
		publishObject(obj, subscriber)
	}
	
	receivedMessage := <-messageChannel
	receivedQos := <-qosChannel

	if receivedMessage == "" {
		t.Errorf("Erro ao receber mensagem")
	}
	if receivedQos == "" {
		t.Errorf("Erro ao receber QoS")
	}


	close(messageChannel)
	close(qosChannel)
}

