package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/tendermint/tendermint/crypto"
)

const (
	// AddrLen defines a valid address length
	AddrLen = 20
)

// ----------------------------------------------------------------------------
// account
// ----------------------------------------------------------------------------

// AccAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type AccAddress []byte

// AccAddressFromHex creates an AccAddress from a hex string.
func AccAddressFromHex(address string) (addr AccAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding Bech32 address failed: must provide an address")
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return AccAddress(bz), nil
}

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func AccAddressFromBech32(address string) (addr AccAddress, err error) {
	bech32PrefixAccAddr := GetConfig().GetBech32AccountAddrPrefix()
	bz, err := GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		return nil, err
	}

	return AccAddress(bz), nil
}

// Returns boolean for whether two AccAddresses are Equal
func (aa AccAddress) Equals(aa2 AccAddress) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Compare(aa.Bytes(), aa2.Bytes()) == 0
}

// Returns boolean for whether an AccAddress is empty
func (aa AccAddress) Empty() bool {
	if aa == nil {
		return true
	}

	aa2 := AccAddress{}
	return bytes.Compare(aa.Bytes(), aa2.Bytes()) == 0
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

// Bytes returns the raw address bytes.
func (aa AccAddress) Bytes() []byte {
	return aa
}

// String implements the Stringer interface.
func (aa AccAddress) String() string {
	bech32PrefixAccAddr := GetConfig().GetBech32AccountAddrPrefix()
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixAccAddr, aa.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (aa AccAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", aa.String())))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(aa))))
	}
}

// ----------------------------------------------------------------------------
// validator operator
// ----------------------------------------------------------------------------

// ValAddress defines a wrapper around bytes meant to present a validator's
// operator. When marshaled to a string or JSON, it uses Bech32.
type ValAddress []byte

// ValAddressFromHex creates a ValAddress from a hex string.
func ValAddressFromHex(address string) (addr ValAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding Bech32 address failed: must provide an address")
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return ValAddress(bz), nil
}

// ValAddressFromBech32 creates a ValAddress from a Bech32 string.
func ValAddressFromBech32(address string) (addr ValAddress, err error) {
	bech32PrefixValAddr := GetConfig().GetBech32ValidatorAddrPrefix()
	bz, err := GetFromBech32(address, bech32PrefixValAddr)
	if err != nil {
		return nil, err
	}

	return ValAddress(bz), nil
}

// Returns boolean for whether two ValAddresses are Equal
func (va ValAddress) Equals(va2 ValAddress) bool {
	if va.Empty() && va2.Empty() {
		return true
	}

	return bytes.Compare(va.Bytes(), va2.Bytes()) == 0
}

// Returns boolean for whether an AccAddress is empty
func (va ValAddress) Empty() bool {
	if va == nil {
		return true
	}

	va2 := ValAddress{}
	return bytes.Compare(va.Bytes(), va2.Bytes()) == 0
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
	bech32PrefixValAddr := GetConfig().GetBech32ValidatorAddrPrefix()
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
		s.Write([]byte(fmt.Sprintf("%s", va.String())))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", va)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(va))))
	}
}

// ----------------------------------------------------------------------------
// consensus node
// ----------------------------------------------------------------------------

// ConsAddress defines a wrapper around bytes meant to present a consensus node.
// When marshaled to a string or JSON, it uses Bech32.
type ConsAddress []byte

// ConsAddressFromHex creates a ConsAddress from a hex string.
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

// ConsAddressFromBech32 creates a ConsAddress from a Bech32 string.
func ConsAddressFromBech32(address string) (addr ConsAddress, err error) {
	bech32PrefixConsAddr := GetConfig().GetBech32ConsensusAddrPrefix()
	bz, err := GetFromBech32(address, bech32PrefixConsAddr)
	if err != nil {
		return nil, err
	}

	return ConsAddress(bz), nil
}

// get ConsAddress from pubkey
func GetConsAddress(pubkey crypto.PubKey) ConsAddress {
	return ConsAddress(pubkey.Address())
}

// Returns boolean for whether two ConsAddress are Equal
func (ca ConsAddress) Equals(ca2 ConsAddress) bool {
	if ca.Empty() && ca2.Empty() {
		return true
	}

	return bytes.Compare(ca.Bytes(), ca2.Bytes()) == 0
}

// Returns boolean for whether an ConsAddress is empty
func (ca ConsAddress) Empty() bool {
	if ca == nil {
		return true
	}

	ca2 := ConsAddress{}
	return bytes.Compare(ca.Bytes(), ca2.Bytes()) == 0
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (ca ConsAddress) Marshal() ([]byte, error) {
	return ca, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (ca *ConsAddress) Unmarshal(data []byte) error {
	*ca = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (ca ConsAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(ca.String())
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (ca *ConsAddress) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	ca2, err := ConsAddressFromBech32(s)
	if err != nil {
		return err
	}

	*ca = ca2
	return nil
}

// Bytes returns the raw address bytes.
func (ca ConsAddress) Bytes() []byte {
	return ca
}

// String implements the Stringer interface.
func (ca ConsAddress) String() string {
	bech32PrefixConsAddr := GetConfig().GetBech32ConsensusAddrPrefix()
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixConsAddr, ca.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (ca ConsAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", ca.String())))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", ca)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(ca))))
	}
}

// ----------------------------------------------------------------------------
// auxiliary
// ----------------------------------------------------------------------------

// Bech32ifyAccPub returns a Bech32 encoded string containing the
// Bech32PrefixAccPub prefix for a given account PubKey.
func Bech32ifyAccPub(pub crypto.PubKey) (string, error) {
	bech32PrefixAccPub := GetConfig().GetBech32AccountPubPrefix()
	return bech32.ConvertAndEncode(bech32PrefixAccPub, pub.Bytes())
}

// MustBech32ifyAccPub returns the result of Bech32ifyAccPub panicing on failure.
func MustBech32ifyAccPub(pub crypto.PubKey) string {
	enc, err := Bech32ifyAccPub(pub)
	if err != nil {
		panic(err)
	}

	return enc
}

// Bech32ifyValPub returns a Bech32 encoded string containing the
// Bech32PrefixValPub prefix for a given validator operator's PubKey.
func Bech32ifyValPub(pub crypto.PubKey) (string, error) {
	bech32PrefixValPub := GetConfig().GetBech32ValidatorPubPrefix()
	return bech32.ConvertAndEncode(bech32PrefixValPub, pub.Bytes())
}

// MustBech32ifyValPub returns the result of Bech32ifyValPub panicing on failure.
func MustBech32ifyValPub(pub crypto.PubKey) string {
	enc, err := Bech32ifyValPub(pub)
	if err != nil {
		panic(err)
	}

	return enc
}

// Bech32ifyConsPub returns a Bech32 encoded string containing the
// Bech32PrefixConsPub prefixfor a given consensus node's PubKey.
func Bech32ifyConsPub(pub crypto.PubKey) (string, error) {
	bech32PrefixConsPub := GetConfig().GetBech32ConsensusPubPrefix()
	return bech32.ConvertAndEncode(bech32PrefixConsPub, pub.Bytes())
}

// MustBech32ifyConsPub returns the result of Bech32ifyConsPub panicing on
// failure.
func MustBech32ifyConsPub(pub crypto.PubKey) string {
	enc, err := Bech32ifyConsPub(pub)
	if err != nil {
		panic(err)
	}

	return enc
}

// GetFromBech32 decodes a bytestring from a Bech32 encoded string.
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
