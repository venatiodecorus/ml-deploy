package hetzner

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var client *hcloud.Client
var clientOnce sync.Once
// var clientInitErr error

func getClient() *hcloud.Client {
	clientOnce.Do(func() {
		client = hcloud.NewClient(hcloud.WithToken(os.Getenv("TF_VAR_hcloud_token")))
	})

	return client
}

func ServerList() []*hcloud.Server {
	client := getClient()
	servers, _, err := client.Server.List(context.Background(), hcloud.ServerListOpts{})
	if err != nil {
		log.Printf("failed to list servers: %s", err)
	}

	return servers
}