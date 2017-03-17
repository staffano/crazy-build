package services

import (
	"log"

	"github.com/staffano/crazy-build/artifact"
)

// A Service2API ..
type Service2API interface {
	artifact.ServiceAPI
	PrintWorld()
}

// Service2Impl implements the Service2API.
type Service2Impl struct {
}

// Allocate this service for an artifact. A token is received
// to handle the service
func (s *Service2Impl) Allocate() int {
	return 0
}

// Deallocate the artifact from the service
func (s *Service2Impl) Deallocate(token int) {}

// IsAvailable for allocation?
func (s *Service2Impl) IsAvailable() bool {
	return true
}

// PrintWorld prints Print
func (s *Service2Impl) PrintWorld() {
	log.Print("World")
}

// Satisfies checks if this instance satisfies the requiresments
func (s *Service2Impl) Satisfies(req string) bool {

	return true
}

func init() {
	artifact.RegisterServiceInstance(&Service2Impl{})
}
