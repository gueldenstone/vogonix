package storage

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type Storage struct {
	*bolt.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// create buckets if they not exists

	return &Storage{
		db,
	}, nil
}

func (s *Storage) AddBucket(bucket string) error {
	return s.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
}

func (s *Storage) UpdateValue(bucket, key, value string) error {
	return s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket '%s' not found", bucket)
		}

		return b.Put([]byte(key), []byte(value))
	})
}

func (s *Storage) GetValue(bucket, key string) (string, error) {
	var value string = ""
	err := s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket '%s' not found", bucket)
		}
		value = string(b.Get([]byte(key)))
		return nil
	})
	return value, err
}
