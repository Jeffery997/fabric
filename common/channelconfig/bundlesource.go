/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	"sync/atomic"

	oldchannelconfig "github.com/hyperledger/fabric/common/config/channel"
	configtxapi "github.com/hyperledger/fabric/common/configtx/api"
	"github.com/hyperledger/fabric/common/policies"
	"github.com/hyperledger/fabric/msp"
)

// BundleSource stores a reference to the current configuration bundle
// It also provides a method to update this bundle.  The assorted methods
// largely pass through to the underlying bundle, but do so through an atomic pointer
// so that gross go-routine reads are not vulnerable to out-of-order execution memory
// type bugs.
type BundleSource struct {
	bundle atomic.Value
}

// NewBundleSource creates a new BundleSource with an initial Bundle value
func NewBundleSource(bundle *Bundle) *BundleSource {
	bs := &BundleSource{}
	bs.bundle.Store(bundle)
	return bs
}

// Update sets a new bundle as the bundle source
func (bs *BundleSource) Update(newBundle *Bundle) {
	bs.bundle.Store(newBundle)
}

// StableBundle returns a pointer to a stable Bundle.
// It is stable because calls to its assorted methods will always return the same
// result, as the underlying data structures are immutable.  For instance, calling
// BundleSource.Orderer() and BundleSource.MSPManager() to get first the list of orderer
// orgs, then querying the MSP for those org definitions could result in a bug because an
// update might replace the underlying Bundle in between.  Therefore, for operations
// which require consistency between the Bundle calls, the caller should first retrieve
// a StableBundle, then operate on it.
func (bs *BundleSource) StableBundle() *Bundle {
	return bs.bundle.Load().(*Bundle)
}

// PolicyManager returns the policy manager constructed for this config
func (bs *BundleSource) PolicyManager() policies.Manager {
	return bs.StableBundle().policyManager
}

// MSPManager returns the MSP manager constructed for this config
func (bs *BundleSource) MSPManager() msp.MSPManager {
	return bs.StableBundle().mspManager
}

// ChannelConfig returns the config.Channel for the chain
func (bs *BundleSource) ChannelConfig() oldchannelconfig.Channel {
	return bs.StableBundle().rootConfig.Channel()
}

// OrdererConfig returns the config.Orderer for the channel
// and whether the Orderer config exists
func (bs *BundleSource) OrdererConfig() (oldchannelconfig.Orderer, bool) {
	result := bs.StableBundle().rootConfig.Orderer()
	return result, result != nil
}

// ConsortiumsConfig() returns the config.Consortiums for the channel
// and whether the consortiums config exists
func (bs *BundleSource) ConsortiumsConfig() (oldchannelconfig.Consortiums, bool) {
	result := bs.StableBundle().rootConfig.Consortiums()
	return result, result != nil
}

// ApplicationConfig returns the configtxapplication.SharedConfig for the channel
// and whether the Application config exists
func (bs *BundleSource) ApplicationConfig() (oldchannelconfig.Application, bool) {
	result := bs.StableBundle().rootConfig.Application()
	return result, result != nil
}

// ConfigtxManager returns the configtx.Manager for the channel
func (bs *BundleSource) ConfigtxManager() configtxapi.Manager {
	return bs.StableBundle().configtxManager
}
