package script

import (
	//"github.com/dfinlab/meter/script/auction"
	"github.com/dfinlab/meter/script/staking"
)

const (
	STAKING_MODULE_NAME = string("staking")
	STAKING_MODULE_ID   = uint32(1000)

	AUCTION_MODULE_NAME = string("auction")
	AUCTION_MODULE_ID   = uint32(1001)
)

func ModuleStakingInit(se *ScriptEngine) *staking.Staking {
	stk := staking.NewStaking(se.chain, se.stateCreator)
	if stk == nil {
		panic("init staking module failed")
	}

	mod := &Module{
		modName:    STAKING_MODULE_NAME,
		modID:      STAKING_MODULE_ID,
		modHandler: stk.PrepareStakingHandler(),
	}
	if err := se.modReg.Register(STAKING_MODULE_ID, mod); err != nil {
		panic("register staking module failed")
	}

	stk.Start()
	se.logger.Info("ScriptEngine", "started moudle", mod.modName)
	return stk
}

/***
func ModuleAuctionInit(se *ScriptEngine) *auction.Auction {
	a := auction.NewAuction(se.chain, se.stateCreator)
	if a == nil {
		panic("init acution module failed")
	}

	mod := &Module{
		modName:    AUCTION_MODULE_NAME,
		modID:      AUCTION_MODULE_ID,
		modHandler: a.PrepareAuctionHandler(),
	}
	if err := se.modReg.Register(AUCTION_MODULE_ID, mod); err != nil {
		panic("register auction module failed")
	}

	a.Start()
	se.logger.Info("ScriptEngine", "started moudle", mod.modName)
	return a
}
***/
