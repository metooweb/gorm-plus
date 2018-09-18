package gorm_plus

import (
	"github.com/jinzhu/gorm"
	"database/sql"
	"fmt"
)

type DB struct {
	db        *gorm.DB
	tx        *gorm.DB
	transTime uint
}

func (t *DB) Inst() *gorm.DB {

	if t.tx != nil {
		return t.tx
	} else {
		return t.db
	}
}

func (t *DB) MustBegin() {

	var err error
	t.transTime++
	if t.tx == nil {
		t.tx = t.db.Begin()
		err = t.tx.Error
	} else {
		err = t.tx.Exec(fmt.Sprintf("SAVEPOINT metooweb_trans_%d", t.transTime)).Error
	}
	if err != nil {
		panic(err)
	}
	return
}

func (t *DB) MustCommit() {
	var err error
	if t.tx == nil {
		err = sql.ErrTxDone
		return
	}

	if t.transTime == 1 {
		err = t.tx.Commit().Error
		t.tx = nil
	}

	t.transTime--

	if err != nil {
		panic(err)
	}

	return

}

func (t *DB) Rollback() {

	var err error

	if t.tx == nil {
		err = sql.ErrTxDone
		return
	}

	if t.transTime == 1 {
		t.tx.Rollback()
		t.tx = nil
	} else {
		err = t.tx.Exec(fmt.Sprintf("ROLLBACK TO metooweb_trans_%d", t.transTime)).Error
	}

	t.transTime--
	if err != nil {
		panic(err)
	}

	return
}

func (t *DB) Take(db *gorm.DB, out interface{}, where ...interface{}) (exist bool) {

	var err error

	if err = db.Take(out, where...).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}
		panic(err)
		return
	}

	exist = true

	return
}

func (t *DB) List(db *gorm.DB, out interface{}, where ...interface{}) {

	if err := db.Find(out, where...).Error; err != nil {
		panic(err)
	}

}

func (t *DB) FindAndCount(db *gorm.DB, out interface{}, page int, limit int) (total int64, err error) {

	if err = db.Count(&total).Error; err != nil {
		return
	}

	if err = db.Limit(limit).Offset((page - 1) * limit).Find(out).Error; err != nil {
		return
	}

	return
}

func (t *DB) Exist(db *gorm.DB) (res bool) {

	dest := &struct {
	}{}

	res = t.Take(db, dest)

	return
}

func (t *DB) Get(dst interface{}, sql string, args ...interface{}) (exist bool) {

	var err error

	if err = t.Inst().Raw(sql, args...).Scan(dst).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}
		panic(err)
		return
	}

	exist = true

	return
}

func (t *DB) Exec(sql string, args ...interface{}) {
	var err error
	if err = t.Inst().Exec(sql, args...).Error; err != nil {
		panic(err)
	}
}

func (t *DB) Create(data interface{}) {

	if err := t.Inst().Create(data).Error; err != nil {
		panic(err)
	}

}

func (t *DB) Update(db *gorm.DB, attrs ...interface{}) {

	if err := db.Update(attrs).Error; err != nil {
		panic(err)
	}

}

func (t *DB) Save(val interface{}) {

	if err := t.Inst().Save(val).Error; err != nil {
		panic(err)
	}

}

func NewDB(db *gorm.DB) *DB {

	return &DB{db: db}
}
