import { gql } from '@apollo/client/core';

export const statsQuery = gql`
    query Stats {
        stats {
            version
            pubSubTopic
            rendezvous
            peerID
            ethereumChainID
            latestBlock {
                number
                hash
            }
            numPeers
            numOrders
            numOrdersV4
            numOrdersIncludingRemoved
            numOrdersIncludingRemovedV4
            numPinnedOrders
            numPinnedOrdersV4
            startOfCurrentUTCDay
            ethRPCRequestsSentInCurrentUTCDay
            ethRPCRateLimitExpiredRequests
            maxExpirationTime
        }
    }
`;

export const addOrdersMutation = gql`
    mutation AddOrders(
        $orders: [NewOrder!]!
        $pinned: Boolean = true
        $opts: AddOrdersOpts = { keepCancelled: false, keepExpired: false, keepFullyFilled: false, keepUnfunded: false }
    ) {
        addOrders(orders: $orders, pinned: $pinned, opts: $opts) {
            accepted {
                order {
                    hash
                    chainId
                    exchangeAddress
                    makerAddress
                    makerAssetData
                    makerAssetAmount
                    makerFeeAssetData
                    makerFee
                    takerAddress
                    takerAssetData
                    takerAssetAmount
                    takerFeeAssetData
                    takerFee
                    senderAddress
                    feeRecipientAddress
                    expirationTimeSeconds
                    salt
                    signature
                    fillableTakerAssetAmount
                }
                isNew
            }
            rejected {
                hash
                code
                message
                order {
                    chainId
                    exchangeAddress
                    makerAddress
                    makerAssetData
                    makerAssetAmount
                    makerFeeAssetData
                    makerFee
                    takerAddress
                    takerAssetData
                    takerAssetAmount
                    takerFeeAssetData
                    takerFee
                    senderAddress
                    feeRecipientAddress
                    expirationTimeSeconds
                    salt
                    signature
                }
            }
        }
    }
`;

export const addOrdersMutationV4 = gql`
    mutation AddOrdersV4(
        $orders: [NewOrderV4!]!
        $pinned: Boolean = true
        $opts: AddOrdersOpts = { keepCancelled: false, keepExpired: false, keepFullyFilled: false, keepUnfunded: false }
    ) {
        addOrdersV4(orders: $orders, pinned: $pinned, opts: $opts) {
            accepted {
                order {
                    chainId
                    exchangeAddress
                    makerToken
                    takerToken
                    makerAmount
                    takerAmount
                    takerTokenFeeAmount
                    maker
                    taker
                    sender
                    feeRecipient
                    pool
                    expiry
                    salt
                    signatureType
                    signatureV
                    signatureR
                    signatureS
                }
                isNew
            }
            rejected {
                code
                message
                hash
                order {
                    chainId
                    exchangeAddress
                    makerToken
                    takerToken
                    makerAmount
                    takerAmount
                    takerTokenFeeAmount
                    maker
                    taker
                    sender
                    feeRecipient
                    pool
                    expiry
                    salt
                    signatureType
                    signatureV
                    signatureR
                    signatureS
                }
            }
        }
    }
`;

export const orderQuery = gql`
    query Order($hash: String!) {
        order(hash: $hash) {
            hash
            chainId
            exchangeAddress
            makerAddress
            makerAssetData
            makerAssetAmount
            makerFeeAssetData
            makerFee
            takerAddress
            takerAssetData
            takerAssetAmount
            takerFeeAssetData
            takerFee
            senderAddress
            feeRecipientAddress
            expirationTimeSeconds
            salt
            signature
            fillableTakerAssetAmount
        }
    }
`;

export const orderQueryV4 = gql`
    query OrderV4($hash: String!) {
        orderv4(hash: $hash) {
            hash
            chainId
            exchangeAddress
            makerToken
            takerToken
            makerAmount
            takerAmount
            takerTokenFeeAmount
            maker
            taker
            sender
            feeRecipient
            pool
            expiry
            salt
            signatureType
            signatureV
            signatureR
            signatureS
            fillableTakerAssetAmount
        }
    }
`;

export const ordersQuery = gql`
    query Orders(
        $filters: [OrderFilter!] = []
        $sort: [OrderSort!] = [{ field: hash, direction: ASC }]
        $limit: Int = 100
    ) {
        orders(filters: $filters, sort: $sort, limit: $limit) {
            hash
            chainId
            exchangeAddress
            makerAddress
            makerAssetData
            makerAssetAmount
            makerFeeAssetData
            makerFee
            takerAddress
            takerAssetData
            takerAssetAmount
            takerFeeAssetData
            takerFee
            senderAddress
            feeRecipientAddress
            expirationTimeSeconds
            salt
            signature
            fillableTakerAssetAmount
        }
    }
`;

export const ordersQueryV4 = gql`
    query Orders(
        $filters: [OrderFilterV4!] = []
        $sort: [OrderSortV4!] = [{ field: hash, direction: ASC }]
        $limit: Int = 100
    ) {
        ordersv4(filters: $filters, sort: $sort, limit: $limit) {
            hash
            chainId
            exchangeAddress
            makerToken
            takerToken
            makerAmount
            takerAmount
            takerTokenFeeAmount
            maker
            taker
            sender
            feeRecipient
            pool
            expiry
            salt
            signatureType
            signatureV
            signatureR
            signatureS
        }
    }
`;

export const orderEventsSubscription = gql`
    subscription {
        orderEvents {
            timestamp
            endState
            order {
                hash
                chainId
                exchangeAddress
                makerAddress
                makerAssetData
                makerAssetAmount
                makerFeeAssetData
                makerFee
                takerAddress
                takerAssetData
                takerAssetAmount
                takerFeeAssetData
                takerFee
                senderAddress
                feeRecipientAddress
                expirationTimeSeconds
                salt
                signature
                fillableTakerAssetAmount
            }
            orderv4 {
                chainId
                exchangeAddress
                makerToken
                takerToken
                makerAmount
                takerAmount
                takerTokenFeeAmount
                maker
                taker
                sender
                feeRecipient
                pool
                expiry
                salt
                signatureType
                signatureV
                signatureR
                signatureS
                hash
                fillableTakerAssetAmount
            }
            contractEvents {
                blockHash
                txHash
                txIndex
                logIndex
                isRemoved
                address
                kind
                parameters
            }
        }
    }
`;
