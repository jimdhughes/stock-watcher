package data

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"
)

const (
	TEST_DB_NAME     = "testdb.bolt"
	TEST_BUCKET_NAME = "bucket1"
)

func TestDbCreateBucket(t *testing.T) {
	defer teardown()
	Init(TEST_DB_NAME)
	_, err := datastore.CreateBucket(TEST_BUCKET_NAME)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestCreateItemInBucket(t *testing.T) {
	setup()
	defer teardown()
	err := datastore.PutItemInBucket(TEST_BUCKET_NAME, "1", "hi")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestGetItemFromBucket(t *testing.T) {
	setup()
	defer teardown()
	err := datastore.PutItemInBucket(TEST_BUCKET_NAME, "1", "hi")
	if err != nil {
		t.Fatalf(err.Error())
	}
	val, err := datastore.GetItemFromBucket(TEST_BUCKET_NAME, "1")
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
	err := datastore.PutItemInBucket(TEST_BUCKET_NAME, "1", "hi")
	if err != nil {
		t.Fatalf(err.Error())
	}
	items, err := datastore.ListItemsInBucket(TEST_BUCKET_NAME, 0, 1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(items) != 1 {
		t.Fatalf("Expceted 1 item. got 0")
	}
}

func TestListMoreItemsThanInBucket(t *testing.T) {
	setup()
	defer teardown()
	err := datastore.PutItemInBucket(TEST_BUCKET_NAME, "1", "hi")
	if err != nil {
		t.Fatalf(err.Error())
	}
	items, err := datastore.ListItemsInBucket(TEST_BUCKET_NAME, 0, 10)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(items) != 1 {
		t.Fatalf("expect exactly 1 item")
	}
}

func TestListItemsStartingAtIndexLargerThanBucketSize(t *testing.T) {
	setup()
	defer teardown()
	err := datastore.PutItemInBucket(TEST_BUCKET_NAME, "1", "hi")
	if err != nil {
		t.Fatalf(err.Error())
	}
	items, err := datastore.ListItemsInBucket(TEST_BUCKET_NAME, 5, 10)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(items) != 0 {
		t.Fatalf("expect exactly 1 item")
	}
}

func TestDeleteItemFromBucket(t *testing.T) {
	setup()
	defer teardown()
	err := datastore.PutItemInBucket(TEST_BUCKET_NAME, "1", "hi")
	if err != nil {
		t.Fatalf(err.Error())
	}
	items, err := datastore.ListItemsInBucket(TEST_BUCKET_NAME, 0, 1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(items) != 1 {
		t.Fatalf("expect exactly 1 item")
	}
	err = datastore.DeleteItemFromBucket(TEST_BUCKET_NAME, "1")
	if err != nil {
		t.Fatalf(err.Error())
	}
	items, err = datastore.ListItemsInBucket(TEST_BUCKET_NAME, 0, 1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(items) > 0 {
		t.Fatalf("expected to now have 0 items")
	}
}

func setup() {
	execpath, _ := os.Executable()
	directory := path.Dir(execpath)
	Init(fmt.Sprintf("%s/%s", directory, TEST_DB_NAME))
	datastore.CreateBucket(TEST_BUCKET_NAME)
}

func teardown() {
	execpath, _ := os.Executable()
	directory := path.Dir(execpath)
	datastore.db.Close()
	err := os.Remove(fmt.Sprintf("%s/%s", directory, TEST_DB_NAME))
	if err != nil {
		log.Println("Error tearingdown")
	}
}
