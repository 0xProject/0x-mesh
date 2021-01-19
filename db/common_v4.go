package db

type OrderFieldV4 string

const (
	OV4FHash                     OrderFieldV4 = "hash"
	OV4FChainID                  OrderFieldV4 = "chainID"
	OV4FExchangeAddress          OrderFieldV4 = "exchangeAddress"
	OV4FMakerToken               OrderFieldV4 = "makerToken"
	OV4FTakerToken               OrderFieldV4 = "takerToken"
	OV4FMakerAmount              OrderFieldV4 = "makerAmount"
	OV4FTakerAmount              OrderFieldV4 = "takerAmount"
	OV4FTakerTokenFeeAmount      OrderFieldV4 = "takerTokenFeeAmount"
	OV4FMaker                    OrderFieldV4 = "maker"
	OV4FTaker                    OrderFieldV4 = "taker"
	OV4FSender                   OrderFieldV4 = "sender"
	OV4FFeeRecipient             OrderFieldV4 = "feeRecipient"
	OV4FPool                     OrderFieldV4 = "pool"
	OV4FExpiry                   OrderFieldV4 = "expiry"
	OV4FSalt                     OrderFieldV4 = "salt"
	OV4FSignature                OrderFieldV4 = "signature"
	OV4FLastUpdated              OrderFieldV4 = "lastUpdated"
	OV4FFillableTakerAssetAmount OrderFieldV4 = "fillableTakerAssetAmount"
	OV4FIsRemoved                OrderFieldV4 = "isRemoved"
	OV4FIsPinned                 OrderFieldV4 = "isPinned"
	OV4FIsUnfillable             OrderFieldV4 = "isUnfillable"
	OV4FIsExpired                OrderFieldV4 = "isExpired"
	OV4FParsedMakerAssetData     OrderFieldV4 = "parsedMakerAssetData"
	OV4FParsedMakerFeeAssetData  OrderFieldV4 = "parsedMakerFeeAssetData"
	OV4FLastValidatedBlockNumber OrderFieldV4 = "lastValidatedBlockNumber"
	OV4FKeepCancelled            OrderFieldV4 = "keepCancelled"
	OV4FKeepExpired              OrderFieldV4 = "keepExpired"
	OV4FKeepFullyFilled          OrderFieldV4 = "keepFullyFilled"
	OV4FKeepUnfunded             OrderFieldV4 = "keepUnfunded"
)

type OrderQueryV4 struct {
	Filters []OrderFilterV4 `json:"filters"`
	Sort    []OrderSortV4   `json:"sort"`
	Limit   uint            `json:"limit"`
	Offset  uint            `json:"offset"`
}

type OrderSortV4 struct {
	Field     OrderFieldV4  `json:"field"`
	Direction SortDirection `json:"direction"`
}

type OrderFilterV4 struct {
	Field OrderFieldV4 `json:"field"`
	Kind  FilterKind   `json:"kind"`
	Value interface{}  `json:"value"`
}
