package meshdb

// var contractAddresses = ethereum.GanacheAddresses

// func TestOrderCRUDOperations(t *testing.T) {
// 	meshDB, err := New("/tmp/meshdb_testing/"+uuid.New().String(), contractAddresses)
// 	require.NoError(t, err)
// 	defer meshDB.Close()

// 	makerAddress := constants.GanacheAccount0
// 	salt := big.NewInt(1548619145450)
// 	o := &zeroex.Order{
// 		ChainID:               big.NewInt(constants.TestChainID),
// 		ExchangeAddress:       contractAddresses.Exchange,
// 		MakerAddress:          makerAddress,
// 		TakerAddress:          constants.NullAddress,
// 		SenderAddress:         constants.NullAddress,
// 		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 		TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 		TakerFeeAssetData:     constants.NullBytes,
// 		MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 		MakerFeeAssetData:     constants.NullBytes,
// 		Salt:                  salt,
// 		MakerFee:              big.NewInt(0),
// 		TakerFee:              big.NewInt(0),
// 		MakerAssetAmount:      big.NewInt(3551808554499581700),
// 		TakerAssetAmount:      big.NewInt(1),
// 		ExpirationTimeSeconds: big.NewInt(1548619325),
// 	}
// 	signedOrder, err := zeroex.SignTestOrder(o)
// 	require.NoError(t, err)

// 	orderHash, err := o.ComputeOrderHash()
// 	require.NoError(t, err)

// 	currentTime := time.Now().UTC()
// 	fiveMinutesFromNow := currentTime.Add(5 * time.Minute)

// 	// Insert
// 	order := &Order{
// 		Hash:                     orderHash,
// 		SignedOrder:              signedOrder,
// 		FillableTakerAssetAmount: big.NewInt(1),
// 		LastUpdated:              currentTime,
// 		IsRemoved:                false,
// 	}
// 	require.NoError(t, meshDB.Orders.Insert(order))
// 	// We need to call ResetHash so that unexported hash field is equal in later
// 	// assertions.
// 	signedOrder.ResetHash()

// 	// Find
// 	foundOrder := &Order{}
// 	require.NoError(t, meshDB.Orders.FindByID(order.ID(), foundOrder))
// 	assert.Equal(t, order, foundOrder)

// 	// Check Indexes
// 	orders, err := meshDB.FindOrdersByMakerAddressAndMaxSalt(makerAddress, salt)
// 	require.NoError(t, err)
// 	assert.Equal(t, []*Order{order}, orders)

// 	orders, err = meshDB.FindOrdersByMakerAddress(makerAddress)
// 	require.NoError(t, err)
// 	assert.Equal(t, []*Order{order}, orders)

// 	orders, err = meshDB.FindOrdersLastUpdatedBefore(fiveMinutesFromNow)
// 	require.NoError(t, err)
// 	assert.Equal(t, []*Order{order}, orders)

// 	// Update
// 	modifiedOrder := foundOrder
// 	modifiedOrder.FillableTakerAssetAmount = big.NewInt(0)
// 	require.NoError(t, meshDB.Orders.Update(modifiedOrder))
// 	foundModifiedOrder := &Order{}
// 	require.NoError(t, meshDB.Orders.FindByID(modifiedOrder.ID(), foundModifiedOrder))
// 	assert.Equal(t, modifiedOrder, foundModifiedOrder)

// 	// Delete
// 	require.NoError(t, meshDB.Orders.Delete(foundModifiedOrder.ID()))
// 	nonExistentOrder := &Order{}
// 	err = meshDB.Orders.FindByID(foundModifiedOrder.ID(), nonExistentOrder)
// 	assert.IsType(t, db.NotFoundError{}, err)
// }

// func TestParseContractAddressesAndTokenIdsFromAssetData(t *testing.T) {
// 	// ERC20 AssetData
// 	erc20AssetData := common.Hex2Bytes("f47261b000000000000000000000000038ae374ecf4db50b0ff37125b591a04997106a32")
// 	singleAssetDatas, err := parseContractAddressesAndTokenIdsFromAssetData(erc20AssetData, contractAddresses)
// 	require.NoError(t, err)
// 	assert.Len(t, singleAssetDatas, 1)
// 	expectedAddress := common.HexToAddress("0x38ae374ecf4db50b0ff37125b591a04997106a32")
// 	assert.Equal(t, expectedAddress, singleAssetDatas[0].Address)
// 	var expectedTokenID *big.Int
// 	assert.Equal(t, expectedTokenID, singleAssetDatas[0].TokenID)

