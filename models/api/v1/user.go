package v1

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	wrapErr "github.com/ilackarms/pkg/errors"
	"golang.org/x/crypto/scrypt"
	"io"
	"reflect"
	"strings"
	"time"
)

// User represents a registered Partikle user
type User struct {
	ID           int64     `orm:"auto"`
	Username     string    `orm:"type(longtext)"`
	PasswordHash string    `orm:"type(longtext)"`
	PasswordSalt string    `orm:"type(longtext)"`
	Created      time.Time `orm:"type(datetime)"`
}

const pwHashBytes = 64

func generateSalt() (string, error) {
	buf := make([]byte, pwHashBytes)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", wrapErr.New("reading random bytes into buf", err)
	}

	return fmt.Sprintf("%x", buf), nil
}

func generatePassHash(password string, salt string) (string, error) {
	h, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, pwHashBytes)
	if err != nil {
		return "", wrapErr.New("generating cryptographic key using salt", err)
	}

	return fmt.Sprintf("%x", h), nil
}

// NewUser returns a *User object with hashed password and salt
func NewUser(username, password string) (*User, error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, wrapErr.New("generating salt", err)
	}
	passHash, err := generatePassHash(password, salt)
	if err != nil {
		return nil, wrapErr.New("generating password hash", err)
	}
	return &User{
		Username:     username,
		PasswordHash: passHash,
		PasswordSalt: salt,
		Created:      time.Now(),
	}, nil
}

func init() {
	orm.RegisterModel(new(User))
}

// AddUser insert a new User into database and returns
// last inserted Id on success.
func AddUser(m *User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserByID retrieves User by Id. Returns error if
// Id doesn't exist
func GetUserByID(id int64) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{ID: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUser retrieves all User matches certain condition. Returns empty list if
// no records exist
func GetAllUser(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []User
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateUserByID updates User by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserByID(m *User) (err error) {
	o := orm.NewOrm()
	v := User{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUser(id int64) (err error) {
	o := orm.NewOrm()
	v := User{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&User{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
