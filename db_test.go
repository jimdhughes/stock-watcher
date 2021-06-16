package main

import (
	"log"
	"os"
	"testing"
)

const (
	TEST_DB_NAME     = "testdb.bolt"
	TEST_BUCKET_NAME = "bucket1"
)

var d Datastore

func TestDbCreateBucket(t *testing.T) {
	d = Datastore{}
	defer teardown()
	d.Init(TEST_DB_NAME)
	_, err := d.CreateBucket(TEST_BUCKET_NAME)
	if err != nil {
		t.Fatalf(err.Error())
	}

}

func TestCreateItemInBucket(t *testing.T) {
	setup()
	defer teardown()
	err := d.AddItemToBucket(TEST_BUCKET_NAME, "1", "hi")
	if err != nil {
		t.Fatalf(err.Error())
	}

}

func TestGetItemFromBucket(t *testing.T) {
	setup()
	defer teardown()
	err := d.AddItemToBucket(TEST_BUCKET_NAME, "1", "hi")
	if err != nil {
		t.Fatalf(err.Error())
	}
	val, err := d.GetItemFromBucket(TEST_BUCKET_NAME, "1")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if val != "hi" {
		t.Fatalf("expected %s - got %s", "hi", val)
	}

}

func TestListBucketItems(t *testing.T) {
	setup()
	defer teardown()
	err := d.AddItemToBucket(TEST_BUCKET_NAME, "1", "hi")
	if err != nil {
		t.Fatalf(err.Error())
	}
	items, err := d.ListItemsInBucket(TEST_BUCKET_NAME, 0, 1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(items) != 1 {
		t.Fatalf("Expceted 1 item. got 0")
	}

}

func setup() {
	d = Datastore{}
	d.Init(TEST_DB_NAME)
	d.CreateBucket(TEST_BUCKET_NAME)
}

func teardown() {
	d.db.Close()
	err := os.Remove(TEST_DB_NAME)
	if err != nil {
		log.Println("Error tearingdown")
	}
}
