package tx

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/tendermint/tendermint/crypto"

	codectypes "plugchain-sdk-go/codec/types"
	"plugchain-sdk-go/types"
	sdk "plugchain-sdk-go/types"
	"plugchain-sdk-go/types/tx/signing"
)

// wrapper is a wrapper around the tx.Tx proto.Message which retain the raw
// body and auth_info bytes.
type wrapper struct {
	tx *Tx

	// bodyBz represents the protobuf encoding of TxBody. This should be encoding
	// from the client using TxRaw if the tx was decoded from the wire
	bodyBz []byte

	// authInfoBz represents the protobuf encoding of TxBody. This should be encoding
	// from the client using TxRaw if the tx was decoded from the wire
	authInfoBz []byte

	txBodyHasUnknownNonCriticals bool
}

var (
	_ ExtensionOptionsTxBuilder = &wrapper{}
	_ codectypes.IntoAny        = &wrapper{}
)

// ExtensionOptionsTxBuilder defines a TxBuilder that can also set extensions.
type ExtensionOptionsTxBuilder interface {
	SetExtensionOptions(...*codectypes.Any)
	SetNonCriticalExtensionOptions(...*codectypes.Any)
}

func newBuilder() *wrapper {
	return &wrapper{
		tx: &Tx{
			Body: &TxBody{},
			AuthInfo: &AuthInfo{
				Fee: &Fee{},
			},
		},
	}
}

func (w *wrapper) GetMsgs() []sdk.Msg {
	return w.tx.GetMsgs()
}

func (w *wrapper) ValidateBasic() error {
	return w.tx.ValidateBasic()
}

func (w *wrapper) getBodyBytes() []byte {
	if len(w.bodyBz) == 0 {
		// if bodyBz is empty, then marshal the body. bodyBz will generally
		// be set to nil whenever SetBody is called so the result of calling
		// this method should always return the correct bytes. Note that after
		// decoding bodyBz is derived from TxRaw so that it matches what was
		// transmitted over the wire
		var err error
		w.bodyBz, err = proto.Marshal(w.tx.Body)
		if err != nil {
			panic(err)
		}
	}
	return w.bodyBz
}

func (w *wrapper) getAuthInfoBytes() []byte {
	if len(w.authInfoBz) == 0 {
		// if authInfoBz is empty, then marshal the body. authInfoBz will generally
		// be set to nil whenever SetAuthInfo is called so the result of calling
		// this method should always return the correct bytes. Note that after
		// decoding authInfoBz is derived from TxRaw so that it matches what was
		// transmitted over the wire
		var err error
		w.authInfoBz, err = proto.Marshal(w.tx.AuthInfo)
		if err != nil {
			panic(err)
		}
	}
	return w.authInfoBz
}

func (w *wrapper) GetSigners() []sdk.AccAddress {
	return w.tx.GetSigners()
}

func (w *wrapper) GetPubKeys() []crypto.PubKey {
	signerInfos := w.tx.AuthInfo.SignerInfos
	pks := make([]crypto.PubKey, len(signerInfos))

	for i, si := range signerInfos {
		// NOTE: it is okay to leave this nil if there is no PubKey in the SignerInfo.
		// PubKey's can be left unset in SignerInfo.
		if si.PublicKey == nil {
			continue
		}

		pk, ok := si.PublicKey.GetCachedValue().(crypto.PubKey)
		if ok {
			pks[i] = pk
		}
	}

	return pks
}

func (w *wrapper) GetGas() uint64 {
	return w.tx.AuthInfo.Fee.GasLimit
}

func (w *wrapper) GetFee() sdk.Coins {
	return w.tx.AuthInfo.Fee.Amount
}

func (w *wrapper) FeePayer() sdk.AccAddress {
	feePayer := w.tx.AuthInfo.Fee.Payer
	if feePayer != "" {
		payerAddr, err := sdk.AccAddressFromBech32(feePayer)
		if err != nil {
			panic(err)
		}
		return payerAddr
	}
	// use first signer as default if no payer specified
	return w.GetSigners()[0]
}

func (w *wrapper) FeeGranter() sdk.AccAddress {
	feePayer := w.tx.AuthInfo.Fee.Granter
	if feePayer != "" {
		granterAddr, err := sdk.AccAddressFromBech32(feePayer)
		if err != nil {
			panic(err)
		}
		return granterAddr
	}
	return nil
}

func (w *wrapper) GetMemo() string {
	return w.tx.Body.Memo
}

func (w *wrapper) GetSignatures() [][]byte {
	return w.tx.Signatures
}

// GetTimeoutHeight returns the transaction's timeout height (if set).
func (w *wrapper) GetTimeoutHeight() uint64 {
	return w.tx.Body.TimeoutHeight
}

func (w *wrapper) GetSignaturesV2() ([]signing.SignatureV2, error) {
	signerInfos := w.tx.AuthInfo.SignerInfos
	sigs := w.tx.Signatures
	pubKeys := w.GetPubKeys()
	n := len(signerInfos)
	res := make([]signing.SignatureV2, n)

	for i, si := range signerInfos {
		// handle nil signatures (in case of simulation)
		if si.ModeInfo == nil {
			res[i] = signing.SignatureV2{
				PubKey: pubKeys[i],
			}
		} else {
			var err error
			sigData, err := ModeInfoAndSigToSignatureData(si.ModeInfo, sigs[i])
			if err != nil {
				return nil, err
			}
			res[i] = signing.SignatureV2{
				PubKey:   pubKeys[i],
				Data:     sigData,
				Sequence: si.GetSequence(),
			}

		}
	}

	return res, nil
}

