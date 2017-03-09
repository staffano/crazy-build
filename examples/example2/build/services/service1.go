package services

import (
	"log"

	"github.com/staffano/crazy-build/artifact"
)

// A Service1API ..
type Service1API interface {
	artifact.ServiceAPI
	PrintHello()
}

// Service1Impl implements the Service1API.
type Service1Impl struct {
}

// Allocate this service for an artifact. A token is received
// to handle the service
func (s *Service1Impl) Allocate() int {
	return 0
}

// Deallocate the artifact from the service
func (s *Service1Impl) Deallocate(token int) {}

// IsAvailable for allocation?
func (s *Service1Impl) IsAvailable() bool {

	return true
}

// PrintHello prints Print
func (s *Service1Impl) PrintHello() {
	log.Print("Hello")
}

// Satisfies checks if this instance satisfies the requiresments
func (s *Service1Impl) Satisfies(req string) bool {

	return true
}

func init() {
	artifact.RegisterServiceInstance(&Service1Impl{})
}
