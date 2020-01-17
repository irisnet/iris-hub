package types

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/types"
	sdk "github.com/irisnet/irishub/types"
)

type BaseToken struct {
	Id              string           `json:"id"`
	Family          AssetFamily      `json:"family"`
	Source          AssetSource      `json:"source"`
	Gateway         string           `json:"gateway"`
	Symbol          string           `json:"symbol"`
	Name            string           `json:"name"`
	Decimal         uint8            `json:"decimal"`
	CanonicalSymbol string           `json:"canonical_symbol"`
	MinUnitAlias    string           `json:"min_unit_alias"`
	InitialSupply   types.Int        `json:"initial_supply"`
	MaxSupply       types.Int        `json:"max_supply"`
	Mintable        bool             `json:"mintable"`
	Owner           types.AccAddress `json:"owner"`
}

func NewBaseToken(symbol, name, minUnit string, decimal uint8, initialSupply types.Int, maxSupply types.Int, mintable bool, owner types.AccAddress) BaseToken {
	symbol = strings.ToLower(strings.TrimSpace(symbol))
	name = strings.TrimSpace(name)
	minUnit = strings.ToLower(strings.TrimSpace(minUnit))

	if maxSupply.IsZero() {
		if mintable {
			maxSupply = sdk.NewInt(int64(MaximumAssetMaxSupply))
		} else {
			maxSupply = initialSupply
		}
	}

	return BaseToken{
		Family:        FUNGIBLE,
		Source:        NATIVE,
		Symbol:        symbol,
		Name:          name,
		Decimal:       decimal,
		MinUnitAlias:  minUnit,
		InitialSupply: initialSupply,
		MaxSupply:     maxSupply,
		Mintable:      mintable,
		Owner:         owner,
	}
}

// FungibleToken
type FungibleToken struct {
	BaseToken `json:"base_token"`
}

func NewFungibleToken(symbol, name, minUnit string, decimal uint8, initialSupply, maxSupply types.Int, mintable bool, owner types.AccAddress) FungibleToken {
	token := FungibleToken{
		BaseToken: NewBaseToken(symbol, name, minUnit, decimal, initialSupply, maxSupply, mintable, owner),
	}
	token.Id = token.GetUniqueID()
	return token
}

func (ft FungibleToken) GetDecimal() uint8 {
	return ft.Decimal
}

func (ft FungibleToken) IsMintable() bool {
	return ft.Mintable
}

func (ft FungibleToken) GetOwner() types.AccAddress {
	return ft.Owner
}

func (ft FungibleToken) GetSource() AssetSource {
	return ft.Source
}

func (ft FungibleToken) GetSymbol() string {
	return ft.Symbol
}

func (ft FungibleToken) GetGateway() string {
	return ft.Gateway
}

func (ft FungibleToken) GetUniqueID() string {
	return strings.ToLower(ft.Symbol)
}

func (ft FungibleToken) GetDenom() string {
	denom, _ := sdk.GetCoinMinDenom(ft.GetUniqueID())
	return denom
}

func (ft FungibleToken) GetInitSupply() types.Int {
	return ft.InitialSupply
}

func (ft FungibleToken) GetCoinType() types.CoinType {
	units := make(types.Units, 2)
	units[0] = types.NewUnit(ft.GetUniqueID(), 0)
	units[1] = types.NewUnit(ft.GetDenom(), ft.Decimal)
	return types.CoinType{
		Name:    ft.GetUniqueID(), // UniqueID == Coin Name
		MinUnit: units[1],
		Units:   units,
		Desc:    ft.Name,
	}
}

// String implements fmt.Stringer
func (ft FungibleToken) String() string {
	ct := ft.GetCoinType()
	initSupply, _ := ct.Convert(types.NewCoin(ft.GetDenom(), ft.InitialSupply).String(), ft.GetUniqueID())
	maxSupply, _ := ct.Convert(types.NewCoin(ft.GetDenom(), ft.MaxSupply).String(), ft.GetUniqueID())
	owner := ""
	if !ft.Owner.Empty() {
		owner = ft.Owner.String()
	}

	return fmt.Sprintf(`FungibleToken %s:
  Name:              %s
  Symbol:            %s
  Scale:             %d
  MinUnit:           %s
  Initial Supply:    %s
  Max Supply:        %s
  Mintable:          %v
  Owner:             %s`,
		ft.GetUniqueID(), ft.Name, ft.Symbol,
		ft.Decimal, ft.MinUnitAlias, initSupply, maxSupply, ft.Mintable, owner)
}

type Tokens []FungibleToken

func (tokens Tokens) String() string {
	if len(tokens) == 0 {
		return "[]"
	}

	out := ""
	for _, token := range tokens {
		out += fmt.Sprintf("%v \n", token.String())
	}
	return out[:len(out)-1]
}

func (tokens Tokens) Validate() sdk.Error {
	if len(tokens) == 0 {
		return nil
	}

	for _, token := range tokens {
		exp := sdk.NewIntWithDecimal(1, int(token.Decimal))
		initialSupply := uint64(token.InitialSupply.Div(exp).Int64())
		maxSupply := uint64(token.MaxSupply.Div(exp).Int64())
		msg := NewMsgIssueToken(token.Symbol, token.MinUnitAlias, token.Name, token.Decimal, initialSupply, maxSupply, token.Mintable, token.Owner)
		if err := ValidateMsgIssueToken(&msg); err != nil {
			return err
		}
	}
	return nil
}

// -----------------------------

func GetTokenID(symbol string) string {
	return strings.ToLower(fmt.Sprintf("i.%s", strings.TrimSpace(symbol)))
}

//CheckSymbol
func CheckSymbol(symbol string) sdk.Error {
	if strings.Contains(strings.ToLower(symbol), sdk.Iris) {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("symbol can not contains : %s", sdk.Iris))
	}
	// check symbol
	if len(symbol) < MinimumAssetSymbolLen || len(symbol) > MaximumAssetSymbolLen {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid asset symbol: %s", symbol))
	}
	if !IsBeginWithAlpha(symbol) || !IsAlphaNumeric(symbol) {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid asset symbol: %s", symbol))
	}
	return nil
}

// CheckTokenID checks if the given token id is valid
func CheckTokenID(id string) sdk.Error {
	prefix, symbol := GetTokenIDParts(id)

	// check symbol
	if err := CheckSymbol(symbol); err != nil {
		return err
	}
	// check prefix
	if prefix != "i" {
		return ErrInvalidTokenID(DefaultCodespace, fmt.Sprintf("invalid token-id: %s", id))
	}

	return nil
}

// GetTokenIDParts returns the source prefix and symbol
func GetTokenIDParts(id string) (prefix string, symbol string) {
	parts := strings.Split(strings.ToLower(id), ".")
	prefix = parts[0]
	symbol = strings.Join(parts[1:], ".")
	return
}
