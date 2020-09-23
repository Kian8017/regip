package regip

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

// USER

const (
	PASS_SALT_LENGTH   = 8
	PASS_BCRYPT_LENGTH = 60
)

type User struct {
	Username string `json:"user"`
	Hash     []byte `json:"hash"`
	Salt     []byte `json:"salt"`
	MetaID
}

func NewUser(username, password string) *User {
	// FIXME: Add error checking code in case they input a username that's too long (over 255 characters)
	var u User
	u.Username = username

	hashed, salted, err := GenerateHash(password)
	if err != nil {
		return nil
	}
	u.Hash = hashed
	u.Salt = salted
	return &u
}

func (u *User) String() string {
	return fmt.Sprintf("User{ID:(%s),salt:(%s),hash:(%s)}|%s", u.ID().String(), string(u.Salt), string(u.Hash), u.Username)
}

func (u *User) MarshalBinary() []byte {
	// (Username length + username) + (Salt length + salt)  +  (hash)
	length := (1 + len(u.Username)) + (1 + PASS_SALT_LENGTH) + PASS_BCRYPT_LENGTH
	buf := make([]byte, length)

	// Add user name
	buf[0] = byte(len(u.Username))
	copy(buf[1:], u.Username)

	// Add salt
	// FIXME
	saltOffset := len(u.Username) + 1
	buf[saltOffset] = byte(PASS_SALT_LENGTH)
	copy(buf[saltOffset+1:], u.Salt)

	// Add Hash
	hashOffset := saltOffset + PASS_SALT_LENGTH + 1
	copy(buf[hashOffset:], u.Hash)

	return buf
}

func (u *User) MarshalString() (string, error) {
	u.M_ID = u.ID()
	raw, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (u *User) Type() byte { return RT_user }

func (u *User) ID() ID { return NewID(byte(RT_user), []byte(u.Username)) }

func UnmarshalUserBinary(raw []byte) (*User, error) {
	var u User

	// Get name
	nameLength := uint(raw[0])
	u.Username = string(raw[1 : nameLength+1])

	// Get salt
	saltOffset := nameLength + 1
	saltLength := uint(raw[saltOffset])
	u.Salt = raw[saltOffset+1 : saltOffset+saltLength+1]

	// Get hash
	u.Hash = raw[saltOffset+saltLength+1:]

	return &u, nil
}

func UnmarshalUserText(raw string) (*User, error) {
	var u User
	err := json.Unmarshal([]byte(raw), &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (u *User) ValidatePassword(pass string) bool {
	key := make([]byte, len(pass)+len(u.Salt))
	copy(key[0:], pass)
	copy(key[len(pass):], u.Salt)
	res := bcrypt.CompareHashAndPassword(u.Hash, key)
	return res == nil
}

// Helper functions

func GenerateHash(pass string) ([]byte, []byte, error) {
	salt := make([]byte, PASS_SALT_LENGTH)
	num, err := rand.Read(salt) // Generate salt
	if num != PASS_SALT_LENGTH || err != nil {
		// Something went wrong...
		return nil, nil, err
	}
	key := make([]byte, PASS_SALT_LENGTH+len(pass))
	copy(key[0:len(pass)], pass)
	copy(key[len(pass):], salt)

	hashed, err := bcrypt.GenerateFromPassword(key, bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}
	return hashed, salt, nil
}
