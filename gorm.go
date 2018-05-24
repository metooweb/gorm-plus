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

func (t *DB) Take(db *gorm.DB, out interface{}, where ...interface{}) (exist bool, err error) {

	if err = db.Take(out, where...).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return
	}

	exist = true

	return
}

func NewDB(db *gorm.DB) *DB {

	return &DB{db: db}
}
