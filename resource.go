package routeros

import (
	"context"
)

// ResourceInterface provides a generic interface for RouterOS resources
type ResourceInterface interface {
	// List retrieves all resources
	List() (*ResourceList, error)

	// ListContext retrieves all resources with context
	ListContext(ctx context.Context) (*ResourceList, error)

	// Get retrieves a specific resource by ID
	Get(id string) (*Resource, error)

	// GetContext retrieves a specific resource by ID with context
	GetContext(ctx context.Context, id string) (*Resource, error)

	// Add creates a new resource
	Add(attrs map[string]string) (string, error)

	// AddContext creates a new resource with context
	AddContext(ctx context.Context, attrs map[string]string) (string, error)

	// Update modifies an existing resource
	Update(id string, attrs map[string]string) error

	// UpdateContext modifies an existing resource with context
	UpdateContext(ctx context.Context, id string, attrs map[string]string) error

	// Delete removes a resource
	Delete(id string) error

	// DeleteContext removes a resource with context
	DeleteContext(ctx context.Context, id string) error

	// Enable enables a resource
	Enable(id string) error

	// EnableContext enables a resource with context
	EnableContext(ctx context.Context, id string) error

	// Disable disables a resource
	Disable(id string) error

	// DisableContext disables a resource with context
	DisableContext(ctx context.Context, id string) error
}

// BaseResource provides a base implementation for resource operations
type BaseResource struct {
	client *Client
	path   string
}

// NewBaseResource creates a new BaseResource
func NewBaseResource(client *Client, path string) *BaseResource {
	return &BaseResource{
		client: client,
		path:   path,
	}
}

// List retrieves all resources
func (r *BaseResource) List() (*ResourceList, error) {
	return r.ListContext(context.Background())
}

