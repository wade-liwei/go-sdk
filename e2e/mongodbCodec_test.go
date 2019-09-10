package e2e

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/binance-chain/go-sdk/common/types"
	txMsg "github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	"github.com/tendermint/tendermint/libs/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	//"github.com/tendermint/tendermint/libs/common"
)

func TestHexRange(t *testing.T) {

	d, err := hex.DecodeString("000102030405060708090A0B0C0D0E0F")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", d)
}

func TestOtherFileFunc(t *testing.T) {
	TestTxInfoSearch(t)
}

func TestTxInfoData(t *testing.T) {
	txInfos, err := TxInfoSearch()
	if err != nil {
		t.Error(err)
	}
	for k, v := range txInfos {
		fmt.Printf("k: %v   txInfo: %v \n", k, v)

		for _, msgV := range v.Tx.GetMsgs() {
			//switch t := t.(type) {
			switch msgVV := msgV.(type) {
			case txMsg.CreateOrderMsg:
				fmt.Printf("sender byte: %v   string: %v \n", msgVV.Sender.Bytes(), msgVV.Sender.String())
			}
		}
	}
}

func TestTxInfoInterface(t *testing.T) {
	//TestTxInfoData(t)
	txInfos, err := TxInfoSearch()
	if err != nil {
		t.Error(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Error(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Error(err)
	}

	//codec -------------------------------------------------------------------------------------------------------
	rb := bson.NewRegistryBuilder()
	//byteArrCod := &tx.InfoTXInterface{}
	tendermintHexCod := &common.HexBytes{}
	accAddressCod := &types.AccAddress{}
	votOptionCod := new(txMsg.VoteOption)
	proposalKindCod := new(txMsg.ProposalKind)
	//rb.RegisterCodec(reflect.TypeOf([]byte("")), byteArrCod)
	rb.RegisterCodec(reflect.TypeOf(common.HexBytes(nil)), tendermintHexCod)
	rb.RegisterCodec(reflect.TypeOf(types.AccAddress(nil)), accAddressCod)
	rb.RegisterCodec(reflect.TypeOf(txMsg.VoteOption('0')), votOptionCod)
	rb.RegisterCodec(reflect.TypeOf(txMsg.ProposalKind('0')), proposalKindCod)

	collection := client.Database("TxMongodbCodec").Collection("TxMongodbCodec", options.Collection().SetRegistry(rb.Build()))
	//codec -------------------------------------------------------------------------------------------------------=
	//collection := client.Database("TxMongodbCodec").Collection("TxMongodbCodec")
	for k, v := range txInfos {
		tmp := tx.InfoTXInterface{}
		tmp.Hash = v.Hash
		tmp.Height = v.Height
		tmp.Result = v.Result
		fmt.Printf("tmp.Hash:  %v   len: %v  type:  %T \n", tmp.Hash.Bytes(), len(tmp.Hash.Bytes()), tmp.Hash)

		msgs := make([]interface{}, 0, len(v.Tx.GetMsgs()))
		for _, msgV := range v.Tx.GetMsgs() {
			msgs = append(msgs, msgV)
		}
		tmp.Tx = msgs

		res, err := collection.InsertOne(ctx, tmp)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("k: %v   txInfo: %v  \n      mongodbId: %v \n", k, tmp.Hash, res.InsertedID)
	}
}

type GetInterface interface {
	Get() string
	EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error
	DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error
}

type Person struct {
	PersonId int64 `bson:"personId123"`
	Name     string
	TestByte []byte
}

func (p *Person) Get() string {
	return p.Name
}

func (e *Person) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	fmt.Printf("---------------------------byte:  %v \n", string(val.Bytes()))
	return vw.WriteString(string(val.Bytes()))
}

// DecodeValue negates the value of ID when reading
func (e *Person) DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	i, err := vr.ReadInt64()
	if err != nil {
		return err
	}

	val.SetInt(i * -1)
	return nil
}

type NewCodec struct {
	ID123          int64 `bson:"_id123"`
	Data           int64
	EmptyInterface interface{}
	People         interface{}
}

func (e *NewCodec) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	return vw.WriteInt64(val.Int())
}

// DecodeValue negates the value of ID when reading
func (e *NewCodec) DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	i, err := vr.ReadInt64()
	if err != nil {
		return err
	}

	val.SetInt(i * -1)
	return nil
}

func TestClientRegistryPassedToCursors(t *testing.T) {
	// register a new codec for the int64 type that does the default encoding for an int64 and negates the value when
	// decoding

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Error(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Error(err)
	}

	rb := bson.NewRegistryBuilder()
	cod := &NewCodec{}
	rb.RegisterCodec(reflect.TypeOf(int64(0)), cod)
	p := &Person{}
	rb.RegisterCodec(reflect.TypeOf([]byte("")), p)

	collection := client.Database("TxMongodbCodec").Collection("TxMongodbCodec", options.Collection().SetRegistry(rb.Build()))
	//collection := client.Database("TxMongodbCodec").Collection("TxMongodbCodec")

	res, err := collection.InsertOne(ctx, NewCodec{ID123: 15, Data: 15, EmptyInterface: 13, People: &Person{Name: "wade", TestByte: []byte("000000000000000000000")}})

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Insert Id: %v \n", res.InsertedID)

}
