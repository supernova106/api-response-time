package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"mime/multipart"
	"os"
)

const (
	// MongoDBUrl is the default mongodb url that will be used to connect to the
	// database.
	mongoDBUrl = "mongodb://10.0.1.105:27017/fluentd"
)

func Connect() (*DB, error) {
	host := os.Getenv("MONGODB_URL")

	if len(host) == 0 {
		host = mongoDBUrl
	}

	session, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)
	db := &DB{session: session}
	err = db.ensureDefaultIndex()

	return db, err
}

type Query bson.M

type DB struct {
	Name    string
	session *mgo.Session
}

// Copy creates a new socket with the same authentication information.
// DB.Close() must be called after use.
// See mgo.Session.Copy()
func (db *DB) Copy() *DB {
	session := db.session.Copy()

	return &DB{
		Name:    db.Name,
		session: session,
	}
}

func (db *DB) Pipe(name string, q interface{}, v interface{}) (err error) {
	err = db.DB().C(name).Pipe(q).All(v)
	return
}

func (db *DB) Close() {
	db.session.Close()
}

func (db *DB) DB() *mgo.Database {
	return db.session.DB(db.Name)
}

func (db *DB) C(name string) *mgo.Collection {
	return db.DB().C(name)
}

func (db *DB) DropC(name string) error {
	return db.C(name).DropCollection()
}

func (db *DB) FindOne(name string, q interface{}, v interface{}) (err error) {
	err = db.DB().C(name).Find(q).One(v)
	return
}

func (db *DB) Count(name string, q interface{}) (count int, err error) {
	count, err = db.DB().C(name).Find(q).Count()

	return
}

func (db *DB) FindID(name string, id string, v interface{}) (err error) {
	objectID := bson.ObjectIdHex(id)
	err = db.DB().C(name).FindId(objectID).One(v)

	return
}

func (db *DB) OpenFile(name string, fileName string) (*mgo.GridFile, error) {
	return db.DB().GridFS(name).Open(fileName)
}

func (db *DB) FindAllPaginated(name string, q interface{}, sortFields []string, v interface{}, page, skip int) (err error) {
	err = db.DB().C(name).Find(q).Skip(page).Limit(skip).Sort(sortFields...).All(v)

	return
}

func (db *DB) FindAll(name string, q interface{}, sortFields []string, v interface{}) (err error) {
	err = db.DB().C(name).Find(q).Sort(sortFields...).All(v)

	return
}

func (db *DB) UpsertSession(name string, q interface{}, v interface{}) (err error) {
	_, err = db.DB().C(name).Upsert(q, v)
	return
}

func (db *DB) Update(name string, q interface{}, v interface{}) (err error) {
	err = db.DB().C(name).Update(q, v)
	return
}

func (db *DB) Insert(name string, v interface{}) (err error) {
	err = db.DB().C(name).Insert(v)
	return
}

func (db *DB) RemoveCollection(name string, q interface{}) (removed int, err error) {
	i, err := db.DB().C(name).RemoveAll(q)
	return i.Removed, err
}

func (db *DB) Upsert(name string, q interface{}, f Query, v interface{}, forceUpdate bool) (updated bool, err error) {
	var findResult interface{}
	updated = false

	if forceUpdate == true {
		_, err = db.DB().C(name).Upsert(q, v)
		updated = true
		return
	}

	err = db.DB().C(name).Find(f).One(&findResult)

	if err == mgo.ErrNotFound {
		err = nil
		_, err = db.DB().C(name).Upsert(q, v)

		if err == nil {
			updated = true
		}
	}

	return
}

func (db *DB) RemoveFile(col, path string) error {
	gridfs := db.DB().GridFS(col)
	return gridfs.Remove(path)
}

func (db *DB) UpsertFile(col, path string, file multipart.File) (id bson.ObjectId, err error) {
	gridfs := db.DB().GridFS(col)
	err = gridfs.Remove(path)
	if err != nil {
		return
	}

	gridFile, err := gridfs.Create(path)
	if err != nil {
		return
	}
	defer gridFile.Close()

	_, err = io.Copy(gridFile, file)
	if err != nil {
		return
	}

	if objectID, ok := gridFile.Id().(bson.ObjectId); ok {
		id = objectID
	} else {
		err = fmt.Errorf("Grid File ID is not available")
	}

	return
}

func (db *DB) RemoveID(name string, id bson.ObjectId) (err error) {
	err = db.DB().C(name).RemoveId(id)

	return
}

func (db *DB) RemoveAll(name string, q interface{}) (err error) {
	err = db.DB().C(name).Remove(q)

	return
}

func (db *DB) EnsureIndexKey(colName string, keys ...string) error {
	for _, key := range keys {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			DropDups:   true,
			Background: true,
		}

		err := db.C(colName).EnsureIndex(index)

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) EnsureUniqueIndexKey(colName string, keys ...string) error {
	for _, key := range keys {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     true,
			DropDups:   true,
			Background: true,
		}

		err := db.C(colName).EnsureIndex(index)
		if err != nil {
			return err
		}

	}

	return nil
}

func (db *DB) ensureDefaultIndex() error {

	err := db.EnsureUniqueIndexKey("users", "userid")
	if err != nil {
		return err
	}

	err = db.EnsureUniqueIndexKey("sessions", "sessionid")
	if err != nil {
		return err
	}

	err = db.EnsureIndexKey("sessions", "lastupdatedtime")
	if err != nil {
		return err
	}

	return nil
}
