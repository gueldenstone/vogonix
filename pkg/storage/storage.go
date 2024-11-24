package storage

import (
	"encoding/json"
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

func (s *Storage) UpdateValue(bucket, key string, value []byte) error {
	return s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket '%s' not found", bucket)
		}

		return b.Put([]byte(key), []byte(value))
	})
}

func (s *Storage) GetValue(bucket, key string) ([]byte, error) {
	var value []byte
	err := s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket '%s' not found", bucket)
		}
		value = b.Get([]byte(key))
		return nil
	})
	return value, err
}

func (s *Storage) UpdateStringValue(bucket, key, value string) error {
	return s.UpdateValue(bucket, key, []byte(value))
}

func (s *Storage) UpdateStructuredValue(bucket, key string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable marshal data %v: %w", data, err)
	}
	return s.UpdateValue(bucket, key, jsonData)
}

func (s *Storage) GetStringValue(bucket, key string) (string, error) {
	b, err := s.GetValue(bucket, key)
	return string(b), err
}
func (s *Storage) GetStructuredValue(bucket, key string, data any) error {
	b, err := s.GetValue(bucket, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, data)
}

func (s *Storage) GetAllData(bucket string) (map[string]string, error) {
	m := make(map[string]string)
	err := s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket '%s' not found", bucket)
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			m[string(k)] = string(v)
		}

		return nil
	})
	return m, err
}

func (s *Storage) GetAllKeys(bucket string) ([]string, error) {
	data, err := s.GetAllData(bucket)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	for k, _ := range data {
		keys = append(keys, k)
	}
	return keys, nil
}
