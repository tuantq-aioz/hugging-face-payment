package common

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgtype"
)

type AccAddress sdk.AccAddress

var _ sdk.Address = AccAddress{}

func (a AccAddress) String() string {
	return common.BytesToAddress(a.Bytes()).String()
}

func (a AccAddress) Equals(address sdk.Address) bool {
	return sdk.AccAddress(a).Equals(address)
}

func (a AccAddress) Empty() bool {
	return sdk.AccAddress(a).Empty()
}

func (a AccAddress) Marshal() ([]byte, error) {
	return sdk.AccAddress(a).Marshal()
}

func (a AccAddress) MarshalJSON() ([]byte, error) {
	return sdk.AccAddress(a).MarshalJSON()
}

func (a *AccAddress) UnmarshalJSON(data []byte) error {
	aa := &sdk.AccAddress{}
	if err := aa.UnmarshalJSON(data); err != nil {
		return err
	}
	*a = AccAddress(*aa)
	return nil
}

func (a AccAddress) Bytes() []byte {
	return sdk.AccAddress(a).Bytes()
}

func (a AccAddress) Format(s fmt.State, verb rune) {
	sdk.AccAddress(a).Format(s, verb)
}

// GORM

func (a *AccAddress) Scan(value interface{}) error {
	if len(value.(string)) == 0 {
		return nil
	}
	addr, err := SdkAccAddressFromString(value.(string))
	if err != nil {
		return err
	}
	*a = AccAddress(addr)
	return nil
}

func (a AccAddress) Value() (driver.Value, error) {
	return a.String(), nil
}

func (AccAddress) GormDataType() string {
	return "TEXT"
}

// PROTO

func (a AccAddress) Size() int {
	return len(a)
}

func (a AccAddress) MarshalTo(data []byte) (n int, err error) {
	if a == nil {
		return 0, nil
	}
	copy(data, a)
	return len(a), nil
}

func (a *AccAddress) Unmarshal(data []byte) error {
	aa := &sdk.AccAddress{}
	if err := aa.Unmarshal(data); err != nil {
		return err
	}
	*a = AccAddress(*aa)
	return nil
}

// SLICE

//type AccAddresses []AccAddress

// GORM

func (aa *AccAddresses) Scan(value interface{}) (err error) {
	var ta pgtype.TextArray
	if err = ta.Scan(value); err != nil {
		return err
	}

	var addrStrs []string
	if err = ta.AssignTo(&addrStrs); err != nil {
		return err
	}

	aa.Addresses = make([]AccAddress, len(ta.Elements))
	for i, e := range addrStrs {
		addr, err := SdkAccAddressFromString(e)
		if err != nil {
			return err
		}
		aa.Addresses[i] = AccAddress(addr)
	}

	return nil
}

func (aa AccAddresses) Value() (driver.Value, error) {
	var ta pgtype.TextArray

	addrStrs := make([]string, len(aa.Addresses))
	for i, e := range aa.Addresses {
		addrStrs[i] = e.String()
	}

	if err := ta.Set(addrStrs); err != nil {
		return nil, err
	}

	return ta.Value()
}

func (AccAddresses) GormDataType() string {
	return "TEXT[]"
}

// Utils

func AccAddressFromBech32(address string) (addr AccAddress, err error) {
	a, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	return AccAddress(a), nil
}

func AccAddressFromString(address string) (addr AccAddress, err error) {
	a, err := SdkAccAddressFromString(address)
	if err != nil {
		return nil, err
	}
	return AccAddress(a), nil
}

func SdkAccAddressFromString(address string) (addr sdk.AccAddress, err error) {
	addr, err = SdkAccAddressFromBech32(address)
	if err == nil {
		return addr, nil
	}

	if !common.IsHexAddress(address) {
		return nil, errors.New("invalid address: must provide bech32 or hex address")
	}

	return sdk.AccAddress(common.HexToAddress(address).Bytes()), nil
}

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func SdkAccAddressFromBech32(address string) (addr sdk.AccAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.AccAddress{}, errors.New("empty address string is not allowed")
	}

	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	bz, err := sdk.GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdk.AccAddress(bz), nil
}
