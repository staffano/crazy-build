package artifact

// ServiceAPI is the API all services have to comply to
// in order to be handled as services
type ServiceAPI interface {

	// Allocate this service for an artifact. A token is received
	// to handle the service.
	Allocate() int

	// Deallocate the artifact from the service
	Deallocate(token int)

	// IsAvailable for allocation?
	IsAvailable() bool

	// Does the service Satisfies the requiremen?
	Satisfies(requirement string) bool
}

// services are all registered services and they all implement
// a specific API interface that is used to find them
var services []*ServiceAPI

// RegisterServiceInstance registers a service instance in the service database
func RegisterServiceInstance(si ServiceAPI) {
	services = append(services, &si)
}
