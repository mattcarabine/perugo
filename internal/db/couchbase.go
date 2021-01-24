package db

//TODO: Test this works
import (
	"errors"
	"github.com/couchbase/gocb/v2"
	"reflect"
	"time"
)

type CB struct {
	bucket *gocb.Bucket
}

func (cb *CB) initDB() error {
	cluster, err := gocb.Connect(
		"localhost",
		gocb.ClusterOptions{
			Username: "Administrator",
			Password: "password",
		})
	if err != nil {
		return err
	}

	bucket := cluster.Bucket("perugo")
	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(30*time.Second, nil)
	if err != nil {
		return err
	}

	cb.bucket = bucket
	return nil
}

func (cb CB) LookupId(id string, result *interface{}) error {
	docType := reflect.TypeOf(result).Name()
	coll := cb.bucket.Collection(docType)
	doc, err := coll.Get(id, &gocb.GetOptions{})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return ErrDBEntityDoesNotExist
		}
		return ErrDBLookupFailed
	}

	err = doc.Content(result)
	if err != nil {
		return ErrDBLookupFailed
	}
	return nil
}

func (cb CB) Store(id string, val *interface{}) error {
	docType := reflect.TypeOf(val).Name()
	coll := cb.bucket.Collection(docType)
	_, err := coll.Upsert(id, val, &gocb.UpsertOptions{})
	if err != nil {
		return ErrDBStoreFailed
	}
	return nil
}
