package xwc_sdk

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"github.com/Whitecoin-XWC/xwc_sdk/src/xwc_sdk/common"
	"github.com/Whitecoin-XWC/xwc_sdk/src/xwc_sdk/xwc"

	"math"
	ierr "github.com/Whitecoin-XWC/xwc_sdk/src/xwc_sdk/err"
)

// xwc API

const (
	XWCBind             = "bind"
	XWCUbind            = "ubind"
	XWCTransfer         = "transfer"
	XWCWithdraw         = "withdraw"
	XWCRegister         = "register"
	XWCMining           = "mining"
	XWCRewards          = "rewards"
	XwcRedeem           = "redeem"
	XWCContractInvoke   = "contractinvoke"
	XWCContractTransfer = "contracttransfer"
	minerInvalidAddress = "InvalidAddress"

	fieldCrossWif           = "crossWif"
	fieldFromAddr           = "fromAddr"
	fieldToAddr             = "toAddr"
	fieldAccountName        = "accountName"
	fieldOrigAddr           = "origAddr"
	fieldRefBlock           = "refBlock"
	fieldChainID            = "chainId"
	fieldPrecision          = "precision"
	fieldTranFee            = "tranFee"
	fieldTranAmt            = "tranAmt"
	fieldCoinAddr           = "coinAddr"
	fieldCoinType           = "coinType"
	fieldAccountId          = "accountId"
	fieldAssetId            = "assetId"
	fieldPayBackList        = "payBackList"
	fieldCitizenName        = "citizenName"
	fieldCitizenAmount      = "amount"
	fieldMainAssetPrecision = "mainCoinPrecision"
	fieldContractAPI        = "contractApi"
	fieldContractID         = "contractId"
	fieldBasicFee           = "basicFee"
	fieldGasPrice           = "gasPrice"
	fieldGasLimit           = "gasLimit"

	xwcPrecisionBits = 5
	xwcPrecision     = 100000
)

func getAssetId(coinType string) string {
	coinType = strings.ToUpper(coinType)
	switch coinType {
	case "XWC":
		return "1.3.0"
	case "BTC":
		return "1.3.1"
	case "LTC":
		return "1.3.2"
	case "HC":
		return "1.3.3"
	case "ETH":
		return "1.3.4"
	case "ERCPAX":
		return "1.3.5"
	case "ERCELF":
		return "1.3.6"
	case "USDT":
		return "1.3.7"
	case "BCH":
		return "1.3.8"
	case "ERCTITAN":
		return "1.3.9"
	default:
		return "1.3.999"

	}
}

// Guarantee guarantee for transaction fee
type Guarantee struct {
	// omit other field, such as coinType, ratio, guaranteeFee
	GuaranteeId string `json:"guaranteeId"`
}

func CalRefInfo(blockHash string) string {
	blockNum := blockHash[:8]
	ref_block_id := blockHash[8:16]
	fmt.Println(ref_block_id)
	var ref_block_num_little uint16
	blockNumBytes, _ := hex.DecodeString(blockNum)

	ref_block_num_little = binary.BigEndian.Uint16(blockNumBytes[2:4])
	ref_block_prefix, _ := hex.DecodeString(ref_block_id)

	end_ref_block_prefix := binary.LittleEndian.Uint32(ref_block_prefix)
	ref_str := fmt.Sprintf("%d,%d", ref_block_num_little, end_ref_block_prefix)

	return ref_str
}

func GetNewPrivate() (privWif string, pubWif string, addr string, err error) {
	return xwc.GetNewPrivate()
}

func validateXWCAddress(address, net string) (bool, error) {
	ok := xwc.ValidateAddress(address, net)
	return ok, nil
}

func dumpXWCPubKey(seed []byte, net string, account, index int) (string, error) {
	wif, err := xwc.ExportWif(seed, uint32(account), uint32(index))
	if err != nil {
		return "", ierr.ErrWrap(err, ierr.ErrDumpXWCKey)
	}
	pub, err := xwc.DerivePubkey(wif)
	if err != nil {
		return "", ierr.ErrWrap(err, ierr.ErrDumpPubKey)
	}
	return pub, nil
}

// dumpXWCPrivateKey dump XWC private key
func dumpXWCPrivateKey(seed []byte, net string, account, index int) (string, error) {
	key, err := xwc.ExportWif(seed, uint32(account), uint32(index))
	if err != nil {
		return key, ierr.ErrWrap(err, ierr.ErrDumpXWCKey)
	}
	return key, nil
}

