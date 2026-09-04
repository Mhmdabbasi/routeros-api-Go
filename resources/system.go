package routeros

import (
	"context"
)

// System provides access to system-related RouterOS commands
type System struct {
	client *Client
	*BaseResource
}

// NewSystem creates a new System resource handler
func NewSystem(client *Client) *System {
	return &System{
		client:       client,
		BaseResource: NewBaseResource(client, "/system"),
	}
}

// Resource returns the base resource for custom operations
func (s *System) Resource() *BaseResource {
	return s.BaseResource
}

// Identity handles router identity operations
type Identity struct {
	client *Client
}

// NewIdentity creates a new Identity handler
func NewIdentity(client *Client) *Identity {
	return &Identity{client: client}
}

// Get returns the current router identity
func (i *Identity) Get() (string, error) {
	return i.GetContext(context.Background())
}

// GetContext returns the current router identity with context
func (i *Identity) GetContext(ctx context.Context) (string, error) {
	responses, err := i.client.SendContext(ctx, "/system/identity/print")
	if err != nil {
		return "", err
	}
	parser := NewResponseParser(responses)
	return parser.Value("name"), nil
}

// Set sets the router identity name
func (i *Identity) Set(name string) error {
	return i.SetContext(context.Background(), name)
}

// SetContext sets the router identity name with context
func (i *Identity) SetContext(ctx context.Context, name string) error {
	_, err := i.client.SendContext(ctx, "/system/identity/set", "=name="+name)
	return err
}

// Clock handles system clock operations
type Clock struct {
	client *Client
}

// NewClock creates a new Clock handler
func NewClock(client *Client) *Clock {
	return &Clock{client: client}
}

// Time returns the current system time
func (c *Clock) Time() (map[string]string, error) {
	return c.TimeContext(context.Background())
}

// TimeContext returns the current system time with context
func (c *Clock) TimeContext(ctx context.Context) (map[string]string, error) {
	responses, err := c.client.SendContext(ctx, "/system/clock/print")
	if err != nil {
		return nil, err
	}
	parser := NewResponseParser(responses)
	return parser.FirstData(), nil
}

// SetTime sets the system time
func (c *Clock) SetTime(time string) error {
	return c.SetTimeContext(context.Background(), time)
}

// SetTimeContext sets the system time with context
func (c *Clock) SetTimeContext(ctx context.Context, time string) error {
	_, err := c.client.SendContext(ctx, "/system/clock/set", "=time="+time)
	return err
}

// SetDate sets the system date
func (c *Clock) SetDate(date string) error {
	return c.SetDateContext(context.Background(), date)
}

// SetDateContext sets the system date with context
func (c *Clock) SetDateContext(ctx context.Context, date string) error {
	_, err := c.client.SendContext(ctx, "/system/clock/set", "=date="+date)
	return err
}

// SetTimeZone sets the system timezone
func (c *Clock) SetTimeZone(tz string) error {
	return c.SetTimeZoneContext(context.Background(), tz)
}

// SetTimeZoneContext sets the system timezone with context
func (c *Clock) SetTimeZoneContext(ctx context.Context, tz string) error {
	_, err := c.client.SendContext(ctx, "/system/clock/set", "=time-zone-name="+tz)
	return err
}

// SetNTP enables or disables NTP client
func (c *Clock) SetNTP(enabled bool) error {
	return c.SetNTPContext(context.Background(), enabled)
}

// SetNTPContext enables or disables NTP client with context
func (c *Clock) SetNTPContext(ctx context.Context, enabled bool) error {
	val := "no"
	if enabled {
		val = "yes"
	}
	_, err := c.client.SendContext(ctx, "/system/clock/set", "=gmt-offset="+val)
	return err
}

// Reboot restarts the router
func (s *System) Reboot() error {
	return s.RebootContext(context.Background())
}

// RebootContext restarts the router with context
func (s *System) RebootContext(ctx context.Context) error {
	_, err := s.client.SendContext(ctx, "/system/reboot")
	return err
}

// Shutdown shuts down the router
func (s *System) Shutdown() error {
	return s.ShutdownContext(context.Background())
}