// 	// ERC721 AssetData
// 	erc721AssetData := common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001")
// 	singleAssetDatas, err = parseContractAddressesAndTokenIdsFromAssetData(erc721AssetData, contractAddresses)
// 	require.NoError(t, err)
// 	assert.Equal(t, 1, len(singleAssetDatas))
// 	expectedAddress = common.HexToAddress("0x1dC4c1cEFEF38a777b15aA20260a54E584b16C48")
// 	assert.Equal(t, expectedAddress, singleAssetDatas[0].Address)
// 	expectedTokenID = big.NewInt(1)
// 	assert.Equal(t, expectedTokenID, singleAssetDatas[0].TokenID)

// 	// Multi AssetData
// 	multiAssetData := common.Hex2Bytes("94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004600000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000x94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004600000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000")
// 	singleAssetDatas, err = parseContractAddressesAndTokenIdsFromAssetData(multiAssetData, contractAddresses)
// 	require.NoError(t, err)
// 	assert.Equal(t, 2, len(singleAssetDatas))
// 	expectedSingleAssetDatas := []singleAssetData{
// 		singleAssetData{
// 			Address: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
// 		},
// 		singleAssetData{
// 			Address: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
// 			TokenID: big.NewInt(1),
// 		},
// 	}
// 	for i, singleAssetData := range singleAssetDatas {
// 		expectedSingleAssetData := expectedSingleAssetDatas[i]
// 		assert.Equal(t, expectedSingleAssetData.Address, singleAssetData.Address)
// 		assert.Equal(t, expectedSingleAssetData.TokenID, singleAssetData.TokenID)
// 	}
// }

// func TestTrimOrdersByExpirationTime(t *testing.T) {
// 	meshDB, err := New("/tmp/meshdb_testing/"+uuid.New().String(), contractAddresses)
// 	require.NoError(t, err)
// 	defer meshDB.Close()

// 	// TODO(albrow): Move these to top of file.
// 	makerAddress := constants.GanacheAccount0

// 	// Note: most of the fields in these orders are the same. For the purposes of
// 	// this test, the only thing that matters is the Salt and ExpirationTime.
// 	rawUnpinnedOrders := []*zeroex.Order{
// 		{
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerFeeAssetData:     constants.NullBytes,
// 			Salt:                  big.NewInt(0),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(100),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 		},
// 		{
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerFeeAssetData:     constants.NullBytes,
// 			Salt:                  big.NewInt(1),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(200),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 		},
// 		{
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerFeeAssetData:     constants.NullBytes,
// 			Salt:                  big.NewInt(2),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(200),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 		},
// 		{
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerFeeAssetData:     constants.NullBytes,
// 			Salt:                  big.NewInt(3),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(300),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 		},
// 	}
// 	rawPinnedOrders := []*zeroex.Order{
// 		{
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerFeeAssetData:     constants.NullBytes,
// 			Salt:                  big.NewInt(0),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(250),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 		},
// 		{
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerFeeAssetData:     constants.NullBytes,
// 			Salt:                  big.NewInt(1),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(350),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 		},
// 	}

// 	insertRawOrders(t, meshDB, rawUnpinnedOrders, false)
// 	pinnedOrders := insertRawOrders(t, meshDB, rawPinnedOrders, true)

// 	// Call CalculateNewMaxExpirationTimeAndTrimDatabase and check the results.
// 	targetMaxOrders := 4
// 	gotExpirationTime, gotRemovedOrders, err := meshDB.TrimOrdersByExpirationTime(targetMaxOrders)
// 	require.NoError(t, err)
// 	assert.Equal(t, "199", gotExpirationTime.String(), "newMaxExpirationTime")
// 	assert.Len(t, gotRemovedOrders, 2, "wrong number of orders removed")

