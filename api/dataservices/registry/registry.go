package registry

import (
	"fmt"

	portainer "github.com/portainer/portainer/api"

	"github.com/rs/zerolog/log"
)

// BucketName represents the name of the bucket where this service stores data.
const BucketName = "registries"

// Service represents a service for managing environment(endpoint) data.
type Service struct {
	connection portainer.Connection
}

func (service *Service) BucketName() string {
	return BucketName
}

// NewService creates a new instance of a service.
func NewService(connection portainer.Connection) (*Service, error) {
	err := connection.SetServiceName(BucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		connection: connection,
	}, nil
}

func (service *Service) Tx(tx portainer.Transaction) ServiceTx {
	return ServiceTx{
		service: service,
		tx:      tx,
	}
}

// Registry returns a registry by ID.
func (service *Service) Registry(ID portainer.RegistryID) (*portainer.Registry, error) {
	var registry portainer.Registry
	identifier := service.connection.ConvertToKey(int(ID))

	err := service.connection.GetObject(BucketName, identifier, &registry)
	if err != nil {
		return nil, err
	}

	return &registry, nil
}

// Registries returns an array containing all the registries.
func (service *Service) Registries() ([]portainer.Registry, error) {
	var registries = make([]portainer.Registry, 0)

	err := service.connection.GetAll(
		BucketName,
		&portainer.Registry{},
		func(obj interface{}) (interface{}, error) {
			registry, ok := obj.(*portainer.Registry)
			if !ok {
				log.Debug().Str("obj", fmt.Sprintf("%#v", obj)).Msg("failed to convert to Registry object")
				return nil, fmt.Errorf("Failed to convert to Registry object: %s", obj)
			}

			registries = append(registries, *registry)

			return &portainer.Registry{}, nil
		})

	return registries, err
}

// Create creates a new registry.
func (service *Service) Create(registry *portainer.Registry) error {
	return service.connection.CreateObject(
		BucketName,
		func(id uint64) (int, interface{}) {
			registry.ID = portainer.RegistryID(id)
			return int(registry.ID), registry
		},
	)
}

// UpdateRegistry updates a registry.
func (service *Service) UpdateRegistry(ID portainer.RegistryID, registry *portainer.Registry) error {
	identifier := service.connection.ConvertToKey(int(ID))
	return service.connection.UpdateObject(BucketName, identifier, registry)
}

// DeleteRegistry deletes a registry.
func (service *Service) DeleteRegistry(ID portainer.RegistryID) error {
	identifier := service.connection.ConvertToKey(int(ID))
	return service.connection.DeleteObject(BucketName, identifier)
}