func getRefChainID(dataJson map[string]interface{}) (ref string, cid string, err error) {
	var ok bool

	ref, ok = dataJson[fieldRefBlock].(string)
	if !ok {
		// bts response is refInfo
		ref, ok = dataJson["refInfo"].(string)
		if !ok {
			err = fmt.Errorf("field %s cannot convert to string", fieldRefBlock)
			return
		}
	}
	if ref == "" {
		err = fmt.Errorf("field %s is empty", fieldRefBlock)
		return
	}

	cid, ok = dataJson[fieldChainID].(string)
	if !ok {
		err = fmt.Errorf("field %s cannot convert to string", fieldChainID)
		return
	}
	if cid == "" {
		err = fmt.Errorf("field %s is empty", fieldChainID)
		return
	}
	return
}

// getAsset get asset info
// jmap1: dataJson
// jmap2: formJson
func getAsset(jmap1, jmap2 map[string]interface{}) (assetId string, amount, fee int64, err error) {
	coinType := getStringField(jmap2, "coinType")

	iasset := jmap2["asset"]
	iassets := jmap1["assets"]
	if iasset == nil {
		if iassets == nil {
			err = fmt.Errorf("not found asset in dataJson")
			return
		}
		assets := iassets.(map[string]interface{})
		iasset = assets[coinType]
	}
	asset := iasset.(map[string]interface{})
	//iprecision, ok := asset[fieldPrecision]
	//if !ok {
	//	err = fmt.Errorf("no found precission in asset")
	//	return
	//}
	precision := getPrecision(asset) // int64(iprecision.(float64))

	amount, err = getInt64(jmap2, fieldTranAmt, precision)
	if err != nil {
		return
	}
	assetId = asset["assetId"].(string)
	fee, err = getInt64(jmap1, fieldTranFee, xwcPrecision)
	if err != nil {
		return
	}
	return
}

//
//// XWCTransaction xwc transaction
//// action: bind
//func XWCTransaction(req *AddrReq, action string, data, form []byte) (buf []byte, err error) {
//	/*defer func() {
//		if r := recover(); r != nil {
//			err = ierr.ErrWrap(fmt.Errorf("%v", r), ierr.ErrXWCTransaction)
//			return
//		}
//	}()
//	*/
//	seed := xwc.MnemonicToSeed(req.Mnemonic, "")
//	wif, err := xwc.ExportWif(seed, uint32(req.Account), uint32(req.Index))
//	if err != nil {
//		err = ierr.ErrWrap(err, ierr.ErrXWCExportWif)
//		return
//	}
//
//	var (
//		sresp    ServerResp
//		dataJson map[string]interface{}
//		formJson map[string]interface{}
//	)
//
//	err = json.Unmarshal(data, &sresp)
//	if err != nil {
//		err = ierr.ErrWrap(err, ierr.ErrUnmarshalJSON)
//		return
//	}
//	dataJson = sresp.Data
//
//	err = json.Unmarshal(form, &formJson)
//	if err != nil {
//		err = ierr.ErrWrap(err, ierr.ErrUnmarshalJSON)
//		return
//	}
//
//	coinAddr := getStringField(formJson, fieldCoinAddr)
//	cointype := getStringField(formJson, fieldCoinType) // fieldCoinAddr bug? fixed!
//	if coinAddr != "" && cointype == CoinHC {
//		var vaddr *AddrResp
//		vaddr, err = getXWCAddress(req)
//		if err != nil {
//			return
//		}
//		if coinAddr != vaddr.Address {
//			err = ierr.ErrWrap(fmt.Errorf("coinAddr(%v) should equal with address derived from param req(%v)", coinAddr, vaddr), ierr.ErrAddressInvalid)
//			return
//		}
//	}
//
//	ref, cid, err := getRefChainID(dataJson)
//	if err != nil {
//		err = ierr.ErrWrap(err, ierr.ErrXWCRefChainID)
//		return
//	}
//
//	action = strings.ToLower(action)
//	fmt.Println("action:", action)
//	switch action {
//	case XWCBind:
//		buf, err = xwcBind(ref, wif, cid, dataJson, formJson)
//
//	case XWCUbind:
//		buf, err = xwcUnbind(ref, wif, cid, dataJson, formJson)
//
//	case XWCTransfer:
//		buf, err = xwcTransfer(ref, wif, cid, dataJson, formJson)
//
//	case XWCWithdraw:
//		buf, err = xwcWithdraw(ref, wif, cid, dataJson, formJson)
//
//	case XWCRegister:
//		buf, err = xwcRegister(ref, wif, cid, dataJson, formJson)
//
//	case XWCMining:
//		buf, err = xwcMining(ref, wif, cid, dataJson, formJson)
//
//	case XWCRewards:
//		buf, err = xwcRewards(ref, wif, cid, dataJson, formJson)
//
//	case XWCContractInvoke:
//		buf, err = xwcContractInvoke(ref, wif, cid, dataJson, formJson)
//
//	case XWCContractTransfer:
//		buf, err = xwcContractTransfer(ref, wif, cid, dataJson, formJson)
//
//	case XwcRedeem:
//		buf, err = xwcRedeem(ref, wif, cid, dataJson, formJson)
//
//	default:
//		err = ierr.ErrWrap(fmt.Errorf("unknown action: %s", action), ierr.ErrXWCInvalidAction)
//		return
//	}
//
//	if err != nil {
//		err = ierr.ErrWrap(err, ierr.ErrXWCTransaction)
//	}
//	return
//}
//
// convert app's param, hc -> HC, (xwc)hc -> HC
func convertXWCSymbol(s string) string {
	var (
		assets    map[string]string
		xwcAssets = map[string]string{}
	)

	s = strings.ToUpper(s)
	assets = map[string]string{
		xwc.CoinBTC:    xwc.CoinBTC,
		xwc.CoinUSDT:   xwc.CoinUSDT,
		xwc.CoinETH:    xwc.CoinETH,
		xwc.CoinHC:     xwc.CoinHC,
		xwc.CoinLTC:    xwc.CoinLTC,
		xwc.CoinPAX:    xwc.CoinERCPAX,
		xwc.CoinERCPAX: xwc.CoinERCPAX,
		xwc.CoinELF:    xwc.CoinERCELF,
		xwc.CoinERCELF: xwc.CoinERCELF,
		xwc.CoinBCH:    xwc.CoinBCH,
	}

	if v, ok := assets[s]; ok {
		return v
	}

	for k, v := range assets {
		xwcAssets[fmt.Sprintf("(XWC)%v", k)] = v
	}
	if v, ok := xwcAssets[s]; ok {
		return v
	}

	for k, v := range assets {
		xwcAssets[fmt.Sprintf("%v(XWC)", k)] = v
		// xwcAssets[fmt.Sprintf("")]
	}
	if v, ok := xwcAssets[s]; ok {
		return v
	}

	return s
}

