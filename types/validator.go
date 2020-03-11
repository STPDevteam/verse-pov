package types

import (
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"

	bls "github.com/dfinlab/meter/crypto/multi_sig"
	"github.com/dfinlab/meter/meter"
	"github.com/ethereum/go-ethereum/crypto"
)

// Volatile state for each Validator
// NOTE: The Accum is not included in Validator.Hash();
// make sure to update that method if changes are made here
type Validator struct {
	Name        string
	Address     meter.Address
	PubKey      ecdsa.PublicKey
	BlsPubKey   bls.PublicKey
	VotingPower int64
	NetAddr     NetAddress
	CommitKey   []byte
}

func NewValidator(name string, address meter.Address, pubKey ecdsa.PublicKey, blsPub bls.PublicKey, votingPower int64) *Validator {
	return &Validator{
		Name:        name,
		Address:     address,
		PubKey:      pubKey,
		BlsPubKey:   blsPub,
		VotingPower: votingPower,
	}
}

// Creates a new copy of the validator so we can mutate accum.
// Panics if the validator is nil.
func (v *Validator) Copy() *Validator {
	vCopy := *v
	return &vCopy
}

func (v *Validator) String() string {
	if v == nil {
		return "nil-Validator"
	}
	pubkey := base64.StdEncoding.EncodeToString(crypto.FromECDSAPub(&v.PubKey))
	pubkey = pubkey[:4] + "..." + pubkey[len(pubkey)-4:]
	return fmt.Sprintf("%15v (address:%v ip:%v pubkey:%v)",
		v.Name,
		v.Address.String(),
		v.NetAddr.IP.String(),
		pubkey,
	)
}
