package users

import (
	"time"

	"github.com/Financial-Times/gourmet/apperror"
	"github.com/Financial-Times/gourmet/example/usersapp/pkg/storage"
	"github.com/doug-martin/goqu/v7"
	"github.com/doug-martin/goqu/v7/exec"
)

const (
	defaultLimit  uint = 100
	defaultOffset uint = 0
)

type Repository interface {
	AddUser(c *User) (*User, error)
	RemoveUser(cID int) error
	FindAllUsers(opts *storage.QueryOptions) ([]User, error)
	FindUserByID(cID int) (User, error)
}

type UserRepository struct {
	db storage.Persistence
}

func NewUserRepository(db storage.Persistence) Repository {
	return &UserRepository{db: db}
}

func (r *UserRepository) AddUser(c *User) (*User, error) {
	created := time.Now().Unix()

	result, err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		c.Created = created
		c.LastUpdated = created
		return tx.From("User").Insert(c)
	});
	if err != nil {
		return nil, apperror.DBError.Wrap(err, "error adding new User")
	}

	cID, _ := result.LastInsertId()
	c.UserID = int(cID)

	return c, nil
}

func (r *UserRepository) RemoveUser(cID int) error {
	_, err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		return tx.From("User").Where(goqu.Ex{"cid": cID}).Delete()
	})

	if err != nil {
		return apperror.DBError.Wrapf(err, "error deleting Users with id %d", cID)
	}
	return nil
}

func (r *UserRepository) FindAllUsers(opts *storage.QueryOptions) (cc []User, err error) {
	if opts.Limit == 0 {
		opts.Limit = defaultLimit
	}

	err = r.db.DB.From("User").
		Limit(opts.Limit).
		Offset(opts.Offset).
		ScanStructs(&cc)

	if err != nil {
		return nil, apperror.DBError.Wrapf(err, "error getting all Users")
	}
	return cc, nil
}

func (r *UserRepository) FindUserByID(uID int) (u User, err error) {
	found, err := r.db.DB.From("User").Where(
		goqu.C("uid").Eq(uID),
	).ScanStruct(&u)

	if !found {
		return u, apperror.NotFound.Newf("user with ID %d not found", uID)
	}

	if err != nil {
		return u, apperror.DBError.Wrapf(err, "error getting user with ID %d", uID)
	}

	return u, nil
}