//
//// getAmount get amount from map[string]interface{}, field can be string or number
//func getAmount(j map[string]interface{}, f string, p int64, def int64) int64 {
//	ip, ok := j[f]
//	if !ok {
//		fmt.Printf("not found field %s", f)
//		return def
//	}
//	fp, ok := ip.(float64)
//	if !ok {
//		// try to parse it as string
//		fmt.Printf("cannot convert field %v to float64, try string\n", f)
//		sp, ok := ip.(string)
//		if !ok {
//			fmt.Printf("field %v is neither float or string\n", f)
//			return def
//		}
//		ip, err := strconv.ParseFloat(sp, 64)
//		if err != nil {
//			fmt.Printf("cannot parse field %v\n", f)
//			return def
//		}
//		fp = float64(ip)
//	}
//	return int64(float64(p) * fp)
//}

func getMainAssetPrecision(dataJson map[string]interface{}) int64 {
	ip, ok := dataJson[fieldMainAssetPrecision]
	if !ok {
		fmt.Println("not found field main coin precision")
		return xwcPrecision
	}
	fp, ok := ip.(float64)
	if !ok {
		// try to parse it as string
		fmt.Println("cannot convert main coin precision to float64, try string")
		sp, ok := ip.(string)
		if !ok {
			fmt.Println("field main coin precision is neither float or string")
			return xwcPrecision
		}
		ip, err := strconv.ParseInt(sp, 10, 64)
		if err != nil {
			fmt.Println("cannot parse field main coin Precision")
			return xwcPrecision
		}
		fp = float64(ip)
	}

	p := int64(fp)
	precision := int64(1)
	for p > 0 {
		precision *= 10
		p = p - 1
	}
	return precision
}

// 5 -> 100000
// 8 -> 100000000
func getPrecision(dataJson map[string]interface{}) int64 {
	return _getPricision(dataJson, fieldPrecision)
}

