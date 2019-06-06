package db

// type Transaction struct {
// 	*writeableCollection
// 	transaction *leveldb.Transaction
// }

// func (c *Collection) OpenTransaction() (*Transaction, error) {
// 	snapshot, err := c.ldb.OpenTransaction()
// 	if err != nil {
// 		return nil, err
// 	}
// 	c.indexMut.RLock()
// 	indexes := make([]*Index, len(c.indexes))
// 	copy(indexes, c.indexes)
// 	c.indexMut.RUnlock()
// 	return &Transaction{
// 		writeableCollection: &writeableCollection{
// 			readOnlyCollection: &readOnlyCollection{
// 				reader:    snapshot,
// 				name:      c.name,
// 				modelType: c.modelType,
// 				indexes:   indexes,
// 			},
// 		},
// 	}, nil
// }

// // Release releases the snapshot. This will not release any ongoing queries,
// // which will still finish unless the database is closed. Other methods should
// // not be called after the snapshot has been released.
// func (s *Snapshot) Release() {
// 	s.snapshot.Release()
// }
