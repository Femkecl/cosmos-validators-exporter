package main

import (
	b64 "encoding/base64"
	"time"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
)

type ValidatorResponse struct {
	Validator Validator `json:"validator"`
}

type Validator struct {
	OperatorAddress   string               `json:"operator_address"`
	ConsensusPubkey   ConsensusPubkey      `json:"consensus_pubkey"`
	Jailed            bool                 `json:"jailed"`
	Status            string               `json:"status"`
	Tokens            string               `json:"tokens"`
	DelegatorShares   string               `json:"delegator_shares"`
	Description       ValidatorDescription `json:"description"`
	UnbondingHeight   string               `json:"unbonding_height"`
	UnbondingTime     time.Time            `json:"unbonding_time"`
	Commission        ValidatorCommission  `json:"commission"`
	MinSelfDelegation string               `json:"min_self_delegation"`
}

type ConsensusPubkey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type ValidatorDescription struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

type ValidatorCommission struct {
	CommissionRates ValidatorCommissionRates `json:"commission_rates"`
	UpdateTime      time.Time                `json:"update_time"`
}

type ValidatorCommissionRates struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
}

type ValidatorInfo struct {
	Address                  string
	Moniker                  string
	Identity                 string
	Website                  string
	SecurityContact          string
	Details                  string
	Tokens                   float64
	TokensUSD                float64
	Jailed                   bool
	Status                   string
	CommissionRate           float64
	CommissionMaxRate        float64
	CommissionMaxChangeRate  float64
	CommissionUpdateTime     time.Time
	UnbondingHeight          int64
	UnbondingTime            time.Time
	MinSelfDelegation        int64
	DelegatorsCount          int64
	SelfDelegation           Balance
	SelfDelegationUSD        float64
	Rank                     uint64
	TotalStake               float64
	Commission               []Balance
	CommissionUSD            float64
	SelfDelegationRewards    []Balance
	SelfDelegationRewardsUSD float64
	WalletBalance            []Balance
	WalletBalanceUSD         float64
	MissedBlocksCount        int64
	IsTombstoned             bool
	JailedUntil              time.Time
	StartHeight              int64
	IndexOffset              int64
}

func (key *ConsensusPubkey) GetValConsAddress(prefix string) (string, error) {
	encCfg := simapp.MakeTestEncodingConfig()
	interfaceRegistry := encCfg.InterfaceRegistry

	sDec, _ := b64.StdEncoding.DecodeString(key.Key)
	pk := codecTypes.Any{
		TypeUrl: key.Type,
		Value:   append([]byte{10, 32}, sDec...),
	}

	var pkProto cryptoTypes.PubKey
	if err := interfaceRegistry.UnpackAny(&pk, &pkProto); err != nil {
		return "", err
	}

	cosmosValCons := types.ConsAddress(pkProto.Address()).String()
	properValCons, err := ChangeBech32Prefix(cosmosValCons, prefix)
	if err != nil {
		return "", err
	}

	return properValCons, nil
}

func NewValidatorInfo(validator Validator) ValidatorInfo {
	return ValidatorInfo{
		Address:                 validator.OperatorAddress,
		Moniker:                 validator.Description.Moniker,
		Identity:                validator.Description.Identity,
		Website:                 validator.Description.Website,
		SecurityContact:         validator.Description.SecurityContact,
		Details:                 validator.Description.Details,
		Tokens:                  StrToFloat64(validator.Tokens),
		Jailed:                  validator.Jailed,
		Status:                  validator.Status,
		CommissionRate:          StrToFloat64(validator.Commission.CommissionRates.Rate),
		CommissionMaxRate:       StrToFloat64(validator.Commission.CommissionRates.MaxRate),
		CommissionMaxChangeRate: StrToFloat64(validator.Commission.CommissionRates.MaxChangeRate),
		CommissionUpdateTime:    validator.Commission.UpdateTime,
		UnbondingHeight:         StrToInt64(validator.UnbondingHeight),
		UnbondingTime:           validator.UnbondingTime,
		MinSelfDelegation:       StrToInt64(validator.MinSelfDelegation),
		MissedBlocksCount:       -1,
	}
}

type ValidatorQuery struct {
	Chain   string
	Address string
	Queries []QueryInfo
	Info    ValidatorInfo
}

type PaginationResponse struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Total string `json:"total"`
}

type QueryInfo struct {
	URL      string
	Duration time.Duration
	Success  bool
}

func (q *ValidatorQuery) GetSuccessfulQueriesCount() float64 {
	var count int64 = 0

	for _, query := range q.Queries {
		if query.Success {
			count++
		}
	}

	return float64(count)
}

func (q *ValidatorQuery) GetFailedQueriesCount() float64 {
	return float64(len(q.Queries)) - q.GetSuccessfulQueriesCount()
}

type ValidatorsResponse struct {
	Validators []Validator `json:"validators"`
}

type Balance struct {
	Amount float64
	Denom  string
}

type BalancesResponse struct {
	Balances types.Coins `json:"balances"`
}

type SigningInfoResponse struct {
	ValSigningInfo ValidatorSigningInfo `json:"val_signing_info"`
}

type ValidatorSigningInfo struct {
	Address             string    `json:"address"`
	StartHeight         string    `json:"start_height"`
	IndexOffset         string    `json:"index_offset"`
	JailedUntil         time.Time `json:"jailed_until"`
	Tombstoned          bool      `json:"tombstoned"`
	MissedBlocksCounter string    `json:"missed_blocks_counter"`
}
