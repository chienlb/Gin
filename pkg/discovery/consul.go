package discovery

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/consul/api"
)

// ConsulClient wraps Consul service discovery
type ConsulClient struct {
	client      *api.Client
	serviceID   string
	serviceName string
	servicePort int
	tags        []string
}

// ConsulConfig contains Consul configuration
type ConsulConfig struct {
	Address     string
	ServiceID   string
	ServiceName string
	ServicePort int
	Tags        []string
	HealthCheck HealthCheck
}

// HealthCheck defines health check configuration
type HealthCheck struct {
	HTTP                           string
	Interval                       string
	Timeout                        string
	DeregisterCriticalServiceAfter string
}

// NewConsulClient creates a new Consul client
func NewConsulClient(cfg ConsulConfig) (*ConsulClient, error) {
	config := api.DefaultConfig()
	config.Address = cfg.Address

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	return &ConsulClient{
		client:      client,
		serviceID:   cfg.ServiceID,
		serviceName: cfg.ServiceName,
		servicePort: cfg.ServicePort,
		tags:        cfg.Tags,
	}, nil
}

// Register registers the service with Consul
func (c *ConsulClient) Register(healthCheckURL string) error {
	registration := &api.AgentServiceRegistration{
		ID:      c.serviceID,
		Name:    c.serviceName,
		Port:    c.servicePort,
		Tags:    c.tags,
		Address: getLocalIP(),
		Check: &api.AgentServiceCheck{
			HTTP:                           healthCheckURL,
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	log.Printf("Service registered with Consul: %s (ID: %s)", c.serviceName, c.serviceID)
	return nil
}

// Deregister deregisters the service from Consul
func (c *ConsulClient) Deregister() error {
	if err := c.client.Agent().ServiceDeregister(c.serviceID); err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	log.Printf("Service deregistered from Consul: %s (ID: %s)", c.serviceName, c.serviceID)
	return nil
}

// DiscoverService discovers services by name
func (c *ConsulClient) DiscoverService(serviceName string) ([]*api.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service: %w", err)
	}

	return services, nil
}

// GetServiceAddress gets a service address (load balanced)
func (c *ConsulClient) GetServiceAddress(serviceName string) (string, error) {
	services, err := c.DiscoverService(serviceName)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances found for service: %s", serviceName)
	}

	// Simple round-robin (in production, use more sophisticated load balancing)
	service := services[time.Now().Unix()%int64(len(services))]
	address := fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)

	return address, nil
}

// WatchService watches for service changes
func (c *ConsulClient) WatchService(serviceName string, callback func([]*api.ServiceEntry)) {
	var lastIndex uint64

	for {
		services, meta, err := c.client.Health().Service(serviceName, "", true, &api.QueryOptions{
			WaitIndex: lastIndex,
			WaitTime:  5 * time.Minute,
		})

		if err != nil {
			log.Printf("Error watching service %s: %v", serviceName, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if meta.LastIndex != lastIndex {
			lastIndex = meta.LastIndex
			callback(services)
		}
	}
}

// GetAllServices gets all registered services
func (c *ConsulClient) GetAllServices() (map[string][]string, error) {
	services, _, err := c.client.Catalog().Services(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get all services: %w", err)
	}

	return services, nil
}

// GetServiceHealth gets service health status
func (c *ConsulClient) GetServiceHealth(serviceName string) (string, error) {
	checks, _, err := c.client.Health().Checks(serviceName, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get service health: %w", err)
	}

	if len(checks) == 0 {
		return "unknown", nil
	}

	// Return the worst status
	worstStatus := "passing"
	for _, check := range checks {
		if check.Status == "critical" {
			return "critical", nil
		}
		if check.Status == "warning" {
			worstStatus = "warning"
		}
	}

	return worstStatus, nil
}

// SetKV sets a key-value pair in Consul KV store
func (c *ConsulClient) SetKV(key string, value []byte) error {
	p := &api.KVPair{Key: key, Value: value}
	_, err := c.client.KV().Put(p, nil)
	if err != nil {
		return fmt.Errorf("failed to set KV: %w", err)
	}
	return nil
}

// GetKV gets a value from Consul KV store
func (c *ConsulClient) GetKV(key string) ([]byte, error) {
	pair, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get KV: %w", err)
	}
	if pair == nil {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return pair.Value, nil
}

// DeleteKV deletes a key from Consul KV store
func (c *ConsulClient) DeleteKV(key string) error {
	_, err := c.client.KV().Delete(key, nil)
	if err != nil {
		return fmt.Errorf("failed to delete KV: %w", err)
	}
	return nil
}

// getLocalIP gets the local IP address
func getLocalIP() string {
	// In production, implement proper IP detection
	return "127.0.0.1"
}