func getPrecisionBits(dataJson map[string]interface{}) int64 {
	ip, ok := dataJson[fieldPrecision]
	if !ok {
		fmt.Println("not found field precision")
		return xwcPrecisionBits
	}
	fp, ok := ip.(float64)
	if !ok {
		// try to parse it as string
		fmt.Println("cannot convert precision to float64, try string")
		sp, ok := ip.(string)
		if !ok {
			fmt.Println("field precision is neither float or string")
			return xwcPrecisionBits
		}
		ip, err := strconv.ParseInt(sp, 10, 64)
		if err != nil {
			fmt.Println("cannot parse field Precision")
			return xwcPrecisionBits
		}
		fp = float64(ip)
	}

	return int64(fp)
}

func _getPricision(dataJson map[string]interface{}, name string) int64 {
	ip, ok := dataJson[name]
	if !ok {
		fmt.Println("not found field precision")
		return xwcPrecision
	}
	fp, ok := ip.(float64)
	if !ok {
		// try to parse it as string
		fmt.Println("cannot convert precision to float64, try string")
		sp, ok := ip.(string)
		if !ok {
			fmt.Println("field precision is neither float or string")
			return xwcPrecision
		}
		ip, err := strconv.ParseInt(sp, 10, 64)
		if err != nil {
			fmt.Println("cannot parse field Precision")
			return xwcPrecision
		}
		fp = float64(ip)
	}

	p := int64(fp)
	precision := int64(1)
	for p > 0 {
		precision *= 10
		p = p - 1
	}
	return precision
}

func getGuaranteeID(dataJson map[string]interface{}) string {
	ig, ok := dataJson["guarantee"]
	if !ok {
		return ""
	}

	g, ok := ig.(map[string]interface{})
	if !ok {
		return ""
	}
	gid, ok := g["guaranteeId"]
	if !ok {
		return ""
	}
	sid, ok := gid.(string)
	if !ok {
		return ""
	}
	return sid
}

// get citizen info from formJson
func getMiner(jmap map[string]interface{}) (id, addr string, err error) {
	icitizen, ok := jmap["citizen"]
	if !ok {
		err = fmt.Errorf("not found field citizen")
		return
	}
	citizen, ok := icitizen.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("cannot convert citizen to map[string]interface{}")
		return
	}
	iid, ok := citizen["citizenId"]
	if !ok {
		err = fmt.Errorf("not found field citizenId")
		return
	}
	iaddr, ok := citizen["citizenAddress"]
	if !ok {
		iaddr = ""
		// err = fmt.Errorf("not found field citizenAddress")
		// return
	}
	id = iid.(string)
	addr = iaddr.(string)
	return
}

// convert string to float64, multiple precision
func getInt64(jmap map[string]interface{}, field string, precision int64) (int64, error) {
	is, ok := jmap[field]
	if !ok {
		return 0, fmt.Errorf("not found field param %v", field)
	}
	s, ok := is.(string)
	if !ok {
		return 0, fmt.Errorf("field %v is invalid format", field)
	}

	ii, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(ii * float64(precision))), nil
}
func getXwcInt64(value string) (int64, error) {
	ii, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(ii * float64(xwcPrecision))), nil
}

func getIntField(jmap map[string]interface{}, field string) (int64, error) {
	is, ok := jmap[field]
	if !ok {
		return 0, fmt.Errorf("not found field param %v", field)
	}
	fi, ok := is.(float64)
	if !ok {
		return 0, fmt.Errorf("convert param %s failed", field)
	}
	return int64(fi), nil
}

func getStringField(jmap map[string]interface{}, field string) string {
	is, ok := jmap[field]
	if !ok {
		fmt.Printf("not found field %v in json map %v\n", field, jmap)
		return ""
	}
	s := is.(string)
	return s
}

func getStringFieldWithDefault(jmap map[string]interface{}, field, def string) string {
	is, ok := jmap[field]
	if !ok {
		fmt.Printf("not found field %v in json map %v, return default %v\n", field, jmap, def)
		return def
	}
	s := is.(string)
	return s
}

