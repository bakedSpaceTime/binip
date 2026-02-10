package db

import (
	"fmt"

	"github.com/bakedSpaceTime/binip/libip/config"
	bolt "go.etcd.io/bbolt"
)

var ipRecordsBucket = "ip_records"
var systemBucket = "system"
var version = "0.1.0"

type Db struct {
	Db     *bolt.DB
	dbFile string
}

func New(c *config.Config) *Db {
	db, err := bolt.Open(c.DbFile, 0600, nil)
	if err != nil {
		panic("cannot open db file")
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ipRecordsBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(systemBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		b := tx.Bucket([]byte(systemBucket))
		err = b.Put([]byte("version"), []byte(version))
		err = b.Put([]byte("app_name"), []byte("binip"))

		return err
	})

	return &Db{
		Db: db,
	}
}

func (db *Db) Close() error {
	return db.Db.Close()
}

func (db *Db) Reset() error {
	return db.Db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(ipRecordsBucket))
	})
}

func (db *Db) String() string {
	bs := make(map[string][]string)

	db.Db.View(
		func(tx *bolt.Tx) error {
			tx.ForEach(
				func(name []byte, b *bolt.Bucket) error {
					ns := string(name)
					bs[ns] = []string{}
					b.ForEach(func(k []byte, v []byte) error {
						bs[ns] = append(bs[ns], fmt.Sprintf("(%s: %s)", k, v))
						return nil
					})
					return nil
				})
			return nil
		})
	return fmt.Sprintf("%v", bs)
}
