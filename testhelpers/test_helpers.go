package testhelpers

import (
	"github.com/astaxie/beego/orm"
	"github.com/ilackarms/pkg/errors"
	_ "github.com/mattn/go-sqlite3" // import your required driver
	"os"
)

func InitTestDB() error {
	if err := orm.RegisterDataBase("default", "sqlite3", "test_data.db"); err != nil {
		return errors.New("registering to sqlite db", err)
	}
	if err := orm.RunSyncdb("default", true, true); err != nil {
		return errors.New("syncing db", err)
	}
	return nil
}

func DestroyTestDB() error {
	if err := os.Remove("test_data.db"); err != nil {
		return errors.New("cleaning up test_data.db", err)
	}
	return nil
}

func RefreshDBState() error {
	if err := orm.RunSyncdb("default", true, true); err != nil {
		return errors.New("syncing db", err)
	}
	return nil
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