// xwcBind bind tunnel address
/*
	form:
	{
		origAddr string	绑定原地址，例如 HC 地址
		coinAddr string	地址，例如HC(XWC)地址
		coinType string	绑定币种，例如HC(XWC)

		// 下面的字段是需要增加的：
		crosswif string  // 这个wif是绑定币种(hc, btc, ltc, eth)的 wif
	}
*/
func XwcBind(ref, wif, chainId, originAddr, coinAddr, coinType, crossWif string) (buf []byte, err error) {
	if coinType == xwc.CoinBCH && strings.HasPrefix(originAddr, "bitcoincash:") == false {
		originAddr = "bitcoincash:" + originAddr
	}
	cwif := crossWif

	return xwc.BuildBindAccountTransaction(ref, wif, coinAddr, 0, originAddr, coinType, cwif, "", chainId)
}
func xwcBind(ref, wif, cid string, dataJson, formJson map[string]interface{}) (buf []byte, err error) {
	origAddr := getStringField(formJson, fieldOrigAddr) // formJson[fieldOrigAddr].(string)
	coinAddr := getStringField(formJson, fieldCoinAddr)
	oct := getStringField(formJson, fieldCoinType)
	if oct == "" {
		err = fmt.Errorf("coinType is empty")
		return
	}
	coinType := convertXWCSymbol(oct)
	if coinType == xwc.CoinBCH && strings.HasPrefix(origAddr, "bitcoincash:") == false {
		origAddr = "bitcoincash:" + origAddr
	}
	cwif := getStringField(formJson, fieldCrossWif)

	return xwc.BuildBindAccountTransaction(ref, wif, coinAddr, 0, origAddr, coinType, cwif, "", cid)
}

func XwcUnbind(ref, wif, chainId, origAddr, coinAddr, coinType, crossWif string) (buf []byte, err error) {
	if coinType == xwc.CoinBCH && strings.HasPrefix(origAddr, "bitcoincash:") == false {
		origAddr = "bitcoincash:" + origAddr
	}
	return xwc.BuildUnBindAccountTransaction(ref, wif, coinAddr, 20000, origAddr, coinType, crossWif, chainId)
}
func xwcUnbind(ref, wif, cid string, dataJson, formJson map[string]interface{}) (buf []byte, err error) {
	origAddr := getStringField(formJson, fieldOrigAddr)
	coinAddr := getStringField(formJson, fieldCoinAddr)
	coinType := convertXWCSymbol(getStringField(formJson, fieldCoinType))
	cwif := getStringField(formJson, fieldCrossWif)

	if coinType == xwc.CoinBCH && strings.HasPrefix(origAddr, "bitcoincash:") == false {
		origAddr = "bitcoincash:" + origAddr
	}

	// todo get fee from dataJson
	return xwc.BuildUnBindAccountTransaction(ref, wif, coinAddr, 20000, origAddr, coinType, cwif, cid)
}

func XwcTransfer(ref, wif, chainId, fromAddr, toAddr, coinType string, transferAmount, fee string, memo string, guaranteeId string) (buf []byte, err error) {
	tranAmt, err := getXwcInt64(transferAmount)
	if err != nil {
		return
	}
	// trans fee is calc with main coin precision
	tranFee, err := getXwcInt64(fee)
	if err != nil {
		return
	}
	gid := guaranteeId

	assetId := getAssetId(coinType)

	return xwc.BuildTransferTransaction(ref, wif, fromAddr, toAddr, memo, assetId, tranAmt, tranFee, coinType, gid, chainId)

}

/*
formJson:
{
    fromAddr      string	交易来源地址
    toAddr	　	string	交易目标地址
    coinType		string	质押币种
    tranAmt		number	转账金额
    tranFee		number	手续费
    bak		number	转账备注(20个字符)
    guarantee	[{}	JSONObject	手续费承兑单(没有为null)
        coinType	string	货币类型
        ratio	number	兑换比例
        guaranteeFee	number	支付手续费的实际金额
        guaranteeId	string	手续费承兑单编号
    ]
}

datajson precision is hc precision
*/
func xwcTransfer(ref, wif, cid string, dataJson, formJson map[string]interface{}) (buf []byte, err error) {
	var sbak string

	from := getStringField(formJson, fieldFromAddr)                       // ["fromAddr"].(string)
	to := getStringField(formJson, fieldToAddr)                           // ["toAddr"].(string)
	coinType := convertXWCSymbol(getStringField(dataJson, fieldCoinType)) // ["coinType"].(string))
	precision := getPrecision(dataJson)                                   // int64(dataJson[fieldPrecision].(float64))
	mainAssetPrecision := getMainAssetPrecision(dataJson)

	bak := formJson["bak"]
	if bak != nil {
		sbak = bak.(string)
	}

	tranAmt, err := getInt64(formJson, fieldTranAmt, precision)
	if err != nil {
		return
	}
	// trans fee is calc with main coin precision
	tranFee, err := getInt64(formJson, fieldTranFee, mainAssetPrecision)
	if err != nil {
		return
	}
	gid := getGuaranteeID(formJson)

	assetId := getStringField(dataJson, fieldAssetId)

	return xwc.BuildTransferTransaction(ref, wif, from, to, sbak, assetId, tranAmt, tranFee, coinType, gid, cid)
}

