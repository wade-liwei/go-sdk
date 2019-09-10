package tx

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/binance-chain/go-sdk/types/msg"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"

	"github.com/tendermint/tendermint/types"
)

const Source int64 = 0

type Tx interface {

	// Gets the Msg.
	GetMsgs() []msg.Msg
	GetSource() int64
	GetData() []byte
	GetMemo() string
	GetSignatures() []StdSignature
	//EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error
	//DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error
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
func (tx StdTx) GetSource() int64              { return tx.Source }
func (tx StdTx) GetData() []byte               { return tx.Data }
func (tx StdTx) GetMemo() string               { return tx.Memo }
func (tx StdTx) GetSignatures() []StdSignature { return tx.Signatures }



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
	Hash     common.HexBytes        `json:"hash"`
	Time     time.Time              `json:"time"`
	Height   int64                  `json:"height"`
	Tx       interface{}            `json:"tx"`
	Result   abci.ResponseDeliverTx `json:"result"`
	Index    uint32                 `json:"index"`
	Proof    types.TxProof          `json:"proof,omitempty"`
	MsgKey   string                 `json:"msg_key"`
	MsgValue string                 `json:"msg_value"`
}

func (e *InfoTXInterface) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {

	if val.IsValid() {
		return vw.WriteString(string(val.Bytes()))
	}
	return errors.New("InfoTXInterface encoder value is invalid.")
}

// DecodeValue negates the value of ID when reading
func (e *InfoTXInterface) DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	txInfoBytesType, b, err := vr.ReadBinary()
	_ = b
	if err != nil {
		return err
	}

	val.SetBytes(txInfoBytesType)

	return nil
}

type Info struct {
	Hash   common.HexBytes        `json:"hash"`
	Height int64                  `json:"height"`
	Tx     Tx                     `json:"tx"`
	Result abci.ResponseDeliverTx `json:"result"`
	Index  uint32                 `json:"index"`
	Proof  types.TxProof          `json:"proof,omitempty"`
}

func (tx Info) String() string {
	return fmt.Sprintf(`
		Hash   : %v
		Height : %v
		Tx     : %v
		Result : %v
		`, tx.Hash, tx.Height, tx.Tx, tx.Result)
}
