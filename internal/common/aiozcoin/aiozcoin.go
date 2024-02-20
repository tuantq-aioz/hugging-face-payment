package aiozcoin

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
)

const (
	// DefaultDenom defines the default coin denomination used in AIOZ Network
	DefaultDenom string = "attoaioz"

	// BaseDenomUnit defines the base denomination unit for AIOZ.
	// 1 aioz = 1x10^{BaseDenomUnit} attoaioz
	BaseDenomUnit = 18

	aiozBaseDenom = DefaultDenom
)

type AiozCoin types.Coins

var AiozCoinZero = AiozCoin{}

func NewAiozCoinFromInt(i types.Int) AiozCoin {
	return AiozCoin(types.NewCoins(types.NewCoin(DefaultDenom, i)))
}

func (c *AiozCoin) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		err := fmt.Errorf("value is not string, %v", value)
		log.Error().Err(err).Send()
		return err
	}

	if strValue == "" || strValue == "0" {
		*c = AiozCoin{}
	} else {
		amount, ok := types.NewIntFromString(value.(string))
		if !ok {
			err := errors.New("invalid string " + value.(string))
			log.Error().Err(err).Send()
			return err
		}
		coins := types.Coins{
			types.Coin{
				Denom:  aiozBaseDenom,
				Amount: amount,
			},
		}

		// coin, err := types.ParseCoinsNormalized(value.(string) + aiozBaseDenom)
		// if err != nil {
		// 	log.Error().Err(err).Send()
		// 	return err
		// }
		*c = AiozCoin(coins)
	}
	return nil
}

func (c AiozCoin) Value() (driver.Value, error) {
	for _, coin := range c {
		if coin.Denom == aiozBaseDenom {
			return coin.Amount.String(), nil
		}
	}
	return "0", nil
}

func (AiozCoin) GormDataType() string {
	return "NUMERIC"
}

func (c AiozCoin) IsAnyGT(coinsB AiozCoin) bool {
	return types.Coins(c).IsAnyGT(types.Coins(coinsB))
}

func (c AiozCoin) IsAllGTE(coinsB AiozCoin) bool {
	return types.Coins(c).IsAllGTE(types.Coins(coinsB))
}

func (c AiozCoin) Add(coinsB AiozCoin) AiozCoin {
	return AiozCoin(types.Coins(c).Add(types.Coins(coinsB)...))
}

func (c AiozCoin) Sub(coinsB AiozCoin) AiozCoin {
	return AiozCoin(types.Coins(c).Sub(types.Coins(coinsB)...))
}

func (c AiozCoin) SafeSub(coinsB AiozCoin) (AiozCoin, bool) {
	r, n := types.Coins(c).SafeSub(types.Coins(coinsB)...)
	return AiozCoin(r), n
}

func (c AiozCoin) IsZero() bool {
	return types.Coins(c).IsZero()
}

func (c AiozCoin) Size() int {
	amount := types.Coins(c).AmountOf(aiozBaseDenom)
	return amount.Size()
}

func (c AiozCoin) MarshalTo(dst []byte) (int, error) {
	amount := types.Coins(c).AmountOf(aiozBaseDenom)
	return amount.MarshalTo(dst)
}

func (c AiozCoin) Marshal() ([]byte, error) {
	amount := types.Coins(c).AmountOf(aiozBaseDenom)
	return amount.Marshal()
}

func (c *AiozCoin) Unmarshal(dAtA []byte) error {
	var amount types.Int
	if err := amount.Unmarshal(dAtA); err != nil {
		return err
	}
	if amount.IsZero() {
		*c = AiozCoin{}
		return nil
	}
	*c = AiozCoin{types.Coin{Amount: amount, Denom: aiozBaseDenom}}
	return nil
}

func (c AiozCoin) Clone() (AiozCoin, error) {
	t, err := c.Marshal()
	if err != nil {
		return nil, err
	}

	var c2 AiozCoin
	err = c2.Unmarshal(t)

	return c2, err
}

func (c AiozCoin) AiozAmount() types.Int {
	return types.Coins(c).AmountOf(aiozBaseDenom)
}
