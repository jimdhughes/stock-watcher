package data

import (
	"encoding/json"
	"errors"
	"log"

	"go.etcd.io/bbolt"
)

const (
	ERROR_BUCKET_DOES_NOT_EXIST = "bucket does not exist"
	ERROR_KEY_NOT_FOUND         = "no value for specified key"
	ERROR_UNMARSHALL_OBJECT     = "unable to unmarshal object"
	ERROR_MARSHALL_OBJECT       = "unable to marshal object"
)

type Datastore struct {
	db *bbolt.DB
}

func (d *Datastore) Init(dbname string) {
	conn, err := bbolt.Open(dbname, 0777, nil)
	if err != nil {
		log.Fatal(err)
	}
	d.db = conn
}

// CreateBucket creates a bucket with the name passed by the bucket parameter
// Returns an error if it is unable to create
func (d *Datastore) CreateBucket(bucketName string) (*bbolt.Bucket, error) {
	var bucket *bbolt.Bucket
	err := d.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		bucket = b
		return err
	})
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// GetItemFromBucket looks for an item in a bucket identified by bucketName with a key identified by the key value
// It will unmarshal the byte array into the receiver interface
func (d *Datastore) GetItemFromBucket(bucketName string, key string) (interface{}, error) {
	var rec interface{}
	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return errors.New(ERROR_BUCKET_DOES_NOT_EXIST)
		}
		val := b.Get([]byte(key))
		if val == nil {
			return errors.New(ERROR_KEY_NOT_FOUND)
		}

		if err := json.Unmarshal(val, &rec); err != nil {
			return errors.New(ERROR_UNMARSHALL_OBJECT)
		}
		return nil
	})
	return rec, err
}

// ListItemsInBucket returns a subset of items in a bucket identified by bucketName
// from the start index to the limit count
func (d *Datastore) ListItemsInBucket(bucketName string, start, limit int) ([]interface{}, error) {
	var values []interface{}
	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return errors.New(ERROR_BUCKET_DOES_NOT_EXIST)
		}
		c := b.Cursor()
		i := 0
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if i > limit {
				break
			}
			if i >= start && i <= start+limit {
				var rec interface{}
				err := json.Unmarshal(v, &rec)
				if err != nil {
					log.Println(err)
					return errors.New(ERROR_UNMARSHALL_OBJECT)
				}
				values = append(values, rec)
			}
		}
		return nil
	})
	return values, err
}

// AddItemToBucket adds an item to a particular bucketName with a a key and any value
func (d *Datastore) PutItemInBucket(bucketName string, key string, value interface{}) error {
	err := d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		bytesValue, err := json.Marshal(value)
		if err != nil {
			log.Println(err)
			return errors.New(ERROR_MARSHALL_OBJECT)
		}
		return b.Put([]byte(key), bytesValue)
	})
	return err
}

// DeleteItemFromBucket deletes a value from a bucket identified by bucketName identified by a key
func (d *Datastore) DeleteItemFromBucket(bucketName string, key string) error {
	err := d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Delete([]byte(key))
	})
	return err
}
