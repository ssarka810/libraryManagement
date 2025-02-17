package db

import (
	"fmt"
	"net/url"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

var db *gorm.DB

func 	Init() *gorm.DB {
	dns := url.URL{
		User:     url.UserPassword("postgres", "postgres"),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", "localhost", 5432),
		Path:     "postgres",
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	logrus.Info("DB Connection  URL: ", dns.String())
	connection, err := gorm.Open("postgres", dns.String())
	if err != nil {
		logrus.Error("Not able to connect the DB , ", err)
	} else {
		db = connection
		db.SingularTable(true)
		db.Debug().AutoMigrate(
			&UserDetails{},
			&BookDetails{},
			&BookCirculationDetails{},
		)
	}

	return connection
}

func GetDB() *gorm.DB {
	return db
}
