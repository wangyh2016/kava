package types

import "encoding/binary"

const (
	// ModuleName is the name of the module
	ModuleName = "liquidator"

	// StoreKey is the store key string for liquidator
	StoreKey = ModuleName

	// RouterKey is the message route for liquidator
	RouterKey = ModuleName

	// QuerierRoute is the querier route for liquidator
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName

	// DefaultCodespace default name for codespace
	DefaultCodespace = ModuleName
)

var sep = []byte(":")

// key value pairs for liquidator module
// 0x00:<collateralAuctionIdBytes> <DepositAmount>
// 0x01:<collateralBytes> <TotalCollateralAtAuction>
// KVStore key prefixes
var (
	CollateralAuctionKeyPrefix = []byte{0x00}
	TotalCollateralPrefix      = []byte{0x01}
)

// GetAuctionIDBytes returns the byte representation of the cdpID
func GetAuctionIDBytes(auctionID uint64) (auctionIDBz []byte) {
	auctionIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(auctionIDBz, auctionID)
	return
}

// GetAuctionIDFromBytes returns auctionID in uint64 format from a byte array
func GetAuctionIDFromBytes(bz []byte) (auctionID uint64) {
	return binary.BigEndian.Uint64(bz)
}

func createKey(bytes ...[]byte) (r []byte) {
	for _, b := range bytes {
		r = append(r, b...)
	}
	return
}
