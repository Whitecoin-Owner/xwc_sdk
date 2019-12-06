/**
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: operation.go.go, Date: 2018-10-31
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package xwc

type Asset struct {
	Xwc_amount   int64  `json:"amount"`
	Xwc_asset_id string `json:"asset_id"`
}

//
// xwc  --- "1.3.0"
// btc --- "1.3.1"
// ltc --- "1.3.2"
// hc  --- "1.3.3"
func DefaultAsset() Asset {
	return Asset{
		0,
		"1.3.0",
	}
}

type Extension struct {
	extension []string
}

type Memo struct {
	Xwc_from    string `json:"from"` //public_key_type  33
	Xwc_to      string `json:"to"`   //public_key_type  33
	Xwc_nonce   uint64 `json:"nonce"`
	Xwc_message string `json:"message"`

	IsEmpty bool   `json:"-"`
	Message string `json:"-"`
}

func DefaultMemo() Memo {

	return Memo{
		"XWC1111111111111111111111111111111114T1Anm",
		"XWC1111111111111111111111111111111114T1Anm",
		0,
		"",
		true,
		"",
	}

}

type Authority struct {
	Xwc_weight_threshold uint32          `json:"weight_threshold"`
	Xwc_account_auths    []interface{}   `json:"account_auths"`
	Xwc_key_auths        [][]interface{} `json:"key_auths"`
	Xwc_address_auths    []interface{}   `json:"address_auths"`

	Key_auths string `json:"-"`
}

func DefaultAuthority() Authority {

	return Authority{
		1,
		[]interface{}{},
		[][]interface{}{{"", 1}},
		[]interface{}{},
		"",
	}
}

type AccountOptions struct {
	Xwc_memo_key              string        `json:"memo_key"`
	Xwc_voting_account        string        `json:"voting_account"`
	Xwc_num_witness           uint16        `json:"num_witness"`
	Xwc_num_committee         uint16        `json:"num_committee"`
	Xwc_votes                 []interface{} `json:"votes"`
	Xwc_miner_pledge_pay_back byte          `json:"miner_pledge_pay_back"`
	Xwc_extensions            []interface{} `json:"extensions"`
}

func DefaultAccountOptions() AccountOptions {

	return AccountOptions{
		"",
		"1.2.5",
		0,
		0,
		[]interface{}{},
		10,
		[]interface{}{},
	}

}

// transfer operation tag is  0
type TransferOperation struct {
	Xwc_fee          Asset  `json:"fee"`
	Xwc_guarantee_id string `json:"guarantee_id,omitempty"`
	Xwc_from         string `json:"from"`
	Xwc_to           string `json:"to"`

	Xwc_from_addr string `json:"from_addr"`
	Xwc_to_addr   string `json:"to_addr"`

	Xwc_amount Asset `json:"amount"`
	Xwc_memo   *Memo `json:"memo,omitempty"`

	Xwc_extensions []interface{} `json:"extensions"`
}

func DefaultTransferOperation() *TransferOperation {

	return &TransferOperation{
		DefaultAsset(),
		"",
		"1.2.0",
		"1.2.0",
		"",
		"",
		DefaultAsset(),
		nil,
		make([]interface{}, 0),
	}
}

// account bind operation tag is 10
type AccountBindOperation struct {
	Xwc_fee               Asset  `json:"fee"`
	Xwc_crosschain_type   string `json:"crosschain_type"`
	Xwc_addr              string `json:"addr"`
	Xwc_account_signature string `json:"account_signature"`
	Xwc_tunnel_address    string `json:"tunnel_address"`
	Xwc_tunnel_signature  string `json:"tunnel_signature"`
	Xwc_guarantee_id      string `json:"guarantee_id,omitempty"`
}

func DefaultAccountBindOperation() *AccountBindOperation {

	return &AccountBindOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
		"",
	}
}

// account unbind operation tag is 11
type AccountUnBindOperation struct {
	Xwc_fee               Asset  `json:"fee"`
	Xwc_crosschain_type   string `json:"crosschain_type"`
	Xwc_addr              string `json:"addr"`
	Xwc_account_signature string `json:"account_signature"`
	Xwc_tunnel_address    string `json:"tunnel_address"`
	Xwc_tunnel_signature  string `json:"tunnel_signature"`
}

func DefaultAccountUnBindOperation() *AccountUnBindOperation {

	return &AccountUnBindOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
	}
}

// withdraw cross chain operation tag is 61
type WithdrawCrosschainOperation struct {
	Xwc_fee              Asset  `json:"fee"`
	Xwc_withdraw_account string `json:"withdraw_account"`
	Xwc_amount           string `json:"amount"`
	Xwc_asset_symbol     string `json:"asset_symbol"`

	Xwc_asset_id           string `json:"asset_id"`
	Xwc_crosschain_account string `json:"crosschain_account"`
	Xwc_memo               string `json:"memo"`
}

func DefaultWithdrawCrosschainOperation() *WithdrawCrosschainOperation {

	return &WithdrawCrosschainOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
		"",
	}
}

//register account operation tag is 5
type RegisterAccountOperation struct {
	Xwc_fee              Asset     `json:"fee"`
	Xwc_registrar        string    `json:"registrar"`
	Xwc_referrer         string    `json:"referrer"`
	Xwc_referrer_percent uint16    `json:"referrer_percent"`
	Xwc_name             string    `json:"name"`
	Xwc_owner            Authority `json:"owner"`
	Xwc_active           Authority `json:"active"`
	Xwc_payer            string    `json:"payer"`

	Xwc_options      AccountOptions `json:"options"`
	Xwc_extensions   interface{}    `json:"extensions"`
	Xwc_guarantee_id string         `json:"guarantee_id,omitempty"`
}

func DefaultRegisterAccountOperation() *RegisterAccountOperation {

	return &RegisterAccountOperation{
		DefaultAsset(),
		"1.2.0",
		"1.2.0",
		0,
		"",
		DefaultAuthority(),
		DefaultAuthority(),
		"",

		DefaultAccountOptions(),
		make(map[string]interface{}, 0),
		"",
	}

}

//lock balance operation tag is 55
type LockBalanceOperation struct {
	Xwc_lock_asset_id     string `json:"lock_asset_id"`
	Xwc_lock_asset_amount int64  `json:"lock_asset_amount"`
	Xwc_contract_addr     string `json:"contract_addr"`

	Xwc_lock_balance_account string `json:"lock_balance_account"`
	Xwc_lockto_miner_account string `json:"lockto_miner_account"`
	Xwc_lock_balance_addr    string `json:"lock_balance_addr"`

	Xwc_fee Asset `json:"fee"`
}

func DefaultLockBalanceOperation() *LockBalanceOperation {

	return &LockBalanceOperation{
		"1.3.0",
		0,
		"",
		"",
		"",
		"",
		DefaultAsset(),
	}
}

//foreclose balance operation tag is 56
type ForecloseBalanceOperation struct {
	Xwc_fee Asset `json:"fee"`

	Xwc_foreclose_asset_id     string `json:"foreclose_asset_id"`
	Xwc_foreclose_asset_amount int64  `json:"foreclose_asset_amount"`

	Xwc_foreclose_miner_account string `json:"foreclose_miner_account"`
	Xwc_foreclose_contract_addr string `json:"foreclose_contract_addr"`

	Xwc_foreclose_account string `json:"foreclose_account"`
	Xwc_foreclose_addr    string `json:"foreclose_addr"`
}

func DefaultForecloseBalanceOperation() *ForecloseBalanceOperation {

	return &ForecloseBalanceOperation{
		DefaultAsset(),
		"1.3.0",
		0,
		"",
		"",
		"",
		"",
	}
}

//obtain pay back operation tag is 73
type ObtainPaybackOperation struct {
	Xwc_pay_back_owner   string          `json:"pay_back_owner"`
	Xwc_pay_back_balance [][]interface{} `json:"pay_back_balance"`
	Xwc_guarantee_id     string          `json:"guarantee_id,omitempty"`
	Xwc_fee              Asset           `json:"fee"`

	citizen_name []string
	obtain_asset []Asset
}

func DefaultObtainPaybackOperation() *ObtainPaybackOperation {

	return &ObtainPaybackOperation{
		"",
		[][]interface{}{{"", DefaultAsset()}},
		"",
		DefaultAsset(),
		nil,
		nil,
	}
}

// contract invoke operation tag is 79
type ContractInvokeOperation struct {
	Xwc_fee           Asset  `json:"fee"`
	Xwc_invoke_cost   uint64 `json:"invoke_cost"`
	Xwc_gas_price     uint64 `json:"gas_price"`
	Xwc_caller_addr   string `json:"caller_addr"`
	Xwc_caller_pubkey string `json:"caller_pubkey"`
	Xwc_contract_id   string `json:"contract_id"`
	Xwc_contract_api  string `json:"contract_api"`
	Xwc_contract_arg  string `json:"contract_arg"`
	//Xwc_extension     []interface{} `json:"extensions"`
	Xwc_guarantee_id string `json:"guarantee_id,omitempty"`
}

func DefaultContractInvokeOperation() *ContractInvokeOperation {

	return &ContractInvokeOperation{
		DefaultAsset(),
		0,
		0,
		"",
		"",
		"",
		"",
		"",
		//make([]interface{}, 0),
		"",
	}
}

// transfer to contract operation tag is 81
type ContractTransferOperation struct {
	Xwc_fee           Asset  `json:"fee"`
	Xwc_invoke_cost   uint64 `json:"invoke_cost"`
	Xwc_gas_price     uint64 `json:"gas_price"`
	Xwc_caller_addr   string `json:"caller_addr"`
	Xwc_caller_pubkey string `json:"caller_pubkey"`
	Xwc_contract_id   string `json:"contract_id"`
	Xwc_amount        Asset  `json:"amount"`
	Xwc_param         string `json:"param"`
	//Xwc_extension     []interface{} `json:"extensions"`
	Xwc_guarantee_id string `json:"guarantee_id,omitempty"`
}

func DefaultContractTransferOperation() *ContractTransferOperation {

	return &ContractTransferOperation{
		DefaultAsset(),
		0,
		0,
		"",
		"",
		"",
		DefaultAsset(),
		"",
		//make([]interface{}, 0),
		"",
	}
}
