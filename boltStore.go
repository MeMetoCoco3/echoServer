package main

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
)

type BoltStore[K comparable, V any] struct {
	db         *bolt.DB
	bucketName []byte
}

type Storer[K comparable, V any] interface {
	Put(K, V) error
	Get(K) (V, error)
	Update(K, V) error
	Delete(K) error
	GetAll() (map[K]V, error)
}

var bucketUsersName = "users_collection"
var bucketEmailName = "email_collection"

const pathUsers = "./users.db"
const pathEmailIndex = "./emailIdx.db"

func NewBoltStore[K comparable, V any](path, bucketName string) (*BoltStore[K, V], error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &BoltStore[K, V]{
		db:         db,
		bucketName: []byte(bucketName),
	}, nil
}

func (s *BoltStore[K, V]) Put(key K, value V) error {
	keyBytes, err := json.Marshal(key)
	if err != nil {
		return fmt.Errorf("Key '%v'  was not correctly converted into json. %v", key, err)
	}
	fmt.Println(keyBytes)
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("Value '%v'  was not correctly converted into json. %v", value, err)

	}

	if err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		return bucket.Put(keyBytes, valueBytes)

	}); err != nil {
		return fmt.Errorf("Error during the transaction:%v", err)
	}

	return nil
}

func (s *BoltStore[K, V]) Get(key K) (V, error) {
	var v V
	var valBytes []byte
	keyBytes, err := json.Marshal(key)
	if err != nil {
		return v, err
	}

	if err = s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		valBytes = bucket.Get(keyBytes)
		if valBytes == nil {
			return fmt.Errorf("Key '%v' is not present in the Store.", key)
		}
		return nil
	}); err != nil {
		return v, err
	}

	err = json.Unmarshal(valBytes, &v)
	if err != nil {
		return v, fmt.Errorf("Err1or on get transaction. %s, %s", err, valBytes)
	}

	return v, nil
}

func (s *BoltStore[K, V]) GetAll() (map[K]V, error) {
	result := make(map[K]V)

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", s.bucketName)
		}

		return bucket.ForEach(func(k, v []byte) error {
			var key K

			switch any(key).(type) {
			case string:
				key = any(string(k)).(K)
			default:
				err := json.Unmarshal(k, &key)
				if err != nil {
					return err
				}
			}

			var value V
			err := json.Unmarshal(v, &value)
			if err != nil {
				return err
			}

			result[key] = value
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *BoltStore[K, V]) Update(key K, value V) error {
	return s.Put(key, value)
}
func (s *BoltStore[K, V]) Delete(key K) error {
	keyBytes, err := json.Marshal(key)
	if err != nil {
		return fmt.Errorf("Key '%v'  was not correctly converted into json. %v", key, err)
	}

	if err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		return bucket.Delete(keyBytes)
	}); err != nil {
		return err
	}

	return nil
}

func (s *BoltStore[K, V]) GetValues(key K, values []K) (map[K]string, error) {
	keyBytes, err := json.Marshal(key)
	if err != nil {
		return nil, fmt.Errorf("GetValues Error while marshal key: %v", err)
	}
	requestedFields := make(map[K]string)

	if err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		userData := bucket.Get(keyBytes)
		if userData == nil {
			return fmt.Errorf("Value with key: %v not found on Database.", key)
		}
		var u User

		err	:= json.Unmarshal(userData, &u)
		if err != nil {
			return fmt.Errorf("GetValues Error while unmarshal Userdata on User %v: %v", userData, err)

		}
		requestedFields = FilterStruct[K, V](u, values)
		return nil

	}); err != nil {
		return nil, err
	}

	return requestedFields, nil
}

// TODO
//func SwitchToBytes[K comparable](k K) []byte {
//	var keyBytes []byte
//	switch any(k).(type) {
//	case string:
//		ky= any(string(k)).(K)
//	default:
//		err := json.Unmarshal(k, &keyBytes)
//		if err != nil {
//			return err
//		}
//	}
//}
