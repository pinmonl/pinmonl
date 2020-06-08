package passwd

import (
	"golang.org/x/crypto/bcrypt"
)

var (
	DefaultHasher = NewBcryptHasher(bcrypt.DefaultCost)
)

func Hash(bs []byte) ([]byte, error) {
	return DefaultHasher.Hash(bs)
}

func HashString(str string) (string, error) {
	return DefaultHasher.HashString(str)
}

func Compare(hashed, source []byte) error {
	return DefaultHasher.Compare(hashed, source)
}

func CompareString(hashed, source string) error {
	return DefaultHasher.CompareString(hashed, source)
}

type BcryptHasher struct {
	Cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	bh := &BcryptHasher{
		Cost: cost,
	}
	return bh
}

func (bh *BcryptHasher) Hash(bs []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(bs, bh.Cost)
}

func (bh *BcryptHasher) HashString(str string) (string, error) {
	hashed, err := bh.Hash([]byte(str))
	return string(hashed), err
}

func (bh *BcryptHasher) Compare(hashed, source []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, source)
}

func (bh *BcryptHasher) CompareString(hashed, source string) error {
	return bh.Compare([]byte(hashed), []byte(source))
}
