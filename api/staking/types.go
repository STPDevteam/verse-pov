package staking

import (
	"bytes"
	"math/big"
	"sort"

	//"encoding/hex"
	"fmt"
	"time"

	"github.com/dfinlab/meter/meter"
	"github.com/dfinlab/meter/script/staking"
)

type Candidate struct {
	Name       string        `json:"name"`
	Addr       meter.Address `json:"addr"`   // the address for staking / reward
	PubKey     string        `json:"pubKey"` // node public key
	IPAddr     string        `json:"ipAddr"` // network addr
	Port       uint16        `json:"port"`
	TotalVotes string        `json:"totalVotes"` // total voting from all buckets
	Commission uint64        `json:"commission"` // commission rate unit "1e09"
	Buckets    []string      `json:"buckets"`    // all buckets voted for this candidate
}

func convertCandidateList(list *staking.CandidateList) []*Candidate {
	candidateList := make([]*Candidate, 0)
	for _, c := range list.ToList() {
		candidateList = append(candidateList, convertCandidate(c))
	}
	sort.SliceStable(candidateList, func(i, j int) bool {
		voteI := new(big.Int)
		voteJ := new(big.Int)
		voteI.SetString(candidateList[i].TotalVotes, 10)
		voteJ.SetString(candidateList[j].TotalVotes, 10)

		return voteI.Cmp(voteJ) >= 0
	})
	return candidateList
}

func convertCandidate(c staking.Candidate) *Candidate {
	buckets := make([]string, 0)
	for _, b := range c.Buckets {
		buckets = append(buckets, b.String())
	}
	return &Candidate{
		Name: string(bytes.Trim(c.Name[:], "\x00")),
		Addr: c.Addr,
		//PubKey:     hex.EncodeToString(c.PubKey),
		PubKey:     string(c.PubKey),
		IPAddr:     string(c.IPAddr),
		Port:       c.Port,
		TotalVotes: c.TotalVotes.String(),
		Commission: c.Commission,
		Buckets:    buckets,
	}
}

type Bucket struct {
	ID         string        `json:"id"`
	Owner      meter.Address `json:"owner"`
	Value      string        `json:"value"`
	Token      uint8         `json:"token"`
	Nonce      uint64        `json:"nonce"`
	CreateTime string        `json:"createTime"`

	Unbounded    bool          `json:"unbounded"`
	Candidate    meter.Address `json:"candidate"`
	Rate         uint8         `json:"rate"`
	Option       uint32        `json:"option"`
	BonusVotes   uint64        `json:"bonusVotes"`
	TotalVotes   string        `json:"totalVotes"`
	MatureTime   string        `json:"matureTime"`
	CalcLastTime string        `json:"calcLastTime"`
}

func convertBucketList(list *staking.BucketList) []*Bucket {
	bucketList := make([]*Bucket, 0)

	for _, b := range list.ToList() {
		bucketList = append(bucketList, &Bucket{
			ID:           b.BucketID.String(),
			Owner:        b.Owner,
			Value:        b.Value.String(),
			Token:        b.Token,
			Nonce:        b.Nonce,
			CreateTime:   fmt.Sprintln(time.Unix(int64(b.CreateTime), 0)),
			Unbounded:    b.Unbounded,
			Candidate:    b.Candidate,
			Rate:         b.Rate,
			Option:       b.Option,
			BonusVotes:   b.BonusVotes,
			TotalVotes:   b.TotalVotes.String(),
			MatureTime:   fmt.Sprintln(time.Unix(int64(b.MatureTime), 0)),
			CalcLastTime: fmt.Sprintln(time.Unix(int64(b.CalcLastTime), 0)),
		})
	}
	return bucketList
}

type Stakeholder struct {
	Holder     meter.Address `json:"holder"`
	TotalStake string        `json:"totalStake"`
	Buckets    []string      `json:"buckets"`
}

func convertStakeholderList(list *staking.StakeholderList) []*Stakeholder {
	stakeholderList := make([]*Stakeholder, 0)
	for _, s := range list.ToList() {
		stakeholderList = append(stakeholderList, convertStakeholder(s))
	}
	return stakeholderList
}

func convertStakeholder(s staking.Stakeholder) *Stakeholder {
	buckets := make([]string, 0)
	for _, b := range s.Buckets {
		buckets = append(buckets, b.String())
	}
	return &Stakeholder{
		Holder:     s.Holder,
		TotalStake: s.TotalStake.String(),
		Buckets:    buckets,
	}
}

type Distributor struct {
	Address meter.Address `json:"address"`
	Shares  uint64        `json:"shares"`
}

type Delegate struct {
	Name        string         `json:"name"`
	Address     meter.Address  `json:"address"`
	PubKey      string         `json:"pubKey"`
	VotingPower string         `json:"votingPower"`
	IPAddr      string         `json:"ipAddr"` // network addr
	Port        uint16         `json:"port"`
	Commission  uint64         `json""commissin"`
	DistList    []*Distributor `json:"distributors"`
}

func convertDelegateList(list *staking.DelegateList) []*Delegate {
	delegateList := make([]*Delegate, 0)
	for _, d := range list.GetDelegates() {
		delegateList = append(delegateList, convertDelegate(*d))
	}
	return delegateList
}

func convertDelegate(d staking.Delegate) *Delegate {
	dists := []*Distributor{}
	for _, dist := range d.DistList {
		dists = append(dists, &Distributor{
			Address: dist.Address,
			Shares:  dist.Shares,
		})
	}

	return &Delegate{
		Name:        string(bytes.Trim(d.Name[:], "\x00")),
		Address:     d.Address,
		PubKey:      string(d.PubKey),
		IPAddr:      string(d.IPAddr),
		Port:        d.Port,
		VotingPower: d.VotingPower.String(),
		Commission:  d.Commission,
		DistList:    dists,
	}
}

/********
type RewardInfo struct {
	Address meter.Address
	Amount  *big.Int
}

type ValidatorReward struct {
	Epoch            uint32
	TotalBaseRewards *big.Int
	ExpectDistribute *big.Int
	ActualDistribute *big.Int
	Info             []*RewardInfo
}
***/
type RewardInfo struct {
	Address meter.Address `json:"address"`
	Amount  uint64        `json:"amount"`
}

type ValidatorReward struct {
	Epoch            uint32        `json:"epoch"`
	BaseReward       uint64        `json:"baseReward"`
	ExpectDistribute uint64        `json:"expecteDistribute"`
	ActualDistribute uint64        `json:"actualDistribute`
	Info             []*RewardInfo `json:"info"`
}

func convertValidatorRewardList(list *staking.ValidatorRewardList) []*ValidatorReward {
	rewardList := make([]*ValidatorReward, 0)
	for _, r := range list.GetList() {
		rewardList = append(rewardList, convertValidatorReward(*r))
	}
	return rewardList
}

func convertValidatorReward(r staking.ValidatorReward) *ValidatorReward {
	info := []*RewardInfo{}
	for _, in := range r.Info {
		info = append(info, &RewardInfo{
			Address: in.Address,
			Amount:  in.Amount.Uint64(),
		})
	}

	return &ValidatorReward{
		Epoch:            r.Epoch,
		ExpectDistribute: r.ExpectDistribute.Uint64(),
		ActualDistribute: r.ActualDistribute.Uint64(),
		Info:             info,
	}
}