func XwcWithdraw(ref, wif, chainId, fromAddr, toAddr, coinType, transferAmt, fee string) (buf []byte, err error) {
	//tranAmt,err := getXwcInt64(transferAmt)
	tranFee, err := getXwcInt64(fee)
	if err != nil {
		return
	}
	assetId := getAssetId(coinType)

	return xwc.BuildWithdrawCrosschainTransaction(ref, wif, fromAddr, tranFee, toAddr, coinType, assetId, transferAmt, chainId)
}

/*
formJson:
{
    fromAddr	　	string	发起方货币地址
    toAddr		string	目标地址
    coinType	　	string	提现币种
    tranAmt		string	提现金额
    tranFee		string	提现手续费
}
*/
func xwcWithdraw(ref, wif, cid string, dataJson, formJson map[string]interface{}) (buf []byte, err error) {
	from := getStringField(formJson, fieldFromAddr) //["fromAddr"].(string)
	to := getStringField(formJson, fieldToAddr)     // ["toAddr"].(string)
	coinType := convertXWCSymbol(getStringField(dataJson, fieldCoinType))

	tranAmt := getStringField(formJson, fieldTranAmt)
	precision := getPrecision(dataJson) // int64(dataJson[fieldPrecision].(float64))

	tranFee, err := getInt64(dataJson, fieldTranFee, precision)
	if err != nil {
		return
	}
	assetId := getStringField(dataJson, fieldAssetId)

	return xwc.BuildWithdrawCrosschainTransaction(ref, wif, from, tranFee, to, coinType, assetId, tranAmt, cid)
}

func XwcRegister(ref, wif, chainId, accountName, accountAddr, fee, guaranteeId string) (buf []byte, err error) {
	fee_value, err := getXwcInt64(fee)
	if err != nil {
		return
	}
	pubkey, err := xwc.DerivePubkey(wif)
	if err != nil {
		return
	}

	return xwc.BuildRegisterAccountTransaction(ref, wif, accountAddr, pubkey, fee_value, guaranteeId, accountName, chainId)
}

/*
formJson:
{
    accountName		string	注册账户名
    coinAddr		string	XWC地址
    coinType		string	货币类型，XWC
    tranFee		number	手续费
    guarantee	{	JSONObject	手续费承兑单(没有为null)
        coinType	string	货币类型
        ratio	number	兑换比例
        guaranteeFee	number	支付手续费的实际金额
        guaranteeId	string	手续费承兑单编号
	}
}
*/
func xwcRegister(ref, wif, cid string, dataJson, formJson map[string]interface{}) (buf []byte, err error) {
	precision := getPrecision(dataJson) // int64(dataJson[fieldPrecision].(float64))
	fee, err := getInt64(dataJson, fieldTranFee, precision)
	if err != nil {
		return
	}

	pubkey, err := xwc.DerivePubkey(wif)
	if err != nil {
		return
	}

	coinAddr := getStringField(formJson, fieldCoinAddr) // ["coinAddr"].(string)
	name := getStringField(formJson, fieldAccountName)  // ["accountName"].(string)
	gid := getGuaranteeID(formJson)
	return xwc.BuildRegisterAccountTransaction(ref, wif, coinAddr, pubkey, fee, gid, name, cid)
}

func XwcMining(ref, wif, chainId, coinType, accountId, addr, miningAmount, fee, minerId string) (buf []byte, err error) {
	tranAmt, err := getXwcInt64(miningAmount)

	if err != nil {
		return
	}
	feeValue, err := getXwcInt64(fee)

	if err != nil {
		return
	}
	assetId := getAssetId(coinType)
	// _ = minerAddr
	return xwc.BuildLockBalanceTransaction(ref, wif, addr, accountId, assetId, tranAmt, feeValue, minerId, minerInvalidAddress, chainId)
}

// xwcMining xwc mining
func xwcMining(ref, wif, cid string, dataJson, formJson map[string]interface{}) (buf []byte, err error) {
	assetId, tranAmt, fee, err := getAsset(dataJson, formJson)
	if err != nil {
		return
	}
	accountId := getStringField(dataJson, fieldAccountId)
	addr := getStringField(formJson, fieldCoinAddr)

	minerID, _, err := getMiner(formJson)
	if err != nil {
		return
	}
	// _ = minerAddr
	return xwc.BuildLockBalanceTransaction(ref, wif, addr, accountId, assetId, tranAmt, fee, minerID, minerInvalidAddress, cid)
}

