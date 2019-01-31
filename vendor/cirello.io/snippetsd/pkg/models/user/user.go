// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user // import "cirello.io/snippetsd/pkg/models/user"
import (
	"context"
	"fmt"

	"cirello.io/errors"
	"golang.org/x/crypto/bcrypt"
)

// User aggregates all the information of a snippet user.
type User struct {
	ID       int64  `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`
	Team     string `db:"team" json:"team"`
}

func (u *User) String() string {
	return fmt.Sprintf("%v", u.Email)
}

// EncryptPassword takes the current value of the password, and encrypts it
// using bcrypt. If it is already encrypted, it will encrypt twice.
func (u *User) EncryptPassword() error {
	// TODO: if possible, prevent double encryption

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return errors.E(err, "cannot encrypt password")
	}
	u.Password = string(hash)
	return nil
}

// NewFromEmail creates a user from a given email.
func NewFromEmail(email, password, team string) (*User, error) {
	// TODO: validate email
	return &User{
		Email:    email,
		Password: password,
		Team:     team,
	}, nil
}

func (u *User) checkPassword(plainPassword string) bool {
	byteHash := []byte(u.Password)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPassword))
	u.Password = ""
	return err == nil
}

// Authenticate finds a user and check their password.
func Authenticate(u *User, email, password string) error {
	if !u.checkPassword(password) {
		return errors.E("invalid credentials")
	}
	return nil
}

type userCtxKeyType struct{}

var userCtxKey = userCtxKeyType{}

// WithContext creates a new context attaching the user.
func WithContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userCtxKey, u)
}

// WhoAmI loads the user attached to the context.
func WhoAmI(ctx context.Context) *User {
	if u, ok := ctx.Value(userCtxKey).(*User); ok {
		return u
	}
	return nil
}
