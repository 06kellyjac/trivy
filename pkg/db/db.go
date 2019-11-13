package db

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/simar7/gokv/encoding"

	"github.com/simar7/gokv/types"

	"github.com/aquasecurity/trivy/pkg/log"

	"golang.org/x/xerrors"

	"github.com/aquasecurity/trivy/pkg/utils"

	bolt "github.com/etcd-io/bbolt"
	gokvbolt "github.com/simar7/gokv/bbolt"
)

const (
	SchemaVersion = 1
)

var (
	db    *bolt.DB
	kv    *gokvbolt.Store
	dbDir string
)

//type OperationsV2 interface {
//	Get(input types.GetItemInput) (bool, error)
//	Set(input types.SetItemInput) error
//	BatchSet(input types.BatchSetItemInput) error
//}

//type DBv2 struct {
//	gokv *gokvbolt.Store
//}

// TODO: Should we move this constructor to gokv lib?
func NewDBv2() (*gokvbolt.Store, error) {
	return NewDBv2WithOptions(gokvbolt.Options{
		RootBucketName: "trivy",
		Path:           dbDir + "trivydb2.db",
		Codec:          encoding.JSON,
	})
}

// TODO: Should we move this constructor to gokv lib?
func NewDBv2WithOptions(options gokvbolt.Options) (*gokvbolt.Store, error) {
	log.Logger.Debug(">>>53")
	configOptions := gokvbolt.Options{}
	if options.RootBucketName != "" {
		configOptions.RootBucketName = options.RootBucketName
	} else {
		configOptions.RootBucketName = "trivy"
	}

	if options.Path != "" {
		configOptions.Path = options.Path
	} else {
		dbDir = filepath.Join(utils.CacheDir(), "db")
		if err := os.MkdirAll(dbDir, 0700); err != nil {
			return nil, xerrors.Errorf("failed to mkdir: %w", err)
		}
		configOptions.Path = dbDir + "trivydb2.db"
	}
	log.Logger.Debug(">>>70db")

	if options.Codec != nil {
		configOptions.Codec = options.Codec
	} else {
		configOptions.Codec = encoding.JSON
	}

	var err error
	// this call blocks.... can only be called once
	kv, err = gokvbolt.NewStore(configOptions)
	if err != nil {
		return nil, err
	}

	log.Logger.Debug(">>>84")
	return kv, nil
}

//func (dbc DBv2) Get(input types.GetItemInput) (bool, error) {
//	return dbc.gokv.Get(input)
//}
//
//func (dbc DBv2) Set(input types.SetItemInput) error {
//	return dbc.gokv.Set(input)
//}
//
//func (dbc DBv2) BatchSet(input types.BatchSetItemInput) error {
//	return dbc.gokv.BatchSet(input)
//}

type Operations interface {
	SetVersion(int) error
	Update(string, string, string, interface{}) error
	BatchUpdate(func(*bolt.Tx) error) error
	PutNestedBucket(*bolt.Tx, string, string, string, interface{}) error
	ForEach(string, string) (map[string][]byte, error)
}

type Config struct {
}

func Init() error {
	dbDir = filepath.Join(utils.CacheDir(), "db")
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return xerrors.Errorf("failed to mkdir: %w", err)
	}

	dbPath := filepath.Join(dbDir, "trivy.db")
	log.Logger.Debugf("db path: %s", dbPath)

	//var err error
	//db, err = bolt.Open(dbPath, 0600, nil)
	//if err != nil {
	//	return xerrors.Errorf("failed to open db: %w", err)
	//}

	var err error
	kv, err = NewDBv2()
	//kv, err = gokvbolt.NewStore(gokvbolt.Options{
	//	RootBucketName: "trivy",
	//	Path:           dbPath,
	//	Codec:          encoding.JSON,
	//})
	if err != nil {
		return xerrors.Errorf("failed to open kv store: %w", err)
	}

	log.Logger.Debug("kv: ", kv)
	return nil

	//gokv, err = gokvbolt.NewStore(gokvbolt.Options{
	//	RootBucketName: "trivy",
	//	Path:           dbPath,
	//	Codec:          encoding.JSON,
	//})
	//if err != nil {
	//	return xerrors.Errorf("failed to initialize gokv: %w", err)
	//}
	//return nil
}

func Close() error {
	//if err := db.Close(); err != nil {
	//	return xerrors.Errorf("failed to close DB: %w", err)
	//}

	if err := kv.Close(); err != nil {
		return xerrors.Errorf("failed to close DB: %w", err)
	}
	return nil
}

