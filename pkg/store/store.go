package store

import (
	"fmt"
	"github.com/linxGnu/grocksdb"
)

// Store represents the RocksDB store
type Store struct {
	db *grocksdb.DB
}



// NewStore initializes and returns a new Store
func NewStore(dbPath string) (*Store, error) {
	opts := grocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	db, err := grocksdb.OpenDb(opts, dbPath)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Save
func (s *Store) SaveItem(prefix string, id string, data []byte) (string, error) {
    wo := grocksdb.NewDefaultWriteOptions()
	defer wo.Destroy()
    fullId := prefix + "-" + id 
	return fullId, s.db.Put(wo, []byte(fullId), data)
}

// Get
func (s *Store) GetItem(id string) ([]byte, error) {
    ro := grocksdb.NewDefaultReadOptions()
    defer ro.Destroy()
    
    rocksSlice, err := s.db.Get(ro, []byte(id))
    if err != nil {
        return nil, err
    }
    defer rocksSlice.Free()

    data := rocksSlice.Data()
    dataCopy := make([]byte, len(data))
    copy(dataCopy, data)

    return dataCopy, nil
}

// GetPrefix
func (s *Store) GetPrefix(prefix string) ([][]byte, error) {
    var datas [][]byte
    ro := grocksdb.NewDefaultReadOptions()
    defer ro.Destroy()
    it := s.db.NewIterator(ro)
    defer it.Close()

    prefixRaw := []byte(prefix)
    for it.Seek(prefixRaw); it.ValidForPrefix(prefixRaw); it.Next() {
        data := it.Value().Data()
        datas = append(datas, data)
    }

    if err := it.Err(); err != nil {
        return nil, err
    }

    return datas, nil
}

// DeleteItem
func (s *Store) DeleteItem(id string) error {
    key := []byte(id)

    // Create a write options object
    wo := grocksdb.NewDefaultWriteOptions()
    defer wo.Destroy()

    // Perform the deletion
    err := s.db.Delete(wo, key)
    if err != nil {
        return fmt.Errorf("failed to delete item with ID %s: %w", id, err)
    }

    return nil
}

// Close closes the store and releases the database resources
func (s *Store) Close() {
	s.db.Close()
}

// Helper function to handle common logic for reading data from RocksDB
func (s *Store) readData(name string) ([]byte, error) {
	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()
	data, err := s.db.Get(ro, []byte(name))
	if err != nil {
		return nil, err
	}
	defer data.Free()
	return data.Data(), nil
}

