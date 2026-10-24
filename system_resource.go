package routeros

import (
	"context"
	"fmt"
)

// Info contains commonly used RouterOS system information.
type Info struct {
	Identity    string
	Version     string
	BuildTime   string
	FreeMemory  string
	TotalMemory string
	CPU         string
	CPUCount    string
	Uptime      string
	ArchName    string
	BoardName   string
	Model       string
	Serial      string
}

// SystemResource implements the system operations used by the root API wrapper.
type SystemResource struct {
	*BaseResource
	client *Client
}

func NewSystemResource(client *Client) *SystemResource {
	return &SystemResource{
		BaseResource: NewBaseResource(client, "/system/resource"),
		client: client,
	}
}

func (s *SystemResource) IdentityContext(ctx context.Context) (string, error) {
	responses, err := s.client.SendContext(ctx, "/system/identity/print")
	if err != nil {
		return "", err
	}
	return NewResponseParser(responses).Value("name"), nil
}

func (s *SystemResource) Identity() (string, error) {
	return s.IdentityContext(context.Background())
}

func (s *SystemResource) ResourceContext(ctx context.Context) (*Resource, error) {
	list, err := s.ListContext(ctx)
	if err != nil {
		return nil, err
	}
	if list.Len() == 0 {
		return nil, ErrResourceNotFound
	}
	return list.First(), nil
}

func (s *SystemResource) Resource() (*Resource, error) {
	return s.ResourceContext(context.Background())
}

func (s *SystemResource) GetInfoContext(ctx context.Context) (*Info, error) {
	identity, err := s.IdentityContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("read system identity: %w", err)
	}
	resource, err := s.ResourceContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("read system resource information: %w", err)
	}
	return &Info{
		Identity: identity,
		Version: resource.Get("version"),
		BuildTime: resource.Get("build-time"),
		FreeMemory: resource.Get("free-memory"),
		TotalMemory: resource.Get("total-memory"),
		CPU: resource.Get("cpu"),
		CPUCount: resource.Get("cpu-count"),
		Uptime: resource.Get("uptime"),
		ArchName: resource.Get("architecture-name"),
		BoardName: resource.Get("board-name"),
		Model: resource.Get("model"),
		Serial: resource.Get("serial-number"),
	}, nil
}

func (s *SystemResource) GetInfo() (*Info, error) {
	return s.GetInfoContext(context.Background())
}

func (s *SystemResource) PackagesContext(ctx context.Context) (*ResourceList, error) {
	responses, err := s.client.SendContext(ctx, "/system/package/print")
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

func (s *SystemResource) Packages() (*ResourceList, error) {
	return s.PackagesContext(context.Background())
}
