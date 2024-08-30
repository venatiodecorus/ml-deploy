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

func VolumeList() []*hcloud.Volume {
	client := getClient()
	volumes, _, err := client.Volume.List(context.Background(), hcloud.VolumeListOpts{})
	if err != nil {
		log.Printf("failed to list volumes: %s", err)
	}

	return volumes
}

func LoadBalancerList() []*hcloud.LoadBalancer {
	client := getClient()
	loadBalancers, _, err := client.LoadBalancer.List(context.Background(), hcloud.LoadBalancerListOpts{})
	if err != nil {
		log.Printf("failed to list load balancers: %s", err)
	}

	return loadBalancers
}

func FloatingIPList() []*hcloud.FloatingIP {
	client := getClient()
	floatingIPs, _, err := client.FloatingIP.List(context.Background(), hcloud.FloatingIPListOpts{})
	if err != nil {
		log.Printf("failed to list floating IPs: %s", err)
	}

	return floatingIPs
}

func NetworkList() []*hcloud.Network {
	client := getClient()
	networks, _, err := client.Network.List(context.Background(), hcloud.NetworkListOpts{})
	if err != nil {
		log.Printf("failed to list networks: %s", err)
	}

	return networks
}

func FirewallList() []*hcloud.Firewall {
	client := getClient()
	firewalls, _, err := client.Firewall.List(context.Background(), hcloud.FirewallListOpts{})
	if err != nil {
		log.Printf("failed to list firewalls: %s", err)
	}

	return firewalls
}

type AllResources struct {
    Servers      []*hcloud.Server
    Volumes      []*hcloud.Volume
    LoadBalancers []*hcloud.LoadBalancer
    FloatingIPs  []*hcloud.FloatingIP
    Networks     []*hcloud.Network
    Firewalls    []*hcloud.Firewall
}

// No list all function in hcloud lib, need to gather each type of resource
func ListAll() AllResources {
	// client := getClient()

	servers := ServerList()
	volumes := VolumeList()
	loadBalancers := LoadBalancerList()
	floatingIPs := FloatingIPList()
	networks := NetworkList()
	firewalls := FirewallList()

	// Return a collection including all these resources
	return AllResources{
        Servers:      servers,
        Volumes:      volumes,
        LoadBalancers: loadBalancers,
        FloatingIPs:  floatingIPs,
        Networks:     networks,
        Firewalls:  firewalls,
	}
}