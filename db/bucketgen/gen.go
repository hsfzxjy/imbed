package main

import (
	"go/format"
	"os"
	"path"
	"strconv"
	"strings"
)

var buckets = [...]string{
	// Notations:
	// _ means concat or hash
	// __ separates the key and value
	"FILES",
	"f_meta",
	"F_SHA__OID",
	"F_FHASH_OID",
	"F_BASENAME_OID",
	"F_URL_OID",
	"F_FHASH_TSSHA__OID",
	"T_FOID_TAG",
	"T_TAG__FOID",
	"T_COID_FOID",
	"CONFIGS",
	"C_SHA__OID",
}

func main() {
	var w strings.Builder
	w.WriteString(`
	// Code generated by "github.com/hsfzxjy/imbed/db/bucketgen"; DO NOT EDIT.

	package db

	import "go.etcd.io/bbolt"

	var bucketNames = [...][]byte{
	`)
	for _, name := range buckets {
		w.WriteString("[]byte(" + strconv.Quote(name) + "),\n")
	}

	w.WriteString(`}

	func (tx*Tx) createAllBuckets() error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name);err != nil{
				return err
			}
		}
		return nil
	}
	`)

	for i, name := range buckets {
		w.WriteString(`
		func (tx*Tx) ` + name + `() *bbolt.Bucket {
			slot := &tx.buckets[` + strconv.Itoa(i) + `]
			slot.Do(func() {
				b := tx.Bucket(bucketNames[` + strconv.Itoa(i) + `])
				if b == nil {
					panic("fatal: bucket ` + name + ` not found, database corrupted")
				}
				slot.Bucket = b
			})
			return slot.Bucket
		}
		`)
	}

	output, err := format.Source([]byte(w.String()))
	check(err)

	filename := os.Getenv("GOFILE")
	dir := path.Dir(filename)
	check(os.WriteFile(path.Join(dir, "db_bucket_gen.go"), output, 0o644))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
