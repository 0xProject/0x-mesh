package wrappers

// HACK(fabio): The `abi-gen` tool is unable to properly name the structs in generated from the
// contract ABI so we add type aliases here to correct the struct names.

// OrderWithoutExchangeAddress is a 0x order representation expected by the smart contracts.
type OrderWithoutExchangeAddress = Struct1

// TraderInfo contains all the order-relevant trader balances & allowances.
type TraderInfo = Struct2

// OrderInfo contains the status and filled amount of an order.
type OrderInfo = Struct3
