// Package pkg provides utils for caching behaviors.
package pkg

import (
	"bytes"
	"encoding/gob"
	"github.com/bluele/gcache"
	"io/ioutil"
	"log"
	"time"
)

// Cache is use to initialize a gCache (https://github.com/bluele/gcache)
// to cache resources from S3 compatible storage.
type Cache struct {
	// Resource larger than the specified max size (bytes) will not be cached.
	MaxSizeCached int64
	// Max number of resources to cache in memory at any one time.
	NumCached int
	// Instance of gCache
	cache gcache.Cache
}

// CachableResource is used as an intermediate struct to be serialized and
// deserialized for caching.
type CachableResource struct {
	// Data to serialize/unserialize.
	Data []byte
	// Metadata for the resource.
	Info ResourceInfo
}

// ToGob converts the Resource into a CachableResource and serialized as a Gob.
func (r *Resource) ToGob() ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	data, err := ioutil.ReadAll(r.Data)
	r.Data = bytes.NewReader(data)
	if err != nil {
		return b.Bytes(), err
	}
	err = e.Encode(CachableResource{Data: data, Info: r.Info})
	return b.Bytes(), err
}

// FromGob converts raw bytes into a CachableResource, which is converted back
// into a Resource.
func (r *Resource) FromGob(data []byte) error {
	d := gob.NewDecoder(bytes.NewReader(data))
	cache := CachableResource{}
	err := d.Decode(&cache)
	if err == nil {
		r.Data = bytes.NewReader(cache.Data)
		r.Info = cache.Info
	}
	return err
}

// NewCache creates a new gCache instance.
func NewCache(NumCached int, MaxSizeCached int64) *Cache {
	gob.Register(CachableResource{})

	cache := gcache.New(NumCached).
		ARC().
		Expiration(time.Hour).
		Build()

	return &Cache{NumCached: NumCached, MaxSizeCached: MaxSizeCached, cache: cache}
}

// GetObjectCache decorates a GetObject function to check the cache before
// actually retrieving the object from the S3 compatible store.
func (h *Cache) GetObjectCache(GetObject func(url string) (Resource, error)) func(url string) (Resource, error) {

	return func(url string) (Resource, error) {

		if h.cache.Has(url) {
			unknown, err := h.cache.Get(url)
			if err == nil {
				log.Printf("loading %s from cache.", url)
				res := Resource{}
				data, ok := unknown.([]byte)
				if ok {
					res.FromGob(data)
					return res, nil
				}
			}
		}

		res, err := GetObject(url)
		if err != nil {
			return Resource{}, err
		}

		// serialize and cache if file is not very big
		if h.MaxSizeCached > 0 && res.Info.Size < h.MaxSizeCached {
			serialized, err := res.ToGob()
			if err == nil {
				h.cache.SetWithExpire(url, serialized, time.Hour)
				log.Printf("saving to %s cache.", url)
			} else {
				log.Printf("Unable to serialize resource[%s]: %s", url, err.Error())
			}
		}
		return res, nil
	}
}
