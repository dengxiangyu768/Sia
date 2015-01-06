package host

import (
	"github.com/NebulousLabs/Sia/consensus"
	"github.com/NebulousLabs/Sia/encoding"
	"github.com/NebulousLabs/Sia/sia/hostdb"
)

// HostAnnounceSelf creates a host announcement transaction, adding
// information to the arbitrary data and then signing the transaction.
func (h *Host) AnnounceHost(freezeVolume consensus.Currency, freezeUnlockHeight consensus.BlockHeight) (t consensus.Transaction, err error) {
	// Get the encoded announcement based on the host settings.
	h.RLock()
	info := h.Settings
	h.RUnlock()
	announcement := string(encoding.MarshalAll(hostdb.HostAnnouncementPrefix, info))

	// Fill out the transaction.
	id, err := h.Wallet.RegisterTransaction(t)
	if err != nil {
		return
	}
	err = h.Wallet.FundTransaction(id, freezeVolume)
	if err != nil {
		return
	}
	info.SpendConditions, info.FreezeIndex, err = h.Wallet.AddTimelockedRefund(id, freezeVolume, freezeUnlockHeight)
	if err != nil {
		return
	}
	err = h.Wallet.AddArbitraryData(id, announcement)
	if err != nil {
		return
	}
	// TODO: Have the wallet manually add a fee? How should this be managed?
	t, err = h.Wallet.SignTransaction(id, true)
	if err != nil {
		return
	}

	return
}
