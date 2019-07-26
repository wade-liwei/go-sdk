package tx

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	// "github.com/binance-chain/go-sdk/common/bech32"
	//
	// binanceTypes "github.com/binance-chain/go-sdk/common/types"
)

const Source int64 = 2

type Tx interface {

	// Gets the Msg.
	GetMsgs() []msg.Msg
	//EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error
	//DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error
}

type MongodbTx interface {
	Tx
}

// StdTx def
type StdTx struct {
	Msgs       []msg.Msg      `json:"msg"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
	Source     int64          `json:"source"`
	Data       []byte         `json:"data"`
}

func (stdTx StdTx) String() string {
	return fmt.Sprintf(`
				Msgs       : %v
				Signatures : %v
				Memo       : %v
				Source     : %v
				Data       : %v
		`, stdTx.Msgs, stdTx.Signatures, stdTx.Memo, stdTx.Source, stdTx.Data)

}

// NewStdTx to instantiate an instance
func NewStdTx(msgs []msg.Msg, sigs []StdSignature, memo string, source int64, data []byte) StdTx {
	return StdTx{
		Msgs:       msgs,
		Signatures: sigs,
		Memo:       memo,
		Source:     source,
		Data:       data,
	}
}

// GetMsgs def
func (tx StdTx) GetMsgs() []msg.Msg { return tx.Msgs }

func (tx StdTx) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	fmt.Printf("codec  encode:   -------------------------------- %v \n", val.Int())
	return vw.WriteInt64(val.Int())
}

//DecodeValue negates the value of ID when reading
func (tx StdTx) DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	fmt.Println("---------------------------------------------------------------------------------------")
	fmt.Printf("val:  %v ----------------------------------------------------\n", val)
	i, err := vr.ReadInt64()
	if err != nil {
		return err
	}
	val.SetInt(i * -1)
	return nil
}

type InfoTXInterface struct {
	Hash   common.HexBytes         `json:"hash"`
	Height int64                   `json:"height"`
	Tx     interface{}             `json:"tx"`
	Result types.ResponseDeliverTx `json:"result"`
}

func (e *InfoTXInterface) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {

	if val.IsValid() {
		return vw.WriteString(string(val.Bytes()))
	}
	return errors.New("InfoTXInterface encoder value is invalid.")
}

// DecodeValue negates the value of ID when reading
func (e *InfoTXInterface) DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	i, err := vr.ReadInt64()
	if err != nil {
		return err
	}
	val.SetInt(i * -1)
	return nil
}

type Info struct {
	Hash   common.HexBytes         `json:"hash"`
	Height int64                   `json:"height"`
	Tx     Tx                      `json:"tx"`
	Result types.ResponseDeliverTx `json:"result"`
}

func (tx Info) String() string {
	return fmt.Sprintf(`
		Hash   : %v
		Height : %v
		Tx     : %v
		Result : %v
		`, tx.Hash, tx.Height, tx.Tx, tx.Result)
}
