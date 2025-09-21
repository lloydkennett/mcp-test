package registry

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Service interface {
	Name() string
	Enabled() bool
	Register(*mcp.Server)
}

type Registry struct {
	services []Service
}

func (r *Registry) Add(s Service) {
	if s == nil {
		return
	}
	r.services = append(r.services, s)
}

func (r *Registry) RegisterAll(server *mcp.Server) {
	for _, s := range r.services {
		if s.Enabled() {
			s.Register(server)
		}
	}
}

func (r *Registry) EnabledServices() []string {
	var out []string
	for _, s := range r.services {
		if s.Enabled() {
			out = append(out, s.Name())
		}
	}
	return out
}

func (r *Registry) DisabledServices() []string {
	var out []string
	for _, s := range r.services {
		if !s.Enabled() {
			out = append(out, s.Name())
		}
	}
	return out
}
