package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

// type Client struct {
var (
	client *azopenai.Client
	clientOnce sync.Once
	clientInitErr error
)
// }

// func NewClient() (bool) {
//     // keyCred := azcore.NewKeyCredential(os.Getenv("OPENAI_KEY"))
// 	// client,err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCred, nil)

// 	// if err != nil {
// 	// 	log.Fatalf("failed to create client: %s", err)
//     //     return nil, err
//     // }


// 	clientOnce.Do(func() {
// 		keyCred := azcore.NewKeyCredential(os.Getenv("OPENAI_KEY"))
// 		var err error
// 		Client,err = azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCred, nil)

// 		if err != nil {
// 			log.Fatalf("failed to create client: %s", err)
// 			return false //nil, err
// 		}
// 	})

// 	return true
// 	// return &Client{Client: client}, nil
// }

func GetClient() (*azopenai.Client, error) {
	clientOnce.Do(func() {
		client, clientInitErr = initializeClient()
	})

	if clientInitErr != nil {
		return nil, clientInitErr
	}
	return client, nil
}

func initializeClient() (*azopenai.Client, error) {
	keyCred := azcore.NewKeyCredential(os.Getenv("OPENAI_KEY"))
	client,err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCred, nil)

	if err != nil {
		log.Fatalf("failed to create client: %s", err)
		return nil, err
	}

	return client, nil
}

func makeRequest(messages []azopenai.ChatRequestMessageClassification) string {
	client, err := GetClient()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
		return ""
	}

	// resp, err := client.GetChatCompletions(context.TODO(), &azopenai.GetChatCompletionsOptions{
	// 	Model: "gpt-4o",
	// 	Prompt: instructions,
	// })
	// messages := []azopenai.ChatRequestMessageClassification{
	// 	&azopenai.ChatRequestSystemMessage{Content: to.Ptr("You are a helpful assistant.")},
	// 	&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(instructions)},
	// }

	modelID := "gpt-4o"
	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages: messages,
		DeploymentName: &modelID,
	}, nil)

	if err != nil {
		log.Fatalf("failed to get chat completions: %s", err)
		return ""
	}

	fmt.Fprintf(os.Stdout, "Response: %s", *resp.Choices[0].Message.Content)
	return *resp.Choices[0].Message.Content
}

func dockerRequest(instructions string) string {
	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{Content: to.Ptr("Your sole task is to generate a Dockerfile based on the instructions provided. You should not return any other formats or unrelated information. You should not ask for clarification. You should not provide additional information")},
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(instructions)},
	}

	return makeRequest(messages)
}

// func GenerateDockerfile(instructions string) {
// 	client, err := GetClient()
// 	if err != nil {
// 		log.Fatalf("failed to get client: %s", err)
// 	}


// }