func (w *wrapper) SetMsgs(msgs ...sdk.Msg) error {
	anys := make([]*codectypes.Any, len(msgs))

	for i, msg := range msgs {
		var err error
		anys[i], err = codectypes.NewAnyWithValue(msg)
		if err != nil {
			return err
		}
	}

	w.tx.Body.Messages = anys

	// set bodyBz to nil because the cached bodyBz no longer matches tx.Body
	w.bodyBz = nil

	return nil
}

// SetTimeoutHeight sets the transaction's height timeout.
func (w *wrapper) SetTimeoutHeight(height uint64) {
	w.tx.Body.TimeoutHeight = height

	// set bodyBz to nil because the cached bodyBz no longer matches tx.Body
	w.bodyBz = nil
}

func (w *wrapper) SetMemo(memo string) {
	w.tx.Body.Memo = memo

	// set bodyBz to nil because the cached bodyBz no longer matches tx.Body
	w.bodyBz = nil
}

func (w *wrapper) SetGasLimit(limit uint64) {
	if w.tx.AuthInfo.Fee == nil {
		w.tx.AuthInfo.Fee = &Fee{}
	}

	w.tx.AuthInfo.Fee.GasLimit = limit

	// set authInfoBz to nil because the cached authInfoBz no longer matches tx.AuthInfo
	w.authInfoBz = nil
}

func (w *wrapper) SetFeeAmount(coins sdk.Coins) {
	if w.tx.AuthInfo.Fee == nil {
		w.tx.AuthInfo.Fee = &Fee{}
	}

	w.tx.AuthInfo.Fee.Amount = coins

	// set authInfoBz to nil because the cached authInfoBz no longer matches tx.AuthInfo
	w.authInfoBz = nil
}

func (w *wrapper) SetFeePayer(feePayer sdk.AccAddress) {
	if w.tx.AuthInfo.Fee == nil {
		w.tx.AuthInfo.Fee = &Fee{}
	}

	w.tx.AuthInfo.Fee.Payer = feePayer.String()

	// set authInfoBz to nil because the cached authInfoBz no longer matches tx.AuthInfo
	w.authInfoBz = nil
}

func (w *wrapper) SetFeeGranter(feeGranter sdk.AccAddress) {
	if w.tx.AuthInfo.Fee == nil {
		w.tx.AuthInfo.Fee = &Fee{}
	}

	w.tx.AuthInfo.Fee.Granter = feeGranter.String()

	// set authInfoBz to nil because the cached authInfoBz no longer matches tx.AuthInfo
	w.authInfoBz = nil
}

func (w *wrapper) SetSignatures(signatures ...signing.SignatureV2) error {
	n := len(signatures)
	signerInfos := make([]*SignerInfo, n)
	rawSigs := make([][]byte, n)

	for i, sig := range signatures {
		var modeInfo *ModeInfo
		modeInfo, rawSigs[i] = SignatureDataToModeInfoAndSig(sig.Data)
		any, err := PubKeyToAny(sig.PubKey)
		if err != nil {
			return err
		}
		signerInfos[i] = &SignerInfo{
			PublicKey: any,
			ModeInfo:  modeInfo,
			Sequence:  sig.Sequence,
		}
	}

	w.setSignerInfos(signerInfos)
	w.setSignatures(rawSigs)

	return nil
}

func (w *wrapper) setSignerInfos(infos []*SignerInfo) {
	w.tx.AuthInfo.SignerInfos = infos
	// set authInfoBz to nil because the cached authInfoBz no longer matches tx.AuthInfo
	w.authInfoBz = nil
}

func (w *wrapper) setSignatures(sigs [][]byte) {
	w.tx.Signatures = sigs
}

func (w *wrapper) GetTx() sdk.Tx {
	return w
}

// GetProtoTx returns the tx as a proto.Message.
func (w *wrapper) AsAny() *codectypes.Any {
	// We're sure here that w.tx is a proto.Message, so this will call
	// codectypes.NewAnyWithValue under the hood.
	return codectypes.UnsafePackAny(w.tx)
}

// WrapTx creates a TxBuilder wrapper around a tx.Tx proto message.
func WrapTx(protoTx *Tx) types.TxBuilder {
	return &wrapper{
		tx: protoTx,
	}
}

func (w *wrapper) GetExtensionOptions() []*codectypes.Any {
	return w.tx.Body.ExtensionOptions
}

func (w *wrapper) GetNonCriticalExtensionOptions() []*codectypes.Any {
	return w.tx.Body.NonCriticalExtensionOptions
}

func (w *wrapper) SetExtensionOptions(extOpts ...*codectypes.Any) {
	w.tx.Body.ExtensionOptions = extOpts
	w.bodyBz = nil
}

func (w *wrapper) SetNonCriticalExtensionOptions(extOpts ...*codectypes.Any) {
	w.tx.Body.NonCriticalExtensionOptions = extOpts
	w.bodyBz = nil
}

// PubKeyToAny converts a crypto.PubKey to a proto Any.
func PubKeyToAny(key crypto.PubKey) (*codectypes.Any, error) {
	protoMsg, ok := key.(proto.Message)
	if !ok {
		return nil, sdk.Wrap(fmt.Errorf("err invalid key, can't proto encode %T", protoMsg))
	}
	return codectypes.NewAnyWithValue(protoMsg)
}
