package indexers

import (
	"errors"
	"github.com/wificoin-project/wfcd/blockchain"
	"github.com/wificoin-project/wfcd/chaincfg/chainhash"
	"github.com/wificoin-project/wfcd/database"
	"github.com/wificoin-project/wfcd/txscript"
	"github.com/wificoin-project/wfcutil"
)

const spentIndexName = "spent index"

var (
	spentIndexKey  = []byte("spentbytxidx")
	spentKeySize   = 32 + 4
	spentValueSize = 32 + 4 + 4 + 8 + 1 + 20
)

// The serialized key format is:
//
//   Field           Type              Size
//   txhash          chainhash.Hash    32 bytes
//   output index    uint32            4 bytes
//   -----
//   Total: 36 bytes
//
// The serialized value format is:
//
//   Field           Type              Size
//   txhash          chainhash.Hash    32 bytes
//   input index     uint32            4 bytes
//   blockheight     uint32            4 bytes
//   wfc amount      int64             8 bytes
//   addr type       uint8             1 byte
//   addr hash       hash160           20 bytes
//   -----
//   Total: 69 bytes
// -----------------------------------------------------------------------------

type SpentIndexValue []byte

func (v SpentIndexValue) TxHash() string {
	var hash chainhash.Hash
	copy(hash[:], v[:32])

	return hash.String()
}

func (v SpentIndexValue) Index() uint32 {
	return byteOrder.Uint32(v[32:36])
}

func (v SpentIndexValue) Height() uint32 {
	return byteOrder.Uint32(v[36:40])
}

func (v SpentIndexValue) Value() int64 {
	return int64(byteOrder.Uint64(v[40:48]))
}

func (v SpentIndexValue) Address() []byte {
	return v[48:]
}

func serializeSpentIndexKey(txid *chainhash.Hash, outputIdx uint32) []byte {
	// Serialize the entry.
	serialized := make([]byte, spentKeySize)

	copy(serialized, txid[:])
	byteOrder.PutUint32(serialized[32:], outputIdx)
	return serialized
}

func serializeSpentIndexvalue(txhash *chainhash.Hash, inputIdx int, height int32, value int64, addr [addrKeySize]byte) []byte {
	serialized := make([]byte, spentValueSize)

	copy(serialized, txhash[:])
	byteOrder.PutUint32(serialized[32:], uint32(inputIdx))
	byteOrder.PutUint32(serialized[36:], uint32(height))
	byteOrder.PutUint64(serialized[40:], uint64(value))
	copy(serialized[48:], addr[:])

	return serialized
}

func dbPutSpentIndexEntry(dbTx database.Tx, key, value []byte) error {
	meta := dbTx.Metadata()
	spentIndex := meta.Bucket(spentIndexKey)
	return spentIndex.Put(key[:], value[:])
}
func dbRemoveSpentIndexEntry(dbTx database.Tx, key []byte) error {

	meta := dbTx.Metadata()
	spentIndex := meta.Bucket(timestampIndexKey)
	return spentIndex.Delete(key)
}

type SpentIndex struct {
	db database.DB
}

var _ Indexer = (*SpentIndex)(nil)

func (idx *SpentIndex) Init() error {
	return nil
}

func (idx *SpentIndex) Key() []byte {
	return spentIndexKey
}

func (idx *SpentIndex) Name() string {
	return spentIndexName
}

func (idx *SpentIndex) Create(dbTx database.Tx) error {
	meta := dbTx.Metadata()
	_, err := meta.CreateBucket(spentIndexKey)

	return err
}

func (idx *SpentIndex) ConnectBlock(dbTx database.Tx, block *wfcutil.Block, stxos []blockchain.SpentTxOut) error {
	blockHeight := block.Height()

	stxoIndex := 0
	for txIdx, tx := range block.Transactions() {
		if txIdx == 0 {
			// is a coinbase
			continue
		}

		txhash := tx.Hash()
		for i, vin := range tx.MsgTx().TxIn {
			var addr [addrKeySize]byte
			pkScript := stxos[stxoIndex].PkScript

			class := txscript.GetScriptClass(pkScript)
			if class == txscript.ScriptHashTy {
				addr[0] = addrKeyTypeScriptHash
				copy(addr[1:], pkScript[3:23])
			} else if class == txscript.PubKeyHashTy {
				addr[0] = addrKeyTypePubKeyHash
				copy(addr[1:], pkScript[2:22])
			} else {
				// unsupported address types.
				addr[0] = 10
			}

			CSpentIndexKey := serializeSpentIndexKey(&vin.PreviousOutPoint.Hash, vin.PreviousOutPoint.Index)
			CSPentIndexValue := serializeSpentIndexvalue(txhash, i, blockHeight, stxos[stxoIndex].Amount, addr)

			if err := dbPutSpentIndexEntry(dbTx, CSpentIndexKey, CSPentIndexValue); err != nil {
				return err
			}

			stxoIndex++
		}
	}

	return nil
}

func (idx *SpentIndex) DisconnectBlock(dbTx database.Tx, block *wfcutil.Block, stox []blockchain.SpentTxOut) error {
	for txIdx, tx := range block.Transactions() {
		if txIdx == 0 {
			continue
		}

		for _, vin := range tx.MsgTx().TxIn {
			CSpentIndexKey := serializeSpentIndexKey(&vin.PreviousOutPoint.Hash, vin.PreviousOutPoint.Index)
			if err := dbRemoveSpentIndexEntry(dbTx, CSpentIndexKey); err != nil {
				return err
			}
		}
	}

	return nil
}

func (idx *SpentIndex) Get(txHash *chainhash.Hash, index uint32) (SpentIndexValue, error) {
	key := serializeSpentIndexKey(txHash, index)
	value := make([]byte, spentValueSize)

	err := idx.db.View(func(dbTx database.Tx) error {
		bucket := dbTx.Metadata().Bucket(spentIndexKey)
		data := bucket.Get(key)
		if len(data) != spentValueSize {
			return errors.New("Unable to get spent info")
		}

		copy(value[:], data)
		return nil
	})

	return SpentIndexValue(value), err
}

func NewSpentIndex(db database.DB) *SpentIndex {
	return &SpentIndex{
		db: db,
	}
}

func DropSpentIndex(db database.DB, interrupt <-chan struct{}) error {
	return dropIndex(db, timestampIndexKey, timestampIndexName, interrupt)
}
