package tx

import (
	"fmt"

	"plugchain-sdk-go/codec"
	sdk "plugchain-sdk-go/types"
	signingtypes "plugchain-sdk-go/types/tx/signing"
)

type config struct {
	handler     sdk.SignModeHandler
	decoder     sdk.TxDecoder
	encoder     sdk.TxEncoder
	jsonDecoder sdk.TxDecoder
	jsonEncoder sdk.TxEncoder
	protoCodec  *codec.ProtoCodec
}

// NewTxConfig returns a new protobuf TxConfig using the provided ProtoCodec and sign modes. The
// first enabled sign mode will become the default sign mode.
func NewTxConfig(protoCodec *codec.ProtoCodec, enabledSignModes []signingtypes.SignMode) sdk.TxConfig {
	return &config{
		handler:     MakeSignModeHandler(enabledSignModes),
		decoder:     DefaultTxDecoder(protoCodec),
		encoder:     DefaultTxEncoder(),
		jsonDecoder: DefaultJSONTxDecoder(protoCodec),
		jsonEncoder: DefaultJSONTxEncoder(),
		protoCodec:  protoCodec,
	}
}

func (g config) NewTxBuilder() sdk.TxBuilder {
	return newBuilder()
}

// WrapTxBuilder returns a builder from provided transaction
func (g config) WrapTxBuilder(newTx sdk.Tx) (sdk.TxBuilder, error) {
	newBuilder, ok := newTx.(*wrapper)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, newTx)
	}

	return newBuilder, nil
}

func (g config) SignModeHandler() sdk.SignModeHandler {
	return g.handler
}

func (g config) TxEncoder() sdk.TxEncoder {
	return g.encoder
}

func (g config) TxDecoder() sdk.TxDecoder {
	return g.decoder
}

func (g config) TxJSONEncoder() sdk.TxEncoder {
	return g.jsonEncoder
}

func (g config) TxJSONDecoder() sdk.TxDecoder {
	return g.jsonDecoder
}
