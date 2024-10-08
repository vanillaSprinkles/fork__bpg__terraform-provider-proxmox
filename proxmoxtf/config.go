/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmoxtf

import (
	"errors"
	"os"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
)

// ProviderConfiguration is the configuration for the provider.
type ProviderConfiguration struct {
	apiClient      api.Client
	sshClient      ssh.Client
	tmpDirOverride string
	idGenerator    cluster.IDGenerator
}

// NewProviderConfiguration creates a new provider configuration.
func NewProviderConfiguration(
	apiClient api.Client,
	sshClient ssh.Client,
	tmpDirOverride string,
	idCfg cluster.IDGeneratorConfig,
) (ProviderConfiguration, error) {
	cfg := ProviderConfiguration{
		apiClient:      apiClient,
		sshClient:      sshClient,
		tmpDirOverride: tmpDirOverride,
	}

	client, err := cfg.GetClient()
	if err != nil {
		return cfg, err
	}

	cfg.idGenerator = cluster.NewIDGenerator(client.Cluster(), idCfg)

	return cfg, nil
}

// GetClient returns the Proxmox API client.
func (c *ProviderConfiguration) GetClient() (proxmox.Client, error) {
	if c.apiClient == nil {
		return nil, errors.New(
			"you must specify the API access details in the provider configuration",
		)
	}

	if c.sshClient == nil {
		return nil, errors.New(
			"you must specify the SSH access details in the provider configuration",
		)
	}

	return proxmox.NewClient(c.apiClient, c.sshClient, c.tmpDirOverride), nil
}

// TempDir returns (possibly overridden) os.TempDir().
func (c *ProviderConfiguration) TempDir() string {
	if c.tmpDirOverride != "" {
		return c.tmpDirOverride
	}

	return os.TempDir()
}

// GetIDGenerator returns the IDGenerator.
func (c *ProviderConfiguration) GetIDGenerator() cluster.IDGenerator {
	return c.idGenerator
}