func Reset() error {
	if err := Close(); err != nil {
		return xerrors.Errorf("failed to reset DB: %w", err)
	}

	if err := os.RemoveAll(dbDir); err != nil {
		return xerrors.Errorf("failed to reset DB: %w", err)
	}

	if err := Init(); err != nil {
		return xerrors.Errorf("failed to reset DB: %w", err)
	}
	return nil
}

func GetVersion() int {
	//var version int
	//found, err := gokv.Get(types.GetItemInput{
	//	BucketName: "metadata",
	//	Key:        "version",
	//	Value:      version,
	//})
	//switch {
	//case err != nil:
	//	// old trivy version
	//	return 1
	//case !found:
	//	// initial run
	//	return 0
	//default:
	//	return version
	//}

	//value, err := Get("trivy", "metadata", "version")

	value, err := Get("metadata", "version")
	if err != nil || len(value) == 0 {
		// initial run
		return 0
	}

	version, err := strconv.Atoi(string(value))
	if err != nil {
		// old trivy version
		return 1
	}
	return version
}

func (dbc Config) SetVersion(version int) error {
	//if err := gokv.Set(types.SetItemInput{
	//	Key:   "metadata",
	//	Value: version,
	//}); err != nil {
	//	return xerrors.Errorf("failed to save DB version: %w", err)
	//}
	//
	//return nil

	if err := kv.Set(types.SetItemInput{
		BucketName: "metadata",
		Key:        "version",
		Value:      version,
	}); err != nil {
		return xerrors.Errorf("failed to save DB version: %w", err)
	}

	//err := dbc.Update("trivy", "metadata", "version", version)
	//if err != nil {
	//	return xerrors.Errorf("failed to save DB version: %w", err)
	//}
	return nil
}

func (dbc Config) Update(rootBucket, nestedBucket, key string, value interface{}) error {
	err := db.Update(func(tx *bolt.Tx) error {
		return dbc.PutNestedBucket(tx, rootBucket, nestedBucket, key, value)
	})
	if err != nil {
		return xerrors.Errorf("error in db update: %w", err)
	}
	return err
}

func (dbc Config) PutNestedBucket(tx *bolt.Tx, rootBucket, nestedBucket, key string, value interface{}) error {
	root, err := tx.CreateBucketIfNotExists([]byte(rootBucket))
	if err != nil {
		return xerrors.Errorf("failed to create a bucket: %w", err)
	}
	return Put(root, nestedBucket, key, value)
}

func Put(root *bolt.Bucket, nestedBucket, key string, value interface{}) error {
	return kv.Set(types.SetItemInput{
		BucketName: nestedBucket,
		Key:        key,
		Value:      value,
	})

	//nested, err := root.CreateBucketIfNotExists([]byte(nestedBucket))
	//if err != nil {
	//	return xerrors.Errorf("failed to create a bucket: %w", err)
	//}
	//v, err := json.Marshal(value)
	//if err != nil {
	//	return xerrors.Errorf("failed to unmarshal JSON: %w", err)
	//}
	//return nested.Put([]byte(key), v)
}

func (dbc Config) BatchUpdate(fn func(tx *bolt.Tx) error) error {
	//kv.BatchSet()

	err := db.Batch(fn)
	if err != nil {
		return xerrors.Errorf("error in batch update: %w", err)
	}
	return nil
}

func Get(nestedBucket, key string) (value []byte, err error) {
	var retVal []byte
	if _, err := kv.Get(types.GetItemInput{
		BucketName: nestedBucket,
		Key:        key,
		Value:      retVal,
	}); err != nil {
		return nil, xerrors.Errorf("failed to get data from db: %w", err)
	}

	return retVal, nil

	//err = db.View(func(tx *bolt.Tx) error {
	//	root := tx.Bucket([]byte(rootBucket))
	//	if root == nil {
	//		return nil
	//	}
	//	nested := root.Bucket([]byte(nestedBucket))
	//	if nested == nil {
	//		return nil
	//	}
	//	value = nested.Get([]byte(key))
	//	return nil
	//})
	//if err != nil {
	//	return nil, xerrors.Errorf("failed to get data from db: %w", err)
	//}
	//return value, nil
}

func (dbc Config) ForEach(rootBucket, nestedBucket string) (value map[string][]byte, err error) {
	value = map[string][]byte{}
	err = db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(rootBucket))
		if root == nil {
			return nil
		}
		nested := root.Bucket([]byte(nestedBucket))
		if nested == nil {
			return nil
		}
		err := nested.ForEach(func(k, v []byte) error {
			value[string(k)] = v
			return nil
		})
		if err != nil {
			return xerrors.Errorf("error in db foreach: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf("failed to get all key/value in the specified bucket: %w", err)
	}
	return value, nil
}
