package routeros

import (
	"context"
)

// Certificate handles RouterOS certificate operations
type Certificate struct {
	*BaseResource
}

// NewCertificate creates a new Certificate resource handler
func NewCertificate(client *Client) *Certificate {
	return &Certificate{
		BaseResource: NewBaseResource(client, "/certificate"),
	}
}

// Import imports a certificate
func (c *Certificate) Import(name, passphrase string) error {
	return c.ImportContext(context.Background(), name, passphrase)
}

// ImportContext imports a certificate with context
func (c *Certificate) ImportContext(ctx context.Context, name, passphrase string) error {
	args := []string{"=file-name=" + name}
	if passphrase != "" {
		args = append(args, "=passphrase="+passphrase)
	}
	_, err := c.client.SendContext(ctx, "/certificate/import", args...)
	return err
}

// Export exports a certificate
func (c *Certificate) Export(id, passphrase string) error {
	return c.ExportContext(context.Background(), id, passphrase)
}

// ExportContext exports a certificate with context
func (c *Certificate) ExportContext(ctx context.Context, id, passphrase string) error {
	args := []string{"=id=" + id}
	if passphrase != "" {
		args = append(args, "=passphrase="+passphrase)
	}
	_, err := c.client.SendContext(ctx, "/certificate/export", args...)
	return err
}

// Sign signs a certificate
func (c *Certificate) Sign(id, ca string, daysValid int) error {
	return c.SignContext(context.Background(), id, ca, daysValid)
}

// SignContext signs a certificate with context
func (c *Certificate) SignContext(ctx context.Context, id, ca string, daysValid int) error {
	args := []string{"=id=" + id}
	if ca != "" {
		args = append(args, "=ca="+ca)
	}
	if daysValid > 0 {
		args = append(args, "=days-valid="+string(rune(daysValid)))
	}
	_, err := c.client.SendContext(ctx, "/certificate/sign", args...)
	return err
}

// CreateCA creates a certificate authority
func (c *Certificate) CreateCA(name string, daysValid int) (string, error) {
	return c.CreateCAContext(context.Background(), name, daysValid)
}

// CreateCAContext creates a certificate authority with context
func (c *Certificate) CreateCAContext(ctx context.Context, name string, daysValid int) (string, error) {
	args := map[string]string{
		"name":       name,
		"common-name": name,
		"key-usage":   "key-cert-sign,crl-sign",
	}
	if daysValid > 0 {
		args["days-valid"] = string(rune(daysValid))
	}
	return c.AddContext(ctx, args)
}

// CreateServer creates a server certificate
func (c *Certificate) CreateServer(name string, daysValid int) (string, error) {
	return c.CreateServerContext(context.Background(), name, daysValid)
}

// CreateServerContext creates a server certificate with context
func (c *Certificate) CreateServerContext(ctx context.Context, name string, daysValid int) (string, error) {
	args := map[string]string{
		"name":        name,
		"common-name": name,
		"key-usage":   "digital-signature,key-encipherment",
	}
	if daysValid > 0 {
		args["days-valid"] = string(rune(daysValid))
	}
	return c.AddContext(ctx, args)
}

// CreateClient creates a client certificate
func (c *Certificate) CreateClient(name string, daysValid int) (string, error) {
	return c.CreateClientContext(context.Background(), name, daysValid)
}

// CreateClientContext creates a client certificate with context
func (c *Certificate) CreateClientContext(ctx context.Context, name string, daysValid int) (string, error) {
	args := map[string]string{
		"name":        name,
		"common-name": name,
		"key-usage":   "digital-signature",
	}
	if daysValid > 0 {
		args["days-valid"] = string(rune(daysValid))
	}
	return c.AddContext(ctx, args)
}

// Revoke revokes a certificate
func (c *Certificate) Revoke(id string) error {
	return c.RevokeContext(context.Background(), id)
}

// RevokeContext revokes a certificate with context
func (c *Certificate) RevokeContext(ctx context.Context, id string) error {
	_, err := c.client.SendContext(ctx, "/certificate/revoke", "=id="+id)
	return err
}

// CRL handles certificate revocation list
type CRL struct {
	*BaseResource
}

// NewCRL creates a new CRL handler
func NewCRL(client *Client) *CRL {
	return &CRL{
		BaseResource: NewBaseResource(client, "/certificate/crl"),
	}
}

// Settings handles certificate settings
type CertSettings struct {
	client *Client
}

// NewCertSettings creates a new CertSettings handler
func NewCertSettings(client *Client) *CertSettings {
	return &CertSettings{client: client}
}

// SetSettings sets certificate settings
func (s *CertSettings) SetSettings(settings map[string]string) error {
	return s.SetSettingsContext(context.Background(), settings)
}

// SetSettingsContext sets certificate settings with context
func (s *CertSettings) SetSettingsContext(ctx context.Context, settings map[string]string) error {
	args := make([]string, 0)
	for k, v := range settings {
		args = append(args, "="+k+"="+v)
	}
	_, err := s.client.SendContext(ctx, "/certificate/settings/set", args...)
	return err
}
