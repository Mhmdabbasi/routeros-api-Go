package routeros

import (
	"context"
)

// Queue handles RouterOS queue operations
type Queue struct {
	client *Client
}

// NewQueue creates a new Queue resource handler
func NewQueue(client *Client) *Queue {
	return &Queue{client: client}
}

// Simple handles simple queue operations
type QueueSimple struct {
	*BaseResource
}

// NewQueueSimple creates a new QueueSimple handler
func NewQueueSimple(client *Client) *QueueSimple {
	return &QueueSimple{
		BaseResource: NewBaseResource(client, "/queue/simple"),
	}
}

// Create creates a new simple queue
func (q *QueueSimple) Create(name, target string, maxLimit string) (string, error) {
	return q.CreateContext(context.Background(), name, target, maxLimit)
}

// CreateContext creates a new simple queue with context
func (q *QueueSimple) CreateContext(ctx context.Context, name, target string, maxLimit string) (string, error) {
	args := map[string]string{
		"name":   name,
		"target": target,
	}
	if maxLimit != "" {
		args["max-limit"] = maxLimit
	}
	return q.AddContext(ctx, args)
}

// SetMaxLimit sets the max limit for a queue
func (q *QueueSimple) SetMaxLimit(id, maxLimit string) error {
	return q.SetMaxLimitContext(context.Background(), id, maxLimit)
}

// SetMaxLimitContext sets the max limit for a queue with context
func (q *QueueSimple) SetMaxLimitContext(ctx context.Context, id, maxLimit string) error {
	return q.UpdateContext(ctx, id, map[string]string{"max-limit": maxLimit})
}

// SetTarget sets the target for a queue
func (q *QueueSimple) SetTarget(id, target string) error {
	return q.SetTargetContext(context.Background(), id, target)
}

// SetTargetContext sets the target for a queue with context
func (q *QueueSimple) SetTargetContext(ctx context.Context, id, target string) error {
	return q.UpdateContext(ctx, id, map[string]string{"target": target})
}

// SetPriority sets the priority for a queue
func (q *QueueSimple) SetPriority(id string, priority int) error {
	return q.SetPriorityContext(context.Background(), id, priority)
}

// SetPriorityContext sets the priority for a queue with context
func (q *QueueSimple) SetPriorityContext(ctx context.Context, id string, priority int) error {
	return q.UpdateContext(ctx, id, map[string]string{"priority": string(rune(priority))})
}

// Tree handles queue tree operations
type QueueTree struct {
	*BaseResource
}

// NewQueueTree creates a new QueueTree handler
func NewQueueTree(client *Client) *QueueTree {
	return &QueueTree{
		BaseResource: NewBaseResource(client, "/queue/tree"),
	}
}

// Create creates a new queue tree entry
func (q *QueueTree) Create(name, parent string, maxLimit string) (string, error) {
	return q.CreateContext(context.Background(), name, parent, maxLimit)
}

// CreateContext creates a new queue tree entry with context
func (q *QueueTree) CreateContext(ctx context.Context, name, parent string, maxLimit string) (string, error) {
	args := map[string]string{
		"name":   name,
		"parent": parent,
	}
	if maxLimit != "" {
		args["max-limit"] = maxLimit
	}
	return q.AddContext(ctx, args)
}

// Type handles queue type operations
type QueueType struct {
	*BaseResource
}

// NewQueueType creates a new QueueType handler
func NewQueueType(client *Client) *QueueType {
	return &QueueType{
		BaseResource: NewBaseResource(client, "/queue/type"),
	}
}

// Interface handles queue interface operations
type QueueInterface struct {
	*BaseResource
}

// NewQueueInterface creates a new QueueInterface handler
func NewQueueInterface(client *Client) *QueueInterface {
	return &QueueInterface{
		BaseResource: NewBaseResource(client, "/queue/interface"),
	}
}

// PCQ handles PCQ queue type operations
type PCQ struct {
	client *Client
}

// NewPCQ creates a new PCQ handler
func NewPCQ(client *Client) *PCQ {
	return &PCQ{client: client}
}