// ListContext retrieves all resources with context
func (r *BaseResource) ListContext(ctx context.Context) (*ResourceList, error) {
	responses, err := r.client.SendContext(ctx, r.path+"/print")
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// Get retrieves a specific resource by ID
func (r *BaseResource) Get(id string) (*Resource, error) {
	return r.GetContext(context.Background(), id)
}

// GetContext retrieves a specific resource by ID with context
func (r *BaseResource) GetContext(ctx context.Context, id string) (*Resource, error) {
	responses, err := r.client.SendContext(ctx, r.path+"/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	list := NewResourceList(responses)
	if list.Len() == 0 {
		return nil, ErrResourceNotFound
	}
	return list.First(), nil
}

// Add creates a new resource
func (r *BaseResource) Add(attrs map[string]string) (string, error) {
	return r.AddContext(context.Background(), attrs)
}

// AddContext creates a new resource with context
func (r *BaseResource) AddContext(ctx context.Context, attrs map[string]string) (string, error) {
	args := make([]string, 0)
	for k, v := range attrs {
		args = append(args, "="+k+"="+v)
	}

	responses, err := r.client.SendContext(ctx, r.path+"/add", args...)
	if err != nil {
		return "", err
	}

	// Extract the new ID from the response
	for _, resp := range responses {
		if resp.Type == ResponseDone {
			// Look for the ID in different possible formats
			if id, ok := resp.Data["ret"]; ok {
				return id, nil
			}
		}
	}

	return "", nil
}

// Update modifies an existing resource
func (r *BaseResource) Update(id string, attrs map[string]string) error {
	return r.UpdateContext(context.Background(), id, attrs)
}

// UpdateContext modifies an existing resource with context
func (r *BaseResource) UpdateContext(ctx context.Context, id string, attrs map[string]string) error {
	args := []string{".id=" + id}
	for k, v := range attrs {
		args = append(args, "="+k+"="+v)
	}

	_, err := r.client.SendContext(ctx, r.path+"/set", args...)
	return err
}

// Delete removes a resource
func (r *BaseResource) Delete(id string) error {
	return r.DeleteContext(context.Background(), id)
}

// DeleteContext removes a resource with context
func (r *BaseResource) DeleteContext(ctx context.Context, id string) error {
	_, err := r.client.SendContext(ctx, r.path+"/remove", ".id="+id)
	return err
}

// Enable enables a resource
func (r *BaseResource) Enable(id string) error {
	return r.EnableContext(context.Background(), id)
}

// EnableContext enables a resource with context
func (r *BaseResource) EnableContext(ctx context.Context, id string) error {
	_, err := r.client.SendContext(ctx, r.path+"/enable", ".id="+id)
	return err
}

// Disable disables a resource
func (r *BaseResource) Disable(id string) error {
	return r.DisableContext(context.Background(), id)
}

// DisableContext disables a resource with context
func (r *BaseResource) DisableContext(ctx context.Context, id string) error {
	_, err := r.client.SendContext(ctx, r.path+"/disable", ".id="+id)
	return err
}

// Find finds resources by attributes
func (r *BaseResource) Find(attrs map[string]string) (*ResourceList, error) {
	return r.FindContext(context.Background(), attrs)
}

// FindContext finds resources by attributes with context
func (r *BaseResource) FindContext(ctx context.Context, attrs map[string]string) (*ResourceList, error) {
	args := make([]string, 0)
	for k, v := range attrs {
		args = append(args, "?="+k+"="+v)
	}

	responses, err := r.client.SendContext(ctx, r.path+"/print", args...)
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// FindOne finds a single resource by attributes
func (r *BaseResource) FindOne(attrs map[string]string) (*Resource, error) {
	return r.FindOneContext(context.Background(), attrs)
}

// FindOneContext finds a single resource by attributes with context
func (r *BaseResource) FindOneContext(ctx context.Context, attrs map[string]string) (*Resource, error) {
	list, err := r.FindContext(ctx, attrs)
	if err != nil {
		return nil, err
	}
	if list.Len() == 0 {
		return nil, ErrResourceNotFound
	}
	return list.First(), nil
}

// FindByName finds a resource by name
func (r *BaseResource) FindByName(name string) (*Resource, error) {
	return r.FindOneContext(context.Background(), map[string]string{"name": name})
}

// FindByNameContext finds a resource by name with context
func (r *BaseResource) FindByNameContext(ctx context.Context, name string) (*Resource, error) {
	return r.FindOneContext(ctx, map[string]string{"name": name})
}

// Print runs the print command with optional arguments
func (r *BaseResource) Print(args ...string) (*ResourceList, error) {
	return r.PrintContext(context.Background(), args...)
}

// PrintContext runs the print command with context
func (r *BaseResource) PrintContext(ctx context.Context, args ...string) (*ResourceList, error) {
	responses, err := r.client.SendContext(ctx, r.path+"/print", args...)
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// Count returns the number of resources
func (r *BaseResource) Count() (int, error) {
	return r.CountContext(context.Background())
}

// CountContext returns the number of resources with context
func (r *BaseResource) CountContext(ctx context.Context) (int, error) {
	list, err := r.ListContext(ctx)
	if err != nil {
		return 0, err
	}
	return list.Len(), nil
}

// Comment sets a comment on a resource
func (r *BaseResource) Comment(id, comment string) error {
	return r.CommentContext(context.Background(), id, comment)
}

// CommentContext sets a comment on a resource with context
func (r *BaseResource) CommentContext(ctx context.Context, id, comment string) error {
	_, err := r.client.SendContext(ctx, r.path+"/comment", ".id="+id, "=comment="+comment)
	return err
}

// Move moves a resource to a new position
func (r *BaseResource) Move(id, destination string) error {
	return r.MoveContext(context.Background(), id, destination)
}

// MoveContext moves a resource to a new position with context
func (r *BaseResource) MoveContext(ctx context.Context, id, destination string) error {
	_, err := r.client.SendContext(ctx, r.path+"/move", ".id="+id, "=destination="+destination)
	return err
}

// Command executes a custom command
func (r *BaseResource) Command(cmd string, args ...string) ([]*Response, error) {
	return r.CommandContext(context.Background(), cmd, args...)
}

// CommandContext executes a custom command with context
func (r *BaseResource) CommandContext(ctx context.Context, cmd string, args ...string) ([]*Response, error) {
	return r.client.SendContext(ctx, r.path+"/"+cmd, args...)
}
