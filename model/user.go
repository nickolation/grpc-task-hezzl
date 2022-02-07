package model

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/gob"
	"os"
	"time"
)

type User struct {
	Id int `db:"id"`

	Username     string     `db:"username"`
	Password     []byte     `db:"password_hash"`
	Gender       string     `db:"gender"`
	Age          int        `db:"age"`
	Description  string     `db:"description"`
	
	Hash         []byte     `db:"user_hash"`
	Date         *time.Time `db:"start_date"`
}

func (u *User) toBytes() []byte {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(u)
	return buf.Bytes()
}

func (u *User) getUserHash() []byte {
	hash := sha256.Sum256(u.toBytes())
	return hash[:]
}

func (u *User) getPasswordHash() []byte {
	hash := sha1.New()
	hash.Write([]byte(u.Password))

	return hash.Sum([]byte(os.Getenv("SALT")))
}

func (u *User) SetHashUserValues() {
	u.Hash = u.getUserHash()
	u.Password = u.getPasswordHash()
}
