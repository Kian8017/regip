package regip

import (
	"errors"
	"fmt"
	badger "github.com/dgraph-io/badger/v2"
	"sync"
)

var (
	DBClosed       = errors.New("database closed")
	PushFailed     = errors.New("push failed")
	AlreadyExists  = errors.New("resource already exists")
	ErrKeyNotFound = errors.New("key not found")
)

type DB struct {
	badgerDB *badger.DB
	listen   *DBListener
}

func (d *DB) Close() {
	d.listen.Stop()
	d.listen.Wait()
	d.badgerDB.Close()
}

// NewDB creates a new database
func NewDB(path string) (*DB, error) {
	currentDB, err := badger.Open(badger.DefaultOptions(path).WithLogger(DBLogger{}))
	if err != nil {
		return nil, err
	}
	var newDB DB
	newDB.badgerDB = currentDB
	newDB.listen = NewDBListener()
	return &newDB, nil
}

func (d *DB) AddRaw(rec Resource, update bool) error {
	ok := d.listen.Add()
	if !ok {
		return DBClosed
	}
	defer d.listen.Done()
	rid := rec.ID()

	err := d.badgerDB.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(rid)
		if err == nil && update == false {
			return AlreadyExists
		} else if err != nil && err != badger.ErrKeyNotFound {
			return err // Got some error that wasn't the key's non existence
		}
		err = txn.Set(rid, rec.MarshalBinary())
		return err
	})
	return err
}

func (d *DB) Add(rec Resource) error {
	return d.AddRaw(rec, false)
}

func (d *DB) Update(rec Resource) error {
	return d.AddRaw(rec, true)
}

func (d *DB) Delete(i ID) error {
	ok := d.listen.Add()
	if !ok {
		return DBClosed
	}
	defer d.listen.Done()
	err := d.badgerDB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(i))
	})
	return err
}

func (d *DB) Get(i ID, lg *Logger) (Resource, error) {
	lgr := lg.Tag("db.Get", CLR_db)
	ok := d.listen.Add()
	if !ok {
		lgr.Error("DB closed")
		return nil, DBClosed
	}
	defer d.listen.Done()

	var rec Resource
	err := d.badgerDB.View(func(txn *badger.Txn) error {
		raw, err := txn.Get(i)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return ErrKeyNotFound
			} else {
				lgr.Error("Got other error txn.Get-ing ", err)
				return err
			}
		}
		err = raw.Value(func(val []byte) error {
			// FIXME: DEBUG
			rec, err = Unmarshal(i[0], val)
			if err != nil {
				lgr.Error("Error unmarshaling value ", string(val),
					" with length ", len(val),
					" with error ", err,
					" with key ", i)
			}
			return err
		})
		return err
	})
	return rec, err
}

func (d *DB) Exists(i ID) bool {
	ok := d.listen.Add()
	if !ok {
		return false
	}
	defer d.listen.Done()

	err := d.badgerDB.View(func(txn *badger.Txn) error {
		_, err := txn.Get(i)
		if err == nil {
			return nil
		} else {
			return ErrKeyNotFound
		}
	})

	if err == nil {
		return true
	} else {
		return false
	}
}

func (d *DB) Flow(prefix byte) *Flow {
	ok := d.listen.Add()
	if !ok {
		return nil
	}

	nf := NewFlow()
	go func() {
		_ = d.badgerDB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var err error
			for it.Seek([]byte{prefix}); it.ValidForPrefix([]byte{prefix}); it.Next() {
				item := it.Item()
				err = item.Value(func(val []byte) error {
					// Make a copy of the buffer, so that subsequent iterations don't affect the previous
					curBuf := make([]byte, len(val))
					copy(curBuf, val)
					current, err := Unmarshal(prefix, curBuf)
					if err != nil {
						return err
					}

					ok := nf.Push(current)
					if !ok {
						return PushFailed
					} else {
						return nil
					}
				})
				if err != nil {
					break
				}
			}
			return err
		})
		nf.Stop()
		d.listen.Done()
	}()
	return nf
}

func (d *DB) IDFlow(prefix byte, lgr *Logger) *Flow {
	ok := d.listen.Add()
	if !ok {
		return nil
	}

	nf := NewFlow()
	go func() {
		_ = d.badgerDB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var err error
			for it.Seek([]byte{prefix}); it.ValidForPrefix([]byte{prefix}); it.Next() {
				k := it.Item().Key()
				if len(k) != ID_LENGTH {
					lgr.Print("Found key entry of invalid length")
					continue
				}
				kc := make([]byte, len(k))
				copy(kc, k)
				i := ID(kc)
				// DEBUG
				lgr.Print("Adding ID ", i.String())
				ok := nf.Push(i)
				if !ok {
					lgr.Error("Push failed, exiting...")
					return PushFailed
				}
			}
			return err
		})
		nf.Stop()
		d.listen.Done()
	}()
	return nf
}

func Unmarshal(t byte, raw []byte) (Resource, error) {
	switch t {
	case RT_record:
		return UnmarshalRecordBinary(raw)
	case RT_country:
		return UnmarshalCountryBinary(raw)
	case RT_user:
		return UnmarshalUserBinary(raw)
	case RT_indexRecord:
		return UnmarshalIndexRecordBinary(raw)
	case RT_trigram:
		return UnmarshalTrigramBinary(raw)
	default:
		panic(fmt.Sprint("Implement db.Get case for resource", t))
	}
}

type DBLogger struct {
}

func (L DBLogger) Errorf(s string, a ...interface{}) {
	fmt.Printf(s, a...)
}

func (L DBLogger) Warningf(s string, a ...interface{}) {
	fmt.Printf(s, a...)
}

func (L DBLogger) Infof(s string, a ...interface{}) {
}

func (L DBLogger) Debugf(s string, a ...interface{}) {
}

type DBListener struct {
	listen sync.WaitGroup
	stop   bool
	mut    sync.Mutex
}

func NewDBListener() *DBListener {
	var l DBListener
	l.listen = sync.WaitGroup{}
	l.mut = sync.Mutex{}
	return &l
}

func (l *DBListener) Add() bool {
	if !l.ok() {
		return false
	}
	l.listen.Add(1)
	return true
}

func (l *DBListener) Done() {
	l.listen.Done()
}

func (l *DBListener) Stop() {
	l.mut.Lock()
	defer l.mut.Unlock()
	l.stop = true
}

func (l *DBListener) Wait() {
	l.listen.Wait()
}

func (l *DBListener) ok() bool {
	l.mut.Lock()
	defer l.mut.Unlock()
	return !l.stop
}
