package ext

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"time"

	"github.com/bluele/gcache"
	core "github.com/e2fyi/minio-web/pkg/core"
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

// CacheRequestsExtension installs the extension to cache all GetObject requests to
// the S3 compatible backend.
func CacheRequestsExtension(NumCached int, MaxSizeCached int64) Extension {
	return func(c *Core) (string, error) {
		cacher := NewCache(NumCached, MaxSizeCached)
		c.ApplyGetObject(cacher.getObjectCache)
		return "caching: enabled", nil
	}
}

// NewCache creates a new Cache object.
func NewCache(NumCached int, MaxSizeCached int64) *Cache {
	gob.Register(CachableResource{})

	cache := gcache.New(NumCached).
		ARC().
		Expiration(5 * time.Minute).
		Build()

	return &Cache{NumCached: NumCached, MaxSizeCached: MaxSizeCached, cache: cache}
}

// toGob converts the Resource into a CachableResource and serialized as a Gob.
func toGob(r *core.Resource) ([]byte, error) {
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

// fromGob converts raw bytes into a CachableResource, which is converted back
// into a Resource.
func fromGob(r *core.Resource, data []byte) error {
	d := gob.NewDecoder(bytes.NewReader(data))
	cache := CachableResource{}
	err := d.Decode(&cache)
	if err == nil {
		r.Data = bytes.NewReader(cache.Data)
		r.Info = cache.Info
	}
	return err
}

// getObjectCache decorates a Handler function to check the cache before
// actually retrieving the object from the S3 compatible store.
func (h *Cache) getObjectCache(GetObject Handler) Handler {

	return func(url string) (Resource, error) {

		if h.cache.Has(url) {
			unknown, err := h.cache.Get(url)
			if err == nil {
				log.Printf("loading %s from cache.", url)
				res := Resource{}
				data, ok := unknown.([]byte)
				if ok {
					fromGob(&res, data)
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
			serialized, err := toGob(&res)
			if err == nil {
				h.cache.SetWithExpire(url, serialized, 5*time.Minute)
				log.Printf("saving to %s cache.", url)
			} else {
				log.Printf("Unable to serialize resource[%s]: %s", url, err.Error())
			}
		}
		return res, nil
	}
}
