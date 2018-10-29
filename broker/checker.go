package broker

import (
	"errors"
	"fmt"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
	uuid "github.com/satori/go.uuid"
)

type Broker struct {
	client osb.Client
}

func NewBroker(url, username, password string) (*Broker, error) {
	config := osb.DefaultClientConfiguration()
	config.URL = url
	config.AuthConfig = &osb.AuthConfig{
		BasicAuthConfig: &osb.BasicAuthConfig{
			Username: username,
			Password: password,
		},
	}

	client, err := osb.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Broker{client}, nil
}

func (b *Broker) Check() error {
	err := b.ping()
	if err != nil {
		return err
	}

	err = b.checkLifecycle()
	if err != nil {
		return err
	}

	return nil
}

func (b *Broker) ping() error {
	_, err := b.client.GetCatalog()
	return err
}

func (b *Broker) checkLifecycle() error {
	fmt.Println("Getting catalog")
	resp, err := b.client.GetCatalog()
	if err != nil {
		return err
	}

	if len(resp.Services) == 0 {
		return errors.New("No services to provision")
	}
	service := resp.Services[0].ID

	if len(resp.Services[0].Plans) == 0 {
		return errors.New("No plans of service to provision")
	}

	plan := resp.Services[0].Plans[0].ID

	guid := uuid.NewV4().String()
	provReq := osb.ProvisionRequest{
		InstanceID:        guid,
		AcceptsIncomplete: false,
		ServiceID:         service,
		PlanID:            plan,
		OrganizationGUID:  "blah",
		SpaceGUID:         "blah",
	}

	fmt.Println("Creating instance")
	_, err = b.client.ProvisionInstance(&provReq)
	if err != nil {
		return err
	}

	deprovReq := osb.DeprovisionRequest{
		InstanceID:        guid,
		AcceptsIncomplete: false,
		ServiceID:         service,
		PlanID:            plan,
	}

	fmt.Println("Deleting instance")

	_, err = b.client.DeprovisionInstance(&deprovReq)
	if err != nil {
		return err
	}
	return nil
}

//idea declarative manifest for this stuff