func getFloatOrString(jm map[string]interface{}, name string) string {
	ii, ok := jm[name]
	if !ok || ii == nil {
		return ""
	}

	if s, ok := ii.(string); ok {
		return s
	}
	f, ok := ii.(float64)
	if ok {
		return fmt.Sprint(f)
	}
	return ""
}

//
func getPayList(formJson map[string]interface{}) (res []string, err error) {
	ips := formJson[fieldPayBackList]
	if ips == nil {
		err = fmt.Errorf("not found field %v", fieldPayBackList)
		return
	}

	pss, ok := ips.([]interface{})
	if !ok {
		err = fmt.Errorf("cannot convert field %v", fieldPayBackList)
		return
	}
	for _, item := range pss {
		ps := item.(map[string]interface{})
		name := ps[fieldCitizenName].(string)
		amt := getFloatOrString(ps, fieldCitizenAmount)
		/*
			iamt, ierr := getInt64(ps, fieldCitizenAmount, precision)
			if ierr != nil {
				err = ierr
				return
			}
		*/

		// assetId
		assetId := getStringField(ps, fieldAssetId)
		// 分红的资产，这里暂时只能是 xwc
		res = append(res, fmt.Sprintf("%s,%v,%s", name, amt, assetId))
	}

	return
}

func XwcRewards(ref, wif, chainId, accountAddr string, payList []string, fee string, guaranteeId string) (buf []byte, err error) {
	feeAmt, err := getXwcInt64(fee)
	if err != nil {
		return
	}
	return xwc.BuildObtainPaybackTransaction(ref, wif, accountAddr, feeAmt, payList, guaranteeId, chainId)
}

// xwcRewards xwc reward
func xwcRewards(ref, wif, cid string, dataJson, formJson map[string]interface{}) (buf []byte, err error) {
	precision := getPrecision(dataJson)

	ps, err := getPayList(dataJson)
	if err != nil {
		return
	}
	addr := getStringField(formJson, fieldCoinAddr)
	//precision, err := getInt64(dataJson, fieldPrecision, xwcPrecision)
	//if err != nil {
	//	return
	//}
	tranFee, err := getInt64(dataJson, fieldTranFee, precision)
	if err != nil {
		return
	}
	gid := getGuaranteeID(formJson)
	return xwc.BuildObtainPaybackTransaction(ref, wif, addr, tranFee, ps, gid, cid)
}

func XwcForecloseBalance(ref, wif, chainId, accountAddr, accountId, assetId, minerId, amount, fee string) (buf []byte, err error) {
	tranAmt, err := getXwcInt64(amount)
	if err != nil {
		return
	}
	tranFee, err := getXwcInt64(fee)
	if err != nil {
		return
	}
	return xwc.BuildRedeemBalanceTransaction(ref, wif, accountAddr, accountId, assetId, tranAmt, tranFee, minerId, minerInvalidAddress, chainId)
}

// xwcRedeem xwc redeem
func xwcRedeem(ref, wif, cid string, dataJson, formJson map[string]interface{}) (buf []byte, err error) {
	accountID := getStringField(dataJson, fieldAccountId)
	precision := getPrecision(dataJson) // int64(dataJson[fieldPrecision].(float64))
	assetId := getStringField(dataJson, fieldAssetId)

	addr := getStringField(formJson, fieldCoinAddr)
	minerID, _, err := getMiner(formJson)
	if err != nil {
		return
	}
	tranAmt, err := getInt64(formJson, fieldTranAmt, precision)
	if err != nil {
		return
	}
	tranFee, err := getInt64(dataJson, fieldTranFee, xwcPrecision)
	if err != nil {
		return
	}

	return xwc.BuildRedeemBalanceTransaction(ref, wif, addr, accountID, assetId, tranAmt, tranFee, minerID, minerInvalidAddress, cid)
}

func XwcContractInvoke(ref, wif, chainId, accountAddr, contractAddr, contractMethod, contractArgs, fee, gasPrice, gasLimit, guaranteeId string) (buf []byte, err error) {
	feeAmt, err := getXwcInt64(fee)
	if err != nil {
		return nil, err
	}
	gasPriceAmt, err := getXwcInt64(gasPrice)
	if err != nil {
		return nil, err
	}
	gasLimitAmt, err := getXwcInt64(gasLimit)
	if err != nil {
		return nil, err
	}
	feeAmt += gasPriceAmt * gasLimitAmt / 100
	return xwc.BuildContractInvokeTransaction(ref, wif, accountAddr, feeAmt, gasPriceAmt, gasLimitAmt, contractAddr, contractMethod, contractArgs, guaranteeId, chainId)
}