// ShutdownContext shuts down the router with context
func (s *System) ShutdownContext(ctx context.Context) error {
	_, err := s.client.SendContext(ctx, "/system/shutdown")
	return err
}

// Resource returns system resource information
func (s *System) Resource() (*Resource, error) {
	return s.ResourceContext(context.Background())
}

// ResourceContext returns system resource information with context
func (s *System) ResourceContext(ctx context.Context) (*Resource, error) {
	responses, err := s.client.SendContext(ctx, "/system/resource/print")
	if err != nil {
		return nil, err
	}
	list := NewResourceList(responses)
	return list.First(), nil
}

// Health returns system health information
func (s *System) Health() (*ResourceList, error) {
	return s.HealthContext(context.Background())
}

// HealthContext returns system health information with context
func (s *System) HealthContext(ctx context.Context) (*ResourceList, error) {
	responses, err := s.client.SendContext(ctx, "/system/health/print")
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// Routerboard returns routerboard information
func (s *System) Routerboard() (*Resource, error) {
	return s.RouterboardContext(context.Background())
}

// RouterboardContext returns routerboard information with context
func (s *System) RouterboardContext(ctx context.Context) (*Resource, error) {
	responses, err := s.client.SendContext(ctx, "/system/routerboard/print")
	if err != nil {
		return nil, err
	}
	list := NewResourceList(responses)
	return list.First(), nil
}

// Packages returns installed packages
func (s *System) Packages() (*ResourceList, error) {
	return s.PackagesContext(context.Background())
}

// PackagesContext returns installed packages with context
func (s *System) PackagesContext(ctx context.Context) (*ResourceList, error) {
	responses, err := s.client.SendContext(ctx, "/system/package/print")
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// History returns command history
func (s *System) History() (*ResourceList, error) {
	return s.HistoryContext(context.Background())
}

// HistoryContext returns command history with context
func (s *System) HistoryContext(ctx context.Context) (*ResourceList, error) {
	responses, err := s.client.SendContext(ctx, "/system/history/print")
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// Undo undoes the last command
func (s *System) Undo() error {
	return s.UndoContext(context.Background())
}

// UndoContext undoes the last command with context
func (s *System) UndoContext(ctx context.Context) error {
	_, err := s.client.SendContext(ctx, "/system/undo")
	return err
}

// Redo redoes the last undone command
func (s *System) Redo() error {
	return s.RedoContext(context.Background())
}

// RedoContext redoes the last undone command with context
func (s *System) RedoContext(ctx context.Context) error {
	_, err := s.client.SendContext(ctx, "/system/redo")
	return err
}

// Backup handles backup operations
type Backup struct {
	client *Client
}

// NewBackup creates a new Backup handler
func NewBackup(client *Client) *Backup {
	return &Backup{client: client}
}

// List returns all backup files
func (b *Backup) List() (*ResourceList, error) {
	return b.ListContext(context.Background())
}

// ListContext returns all backup files with context
func (b *Backup) ListContext(ctx context.Context) (*ResourceList, error) {
	responses, err := b.client.SendContext(ctx, "/system/backup/print")
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// Save creates a new backup
func (b *Backup) Save(name string) error {
	return b.SaveContext(context.Background(), name)
}

// SaveContext creates a new backup with context
func (b *Backup) SaveContext(ctx context.Context, name string) error {
	args := []string{}
	if name != "" {
		args = append(args, "=name="+name)
	}
	_, err := b.client.SendContext(ctx, "/system/backup/save", args...)
	return err
}

// Load loads a backup
func (b *Backup) Load(name string) error {
	return b.LoadContext(context.Background(), name)
}

// LoadContext loads a backup with context
func (b *Backup) LoadContext(ctx context.Context, name string) error {
	_, err := b.client.SendContext(ctx, "/system/backup/load", "=name="+name)
	return err
}

// Script handles system scripts
type Script struct {
	*BaseResource
}

// NewScript creates a new Script handler
func NewScript(client *Client) *Script {
	return &Script{
		BaseResource: NewBaseResource(client, "/system/script"),
	}
}

// Run executes a script by name
func (s *Script) Run(name string) error {
	return s.RunContext(context.Background(), name)
}

// RunContext executes a script by name with context
func (s *Script) RunContext(ctx context.Context, name string) error {
	_, err := s.client.SendContext(ctx, "/system/script/run", "=source="+name)
	return err
}

// Environment handles script environment variables
type Environment struct {
	*BaseResource
}

// NewEnvironment creates a new Environment handler
func NewEnvironment(client *Client) *Environment {
	return &Environment{
		BaseResource: NewBaseResource(client, "/system/script/environment"),
	}
}

// Scheduler handles scheduled tasks
type Scheduler struct {
	*BaseResource
}

// NewScheduler creates a new Scheduler handler
func NewScheduler(client *Client) *Scheduler {
	return &Scheduler{
		BaseResource: NewBaseResource(client, "/system/scheduler"),
	}
}

// User handles system users
type User struct {
	*BaseResource
}

// NewUser creates a new User handler
func NewUser(client *Client) *User {
	return &User{
		BaseResource: NewBaseResource(client, "/user"),
	}
}

// SetPassword sets a user's password
func (u *User) SetPassword(username, password string) error {
	return u.SetPasswordContext(context.Background(), username, password)
}

// SetPasswordContext sets a user's password with context
func (u *User) SetPasswordContext(ctx context.Context, username, password string) error {
	user, err := u.FindByNameContext(ctx, username)
	if err != nil {
		return err
	}
	_, err = u.client.SendContext(ctx, "/user/set", ".id="+user.ID, "=password="+password)
	return err
}

// Group handles user groups
type UserGroup struct {
	*BaseResource
}

// NewUserGroup creates a new UserGroup handler
func NewUserGroup(client *Client) *UserGroup {
	return &UserGroup{
		BaseResource: NewBaseResource(client, "/user/group"),
	}
}

// Logs returns system logs
func (s *System) Logs() (*ResourceList, error) {
	return s.LogsContext(context.Background())
}

// LogsContext returns system logs with context
func (s *System) LogsContext(ctx context.Context) (*ResourceList, error) {
	responses, err := s.client.SendContext(ctx, "/log/print")
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// Info contains commonly used system information
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

// GetInfo returns system information
func (s *System) GetInfo() (*Info, error) {
	return s.GetInfoContext(context.Background())
}

// GetInfoContext returns system information with context
func (s *System) GetInfoContext(ctx context.Context) (*Info, error) {
	// Get identity
	identity, err := NewIdentity(s.client).GetContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get resource info
	resource, err := s.ResourceContext(ctx)
	if err != nil {
		return nil, err
	}

	info := &Info{
		Identity:    identity,
		Version:     resource.Get("version"),
		BuildTime:   resource.Get("build-time"),
		FreeMemory:  resource.Get("free-memory"),
		TotalMemory: resource.Get("total-memory"),
		CPU:         resource.Get("cpu"),
		CPUCount:    resource.Get("cpu-count"),
		Uptime:      resource.Get("uptime"),
		ArchName:    resource.Get("architecture-name"),
		BoardName:   resource.Get("board-name"),
		Model:       resource.Get("model"),
		Serial:      resource.Get("serial-number"),
	}

	return info, nil
}

// Export exports configuration
func (s *System) Export() (string, error) {
	return s.ExportContext(context.Background())
}

// ExportContext exports configuration with context
func (s *System) ExportContext(ctx context.Context) (string, error) {
	responses, err := s.client.SendContext(ctx, "/export")
	if err != nil {
		return "", err
	}

	var export string
	for _, resp := range responses {
		if resp.Type == ResponseReply {
			// Export data might be in different fields
			for _, v := range resp.Data {
				export += v + "\n"
			}
		}
	}
	return export, nil
}

// ResetConfiguration resets the router configuration
func (s *System) ResetConfiguration() error {
	return s.ResetConfigurationContext(context.Background())
}

// ResetConfigurationContext resets the router configuration with context
func (s *System) ResetConfigurationContext(ctx context.Context) error {
	_, err := s.client.SendContext(ctx, "/system/reset-configuration")
	return err
}
