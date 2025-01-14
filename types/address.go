package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"plugchain-sdk-go/utils/bech32"
)

//Bech32 Address
type AccAddress []byte

//Return by address accaddress
func AccAddressFromBech32(address string) (AccAddress, Error) {
	bech32PrefixAccAddr := GetAddrPrefixCfg().GetBech32AccountAddrPrefix()
	bz, err := bech32.GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		return nil, Wrap(err)
	}

	return AccAddress(bz), nil
}

//Verify address
func ValidateAccAddress(address string) Error {
	bech32PrefixAccAddr := GetAddrPrefixCfg().GetBech32AccountAddrPrefix()
	_, err := bech32.GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		return Wrap(err)
	}
	return nil
}

//Get by address AccAddress
func MustAccAddressFromBech32(address string) AccAddress {
	addr, err := AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}
	return addr
}

//return string address
func (aa AccAddress) String() string {
	if aa.Empty() {
		return ""
	}

	bech32PrefixAccAddr := GetAddrPrefixCfg().GetBech32AccountAddrPrefix()
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixAccAddr, aa.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

//Compare whether the two addresses are the same
func (aa AccAddress) Equals(aa2 AccAddress) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

//Determine whether it is an empty address
func (aa AccAddress) Empty() bool {
	if aa == nil {
		return true
	}

	aa2 := AccAddress{}
	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa AccAddress) Marshal() ([]byte, error) {
	return aa, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *AccAddress) Unmarshal(data []byte) error {
	*aa = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa AccAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *AccAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := AccAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// Returns the original byte
func (aa AccAddress) Bytes() []byte {
	return aa
}

type ValAddress []byte

func ValAddressFromBech32(address string) (ValAddress, Error) {
	bech32PrefixValAddr := GetAddrPrefixCfg().GetBech32ValidatorAddrPrefix()
	bz, err := bech32.GetFromBech32(address, bech32PrefixValAddr)
	if err != nil {
		return nil, Wrap(err)
	}

	return ValAddress(bz), nil
}

// Returns boolean for whether two ValAddresses are equal
func (va ValAddress) Equals(va2 ValAddress) bool {
	if va.Empty() && va2.Empty() {
		return true
	}

	return bytes.Equal(va.Bytes(), va2.Bytes())
}

// Returns boolean for whether an AccAddress is empty
func (va ValAddress) Empty() bool {
	if va == nil {
		return true
	}

	va2 := ValAddress{}
	return bytes.Equal(va.Bytes(), va2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (va ValAddress) Marshal() ([]byte, error) {
	return va, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (va *ValAddress) Unmarshal(data []byte) error {
	*va = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (va ValAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(va.String())
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (va *ValAddress) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	va2, err := ValAddressFromBech32(s)
	if err != nil {
		return err
	}

	*va = va2
	return nil
}

// Bytes returns the raw address bytes.
func (va ValAddress) Bytes() []byte {
	return va
}

// String implements the Stringer interface.
func (va ValAddress) String() string {
	bech32PrefixValAddr := GetAddrPrefixCfg().GetBech32ValidatorAddrPrefix()
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixValAddr, va.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (va ValAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(va.String()))
	case 'p':
		_, _ = s.Write([]byte(fmt.Sprintf("%p", va)))
	default:
		_, _ = s.Write([]byte(fmt.Sprintf("%X", []byte(va))))
	}
}

//Address byte
type ConsAddress []byte

func (ca ConsAddress) String() string {
	bech32PrefixConsAddr := GetAddrPrefixCfg().GetBech32ConsensusAddrPrefix()
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixConsAddr, ca.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

func (ca ConsAddress) Bytes() []byte {
	return ca
}

//Convert address to ConsAddress format
func ConsAddressFromHex(address string) (addr ConsAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding Bech32 address failed: must provide an address")
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return ConsAddress(bz), nil
}

type Bech32PubKeyType string

const (
	Bech32PubKeyTypeAccPub  Bech32PubKeyType = "accpub"
	Bech32PubKeyTypeValPub  Bech32PubKeyType = "valpub"
	Bech32PubKeyTypeConsPub Bech32PubKeyType = "conspub"
)

// Bech32ifyPubKey returns a Bech32 encoded string containing the appropriate
// prefix based on the key type provided for a given PublicKey.
func Bech32ifyPubKey(pkt Bech32PubKeyType, pubkey TmPubKey) (string, error) {
	var bech32Prefix string

	switch pkt {
	case Bech32PubKeyTypeAccPub:
		bech32Prefix = GetAddrPrefixCfg().GetBech32AccountPubPrefix()

	case Bech32PubKeyTypeValPub:
		bech32Prefix = GetAddrPrefixCfg().GetBech32ValidatorPubPrefix()

	case Bech32PubKeyTypeConsPub:
		bech32Prefix = GetAddrPrefixCfg().GetBech32ConsensusPubPrefix()

	}

	return bech32.ConvertAndEncode(bech32Prefix, pubkey.Bytes())
}

func GetPubKeyFromBech32(pkt Bech32PubKeyType, pubkeyStr string) (TmPubKey, error) {
	var bech32Prefix string

	switch pkt {
	case Bech32PubKeyTypeAccPub:
		bech32Prefix = GetAddrPrefixCfg().GetBech32AccountPubPrefix()

	case Bech32PubKeyTypeValPub:
		bech32Prefix = GetAddrPrefixCfg().GetBech32ValidatorPubPrefix()

	case Bech32PubKeyTypeConsPub:
		bech32Prefix = GetAddrPrefixCfg().GetBech32ConsensusPubPrefix()

	}

	bz, err := GetFromBech32(pubkeyStr, bech32Prefix)
	if err != nil {
		return nil, err
	}

	return PubKeyFromBytes(bz)
}

func GetFromBech32(bech32str, prefix string) ([]byte, error) {
	if len(bech32str) == 0 {
		return nil, errors.New("decoding Bech32 address failed: must provide an address")
	}

	hrp, bz, err := bech32.DecodeAndConvert(bech32str)
	if err != nil {
		return nil, err
	}

	if hrp != prefix {
		return nil, fmt.Errorf("invalid Bech32 prefix; expected %s, got %s", prefix, hrp)
	}

	return bz, nil
}
