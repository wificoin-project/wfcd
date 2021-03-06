package indexers

import (
	"github.com/wificoin-project/wfcd/blockchain"
	"github.com/wificoin-project/wfcd/chaincfg/chainhash"
	"github.com/wificoin-project/wfcd/database"
	"github.com/wificoin-project/wfcutil"
)

const timestampIndexName = "timestamp index"

var (
	timestampIndexKey = []byte("hashbytimestampidx")
)

// The serialized key format is:
//
//   Field           Type      Size
//   block time      uint64    8 bytes
//   -----
//   Total: 8 bytes
//
// The serialized value format is:
//
//   Field           Type      Size
//   block hash      uint32    32 bytes
//   -----
//   Total: 32 bytes
// -----------------------------------------------------------------------------

func dbRemoveTimestampIndexEntry(dbTx database.Tx, blockTime int64) error {
	var serializedTimestamp [8]byte
	byteOrder.PutUint64(serializedTimestamp[:], uint64(blockTime))

	meta := dbTx.Metadata()
	timeIndex := meta.Bucket(timestampIndexKey)
	return timeIndex.Delete(serializedTimestamp[:])
}

func dbPutTimestampIndexEntry(dbTx database.Tx, hash *chainhash.Hash, blockTime int64) error {
	var serializedTimestamp [8]byte
	byteOrder.PutUint64(serializedTimestamp[:], uint64(blockTime))

	meta := dbTx.Metadata()
	timeIndex := meta.Bucket(timestampIndexKey)
	return timeIndex.Put(serializedTimestamp[:], hash[:])
}

func dbReadTimestampIndex(dbTx database.Tx, high, low uint64, hashes *[]string) error {
	timIndex := dbTx.Metadata().Bucket(timestampIndexKey)
	return timIndex.ForEach(func(k, v []byte) error {
		timestamp := byteOrder.Uint64(k)
		if timestamp >= low && timestamp < high {
			var hash chainhash.Hash
			copy(hash[:], v)
			*hashes = append(*hashes, hash.String())
		}
		return nil
	})
}

type TimestampIndex struct {
	db database.DB
}

var _ Indexer = (*TimestampIndex)(nil)

func (idx *TimestampIndex) Init() error {
	return nil
}

func (idx *TimestampIndex) Key() []byte {
	return timestampIndexKey
}

func (idx *TimestampIndex) Name() string {
	return timestampIndexName
}

func (idx *TimestampIndex) Create(dbTx database.Tx) error {
	meta := dbTx.Metadata()
	_, err := meta.CreateBucket(timestampIndexKey)
	return err
}

func (idx *TimestampIndex) ConnectBlock(dbTx database.Tx, block *wfcutil.Block,
	stxos []blockchain.SpentTxOut) error {
	newTimestamp := block.MsgBlock().Header.Timestamp.Unix()

	err := dbPutTimestampIndexEntry(dbTx, block.Hash(), newTimestamp)
	if err != nil {
		return err
	}

	return nil
}
func (idx *TimestampIndex) DisconnectBlock(dbTx database.Tx, block *wfcutil.Block,
	stxos []blockchain.SpentTxOut) error {
	return dbRemoveTimestampIndexEntry(dbTx, block.MsgBlock().Header.Timestamp.Unix())
}

func (idx *TimestampIndex) Get(high, low uint64) ([]string, error) {
	var hashes = []string{}
	var err error

	if high > low {
		err = idx.db.View(func(dbTx database.Tx) error {
			return dbReadTimestampIndex(dbTx, high, low, &hashes)
		})
	}

	return hashes, err
}

func NewTimeStampIndex(db database.DB) *TimestampIndex {
	return &TimestampIndex{db: db}
}

func DropTimeStampIndex(db database.DB, interrupt <-chan struct{}) error {
	return dropIndex(db, timestampIndexKey, timestampIndexName, interrupt)
}