// 	// Check that the expiration time of each removed order is >= the new max.
// 	for _, removedOrder := range gotRemovedOrders {
// 		expirationTimeOfRemovedOrder := removedOrder.SignedOrder.ExpirationTimeSeconds
// 		assert.True(t, expirationTimeOfRemovedOrder.Cmp(gotExpirationTime) != -1, "an order was removed with expiration time (%s) less than the new max (%s)", expirationTimeOfRemovedOrder, gotExpirationTime)
// 	}
// 	var remainingOrders []*Order
// 	require.NoError(t, meshDB.Orders.FindAll(&remainingOrders))
// 	assert.Len(t, remainingOrders, 4, "wrong number of orders remaining")

// 	// Check that the expiration time of each remaining order is <= the new max.
// 	for _, remainingOrder := range remainingOrders {
// 		if !remainingOrder.IsPinned {
// 			// Unpinned orders should not have an expiration time greater than the
// 			// new max.
// 			expirationTimeOfRemainingOrder := remainingOrder.SignedOrder.ExpirationTimeSeconds
// 			newMaxPlusOne := big.NewInt(0).Add(gotExpirationTime, big.NewInt(1))
// 			assert.True(t, expirationTimeOfRemainingOrder.Cmp(newMaxPlusOne) != 1, "a remaining order had an expiration time (%s) greater than the new max + 1 (%s)", expirationTimeOfRemainingOrder, newMaxPlusOne)
// 		}
// 	}

// 	// Check that the pinned orders are still in the database.
// 	for _, pinnedOrder := range pinnedOrders {
// 		require.NoError(t, meshDB.Orders.FindByID(pinnedOrder.Hash.Bytes(), &Order{}))
// 	}

// 	// Trying to trim orders when the database is full of pinned orders should
// 	// return an error.
// 	_, _, err = meshDB.TrimOrdersByExpirationTime(1)
// 	assert.EqualError(t, err, ErrDBFilledWithPinnedOrders.Error(), "expected ErrFilledWithPinnedOrders when targetMaxOrders is less than the number of pinned orders")
// }

// func TestFindOrdersByMakerAddressMakerFeeAssetAddressTokenID(t *testing.T) {
// 	meshDB, err := New("/tmp/meshdb_testing/"+uuid.New().String(), contractAddresses)
// 	require.NoError(t, err)
// 	defer meshDB.Close()

// 	makerAddress := constants.GanacheAccount0
// 	nextSalt := big.NewInt(1548619145450)

// 	zeroexOrders := []*zeroex.Order{
// 		// No Maker fee
// 		&zeroex.Order{
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			MakerFeeAssetData:     constants.NullBytes,
// 			Salt:                  nextSalt.Add(nextSalt, big.NewInt(1)),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(1548619325),
// 		},
// 		// ERC20 maker fee
// 		&zeroex.Order{
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			MakerFeeAssetData:     common.Hex2Bytes("f47261b000000000000000000000000038ae374ecf4db50b0ff37125b591a04997106a32"),
// 			Salt:                  nextSalt.Add(nextSalt, big.NewInt(1)),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(1548619325),
// 		},
// 		// ERC721 maker fee with token id = 1
// 		&zeroex.Order{
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			MakerFeeAssetData:     common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			Salt:                  nextSalt.Add(nextSalt, big.NewInt(1)),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(1548619325),
// 		},
// 		// ERC721 maker fee with token id = 2
// 		&zeroex.Order{
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			MakerFeeAssetData:     common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000002"),
// 			Salt:                  nextSalt.Add(nextSalt, big.NewInt(1)),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(1548619325),
// 		},
// 	}
// 	orders := make([]*Order, len(zeroexOrders))
// 	for i, o := range zeroexOrders {
// 		signedOrder, err := zeroex.SignTestOrder(o)
// 		require.NoError(t, err)
// 		orderHash, err := o.ComputeOrderHash()
// 		require.NoError(t, err)

// 		orders[i] = &Order{
// 			Hash:                     orderHash,
// 			SignedOrder:              signedOrder,
// 			FillableTakerAssetAmount: big.NewInt(1),
// 			LastUpdated:              time.Now().UTC(),
// 			IsRemoved:                false,
// 		}
// 		require.NoError(t, meshDB.Orders.Insert(orders[i]))
// 		// We need to call ResetHash so that unexported hash field is equal in later
// 		// assertions.
// 		signedOrder.ResetHash()
// 	}

