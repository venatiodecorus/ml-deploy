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

var (
	oClient *azopenai.Client
	clientOnce sync.Once
	clientInitErr error
)

func GetClient() (*azopenai.Client, error) {
	clientOnce.Do(func() {
		oClient, clientInitErr = initializeClient()
	})

	if clientInitErr != nil {
		return nil, clientInitErr
	}
	return oClient, nil
}

func initializeClient() (*azopenai.Client, error) {
	keyCred := azcore.NewKeyCredential(os.Getenv("OPENAI_KEY"))
	client,err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCred, nil)

	if err != nil {
		log.Printf("failed to create client: %s", err)
		return nil, err
	}

	return client, nil
}

func makeRequest(messages []azopenai.ChatRequestMessageClassification) (string, error) {
	client, err := GetClient()
	if err != nil {
		log.Printf("failed to get client: %s", err)
		return "",err
	}

	modelID := "gpt-4o"
	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages: messages,
		DeploymentName: &modelID,
	}, nil)

	if err != nil {
		log.Printf("failed to get chat completions: %s", err)
		return "",err
	}

	fmt.Printf("Response: %s", *resp.Choices[0].Message.Content)
	return *resp.Choices[0].Message.Content, nil
}

func dockerRequest(instructions string) (string, error) {
	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{Content: to.Ptr("Your sole task is to generate a Dockerfile based on the instructions provided. You should not return any other formats or unrelated information. You should not ask for clarification. You should not provide additional information")},
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(instructions)},
	}

	return makeRequest(messages)
}

func terraformRequest(instructions string) (string, error) {
	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{Content: to.Ptr("Your only function is to generate a Terraform file based on the instructions provided. You should not return anything other than a valid file, or ask for any additional information. The file should use the hetznercloud/hcloud provider, and always include the required_providers definition in the file. The file should accept the Hetzner API token from an environment variable named HCLOUD_TOKEN. The file should be returned as raw text, with no JSON wrappers. Don't include an SSH key.")},
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(instructions)},
	}

	return makeRequest(messages)
}