func xwcContractInvoke(ref, wif, cid string, dataJson, formJson map[string]interface{}) (buf []byte, err error) {
	addr := getStringField(formJson, fieldFromAddr)
	gid := getGuaranteeID(formJson)

	// fee, gas price, gas limit
	fee, err := getIntField(dataJson, fieldBasicFee)
	if err != nil {
		return nil, err
	}
	gasPrice, err := getIntField(dataJson, fieldGasPrice)
	if err != nil {
		return nil, err
	}
	gasLimit, err := getIntField(dataJson, fieldGasLimit)
	if err != nil {
		return nil, err
	}

	fee += gasPrice * gasLimit / 100

	// contract api, id
	contractAPI := getStringField(dataJson, fieldContractAPI)
	contractId := getStringField(dataJson, fieldContractID)
	// dest addr, amount, memo
	precisionBits := getPrecisionBits(dataJson)
	toAddr := getStringField(formJson, fieldToAddr)
	samt := getStringField(formJson, fieldTranAmt)
	amount, err := common.ConvertToStringWithPrecision(samt, precisionBits)
	// amount, err := getInt64(formJson, fieldTranAmt, precision)
	if err != nil {
		return
	}
	bak := getStringField(dataJson, "bak") // TODO: 有的bak在from中
	if bak == "" {
		bak = getStringField(formJson, "bak")
	}
	contractArgs := toAddr + "," + fmt.Sprint(amount)
	if bak != "" {
		contractArgs += "," + bak
	}

	return xwc.BuildContractInvokeTransaction(ref, wif, addr, fee, gasPrice, gasLimit, contractId, contractAPI, contractArgs, gid, cid)
}
func XwcTransferToContract(ref, wif, chainId, accountAddr, contractAddr, amount, assetId, memo, fee, gasPrice, gasLimit, guaranteeId string) (buf []byte, err error) {
	// fee, gas price, gas limit
	feeAmt, err := getXwcInt64(fee)
	if err != nil {
		return nil, err
	}
	gasPriceAmt, err := getXwcInt64(gasPrice)
	if err != nil {
		return nil, err
	}
	gasLimitAmt, err := getXwcInt64(gasLimit)
	if err != nil {
		return nil, err
	}
	amountAmt, err := getXwcInt64(amount)
	if err != nil {
		return nil, err
	}
	feeAmt += gasPriceAmt * gasLimitAmt / 100

	return xwc.BuildContractTransferTransaction(ref, wif, accountAddr, feeAmt, amountAmt, assetId, gasPriceAmt, gasLimitAmt, contractAddr, memo, guaranteeId, chainId)
}
func xwcContractTransfer(ref, wif, cid string, dataJson, formJson map[string]interface{}) (buf []byte, err error) {
	gid := getGuaranteeID(formJson)
	addr := getStringField(formJson, fieldFromAddr)

	// fee, gas price, gas limit
	fee, err := getIntField(dataJson, fieldBasicFee)
	if err != nil {
		return nil, err
	}
	gasPrice, err := getIntField(dataJson, fieldGasPrice)
	if err != nil {
		return nil, err
	}
	gasLimit, err := getIntField(dataJson, fieldGasLimit)
	if err != nil {
		return nil, err
	}

	fee += gasPrice * gasLimit / 100

	// contract id
	contractId := getStringField(dataJson, fieldContractID)
	// TODO: 这里没加上转账到合约的备注，不能直接用''
	// dest addr, amount, memo
	precision := getPrecision(dataJson)
	// toAddr := getStringField(formJson, fieldToAddr)
	amount, err := getInt64(formJson, fieldTranAmt, precision)
	if err != nil {
		return nil, err
	}
	assetID := getStringField(formJson, fieldAssetId)

	memo := getStringField(dataJson, "memo")
	if memo == "" {
		memo = getStringField(formJson, "memo")
	}

	return xwc.BuildContractTransferTransaction(ref, wif, addr, fee, amount, assetID, gasPrice, gasLimit, contractId, memo, gid, cid)
}

func convertString(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func convertStringToInt64WithPrecision(s string, precision int64) (int64, error) {
	ii, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(ii * float64(precision))), nil
}

func convertPrecisionBits(bit int) int64 {
	precision := int64(1)
	for bit > 0 {
		precision *= 10
		bit = bit - 1
	}
	return precision
}
