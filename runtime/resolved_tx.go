// Copyright (c) 2020 The Meter.io developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package runtime

import (
	"fmt"
	"math/big"

	"github.com/dfinlab/meter/builtin"
	"github.com/dfinlab/meter/meter"
	"github.com/dfinlab/meter/state"
	"github.com/dfinlab/meter/tx"
	"github.com/dfinlab/meter/xenv"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/pkg/errors"
)

// ResolvedTransaction resolve the transaction according to given state.
type ResolvedTransaction struct {
	tx           *tx.Transaction
	Origin       meter.Address
	IntrinsicGas uint64
	Clauses      []*tx.Clause
}

// ResolveTransaction resolves the transaction and performs basic validation.
func ResolveTransaction(tx *tx.Transaction) (*ResolvedTransaction, error) {
	origin, err := tx.Signer()
	if err != nil {
		return nil, err
	}
	intrinsicGas, err := tx.IntrinsicGas()
	if err != nil {
		return nil, err
	}
	if tx.Gas() < intrinsicGas {
		return nil, fmt.Errorf("intrinsic gas (%d) exceeds provided gas (%d)", intrinsicGas, tx.Gas())
	}

	clauses := tx.Clauses()
	sumValue := new(big.Int)
	for _, clause := range clauses {
		value := clause.Value()
		if value.Sign() < 0 {
			return nil, errors.New("clause with negative value")
		}

		sumValue.Add(sumValue, value)
		if sumValue.Cmp(math.MaxBig256) > 0 {
			return nil, errors.New("tx value too large")
		}
	}

	return &ResolvedTransaction{
		tx,
		origin,
		intrinsicGas,
		clauses,
	}, nil
}

// CommonTo returns common 'To' field of clauses if any.
// Nil returned if no common 'To'.
func (r *ResolvedTransaction) CommonTo() *meter.Address {
	if len(r.Clauses) == 0 {
		return nil
	}

	firstTo := r.Clauses[0].To()
	if firstTo == nil {
		return nil
	}

	for _, clause := range r.Clauses[1:] {
		to := clause.To()
		if to == nil {
			return nil
		}
		if *to != *firstTo {
			return nil
		}
	}
	return firstTo
}

// BuyGas consumes energy to buy gas, to prepare for execution.
func (r *ResolvedTransaction) BuyGas(state *state.State, blockTime uint64) (
	baseGasPrice *big.Int,
	gasPrice *big.Int,
	payer meter.Address,
	returnGas func(uint64), err error) {

	baseGasPrice = builtin.Params.Native(state).Get(meter.KeyBaseGasPrice)
	gasPrice = r.tx.GasPrice(baseGasPrice)

	//energy := builtin.MeterTracker.Native(state)
	doReturnGas := func(rgas uint64) *big.Int {
		returnedEnergy := new(big.Int).Mul(new(big.Int).SetUint64(rgas), gasPrice)
		//energy.Add(payer, returnedEnergy)
		state.AddEnergy(payer, returnedEnergy)
		return returnedEnergy
	}

	prepaid := new(big.Int).Mul(new(big.Int).SetUint64(r.tx.Gas()), gasPrice)
	commonTo := r.CommonTo()
	if commonTo != nil {
		binding := builtin.Prototype.Native(state).Bind(*commonTo)
		credit := binding.UserCredit(r.Origin, blockTime)
		if credit.Cmp(prepaid) >= 0 {
			doReturnGasAndSetCredit := func(rgas uint64) {
				returnedEnergy := doReturnGas(rgas)
				usedEnergy := new(big.Int).Sub(prepaid, returnedEnergy)
				binding.SetUserCredit(r.Origin, new(big.Int).Sub(credit, usedEnergy), blockTime)
			}
			// has enough credit
			if sponsor := binding.CurrentSponsor(); binding.IsSponsor(sponsor) {
				// deduct from sponsor, if any
				/*if energy.Sub(sponsor, prepaid) {*/
				if state.SubEnergy(sponsor, prepaid) {
					return baseGasPrice, gasPrice, sponsor, doReturnGasAndSetCredit, nil
				}
			}
			// deduct from To
			/*if energy.Sub(*commonTo, prepaid) {*/
			if state.SubEnergy(*commonTo, prepaid) {
				return baseGasPrice, gasPrice, *commonTo, doReturnGasAndSetCredit, nil
			}
		}
	}

	// fallback to deduct from tx origin
	// XXX reward transaction (Origin is nil) should not to deduct
	if !r.Origin.IsZero() {
		/*if energy.Sub(r.Origin, prepaid) {*/
		if state.SubEnergy(r.Origin, prepaid) {
			return baseGasPrice, gasPrice, r.Origin, func(rgas uint64) { doReturnGas(rgas) }, nil
		} else {
			return nil, nil, meter.Address{}, nil, errors.New("insufficient energy")
		}
	} else {
		return baseGasPrice, gasPrice, r.Origin, func(rgas uint64) { doReturnGas(rgas) }, nil
	}
}

// ToContext create a tx context object.
func (r *ResolvedTransaction) ToContext(gasPrice *big.Int, blockNumber uint32, getID func(uint32) meter.Bytes32) *xenv.TransactionContext {
	return &xenv.TransactionContext{
		ID:         r.tx.ID(),
		Origin:     r.Origin,
		GasPrice:   gasPrice,
		ProvedWork: r.tx.ProvedWork(blockNumber, getID),
		BlockRef:   r.tx.BlockRef(),
		Expiration: r.tx.Expiration(),
		Nonce:      r.tx.Nonce(),
	}
}