// 	testCases := []struct {
// 		makerFeeAssetAddress common.Address
// 		makerFeeTokenID      *big.Int
// 		expectedOrders       []*Order
// 	}{
// 		{
// 			makerFeeAssetAddress: constants.NullAddress,
// 			makerFeeTokenID:      nil,
// 			expectedOrders:       orders[0:1],
// 		},
// 		{
// 			makerFeeAssetAddress: common.HexToAddress("0x38ae374ecf4db50b0ff37125b591a04997106a32"),
// 			makerFeeTokenID:      nil,
// 			expectedOrders:       orders[1:2],
// 		},
// 		{
// 			// Since no token id was specified, this query should match all token ids.
// 			makerFeeAssetAddress: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
// 			makerFeeTokenID:      nil,
// 			expectedOrders:       orders[2:4],
// 		},
// 		{
// 			makerFeeAssetAddress: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
// 			makerFeeTokenID:      big.NewInt(1),
// 			expectedOrders:       orders[2:3],
// 		},
// 		{
// 			makerFeeAssetAddress: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
// 			makerFeeTokenID:      big.NewInt(2),
// 			expectedOrders:       orders[3:4],
// 		},
// 	}
// 	for i, tc := range testCases {
// 		foundOrders, err := meshDB.FindOrdersByMakerAddressMakerFeeAssetAddressAndTokenID(makerAddress, tc.makerFeeAssetAddress, tc.makerFeeTokenID)
// 		require.NoError(t, err)
// 		assert.Equal(t, tc.expectedOrders, foundOrders, "test case %d", i)
// 	}
// }

// func insertRawOrders(t *testing.T, meshDB *MeshDB, rawOrders []*zeroex.Order, isPinned bool) []*Order {
// 	results := make([]*Order, len(rawOrders))
// 	for i, order := range rawOrders {
// 		// Sign, compute order hash, and insert.
// 		signedOrder, err := zeroex.SignTestOrder(order)
// 		require.NoError(t, err)
// 		orderHash, err := order.ComputeOrderHash()
// 		require.NoError(t, err)

// 		order := &Order{
// 			Hash:                     orderHash,
// 			SignedOrder:              signedOrder,
// 			FillableTakerAssetAmount: big.NewInt(1),
// 			LastUpdated:              time.Now(),
// 			IsRemoved:                false,
// 			IsPinned:                 isPinned,
// 		}
// 		results[i] = order
// 		require.NoError(t, meshDB.Orders.Insert(order))
// 	}
// 	return results
// }

// func TestPruneMiniHeadersAboveRetentionLimit(t *testing.T) {
// 	t.Parallel()

// 	meshDB, err := New("/tmp/meshdb_testing/"+uuid.New().String(), contractAddresses)
// 	require.NoError(t, err)
// 	defer meshDB.Close()

// 	txn := meshDB.MiniHeaders.OpenTransaction()
// 	defer func() {
// 		_ = txn.Discard()
// 	}()

// 	miniHeadersToAdd := miniHeadersMaxPerPage*2 + defaultMiniHeaderRetentionLimit + 1
// 	for i := 0; i < miniHeadersToAdd; i++ {
// 		miniHeader := &miniheader.MiniHeader{
// 			Hash:      common.BigToHash(big.NewInt(int64(i))),
// 			Number:    big.NewInt(int64(i)),
// 			Timestamp: time.Now().Add(time.Duration(i)*time.Second - 5*time.Hour),
// 		}
// 		require.NoError(t, txn.Insert(miniHeader))
// 	}
// 	require.NoError(t, txn.Commit())

// 	require.NoError(t, meshDB.PruneMiniHeadersAboveRetentionLimit())
// 	remainingMiniHeaders, err := meshDB.MiniHeaders.Count()
// 	assert.Equal(t, defaultMiniHeaderRetentionLimit, remainingMiniHeaders, "wrong number of MiniHeaders remaining")
// }
