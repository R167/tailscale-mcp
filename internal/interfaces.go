package internal

import (
	"context"

	tailscale "tailscale.com/client/tailscale/v2"
)

// TailscaleClient defines the interface for Tailscale API operations
type TailscaleClient interface {
	Devices() DevicesResource
	PolicyFile() PolicyFileResource
	Keys() KeysResource
}

// DevicesResource defines the interface for device operations
type DevicesResource interface {
	ListWithAllFields(ctx context.Context) ([]tailscale.Device, error)
	GetWithAllFields(ctx context.Context, deviceID string) (*tailscale.Device, error)
	SubnetRoutes(ctx context.Context, deviceID string) (*tailscale.DeviceRoutes, error)
}

// PolicyFileResource defines the interface for ACL operations
type PolicyFileResource interface {
	Get(ctx context.Context) (*tailscale.ACL, error)
}

// KeysResource defines the interface for API key operations
type KeysResource interface {
	List(ctx context.Context, all bool) ([]tailscale.Key, error)
}

// TailscaleClientAdapter wraps the real Tailscale client to implement our interface
type TailscaleClientAdapter struct {
	*tailscale.Client
}

func (t *TailscaleClientAdapter) Devices() DevicesResource {
	return &DevicesResourceAdapter{t.Client.Devices()}
}

func (t *TailscaleClientAdapter) PolicyFile() PolicyFileResource {
	return &PolicyFileResourceAdapter{t.Client.PolicyFile()}
}

func (t *TailscaleClientAdapter) Keys() KeysResource {
	return &KeysResourceAdapter{t.Client.Keys()}
}

// DevicesResourceAdapter adapts the real DevicesResource
type DevicesResourceAdapter struct {
	*tailscale.DevicesResource
}

func (d *DevicesResourceAdapter) ListWithAllFields(ctx context.Context) ([]tailscale.Device, error) {
	return d.DevicesResource.ListWithAllFields(ctx)
}

func (d *DevicesResourceAdapter) GetWithAllFields(ctx context.Context, deviceID string) (*tailscale.Device, error) {
	return d.DevicesResource.GetWithAllFields(ctx, deviceID)
}

func (d *DevicesResourceAdapter) SubnetRoutes(ctx context.Context, deviceID string) (*tailscale.DeviceRoutes, error) {
	return d.DevicesResource.SubnetRoutes(ctx, deviceID)
}

// PolicyFileResourceAdapter adapts the real PolicyFileResource
type PolicyFileResourceAdapter struct {
	*tailscale.PolicyFileResource
}

func (p *PolicyFileResourceAdapter) Get(ctx context.Context) (*tailscale.ACL, error) {
	return p.PolicyFileResource.Get(ctx)
}

// KeysResourceAdapter adapts the real KeysResource
type KeysResourceAdapter struct {
	*tailscale.KeysResource
}

func (k *KeysResourceAdapter) List(ctx context.Context, all bool) ([]tailscale.Key, error) {
	return k.KeysResource.List(ctx, all)
}
