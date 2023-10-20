package datasource

import (
	"Projeect/internal/model"
	"Projeect/utils"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

type PSQL struct {
	db *sql.DB
}

var (
	uri = os.Getenv("PS_URI")
)

func InitializePSQL() (psql *PSQL, err error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		logrus.Errorf("connection to database failed: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logrus.Errorf("database ping has faild: %v", err)
		return nil, err
	}

	db.Exec(`create table users (
    national_code varchar(20) PRIMARY KEY ,
    name varchar(40),
    email varchar(50),
    ip varchar(20),
    status varchar(50)
)`)

	return &PSQL{db: db}, nil
}

func (psql *PSQL) SaveUser(user model.User) (id string, err error) {
	_, err = psql.db.Exec("INSERT INTO users (national_code, name, email, ip, status) VALUES ($1, $2, $3, $4, $5)",
		user.NationalCode, user.Name, user.Email, user.IPAddress, "pending")
	return user.NationalCode, err
}

func (psql *PSQL) FetchUser(nationalCode string) (*model.User, error) {
	var user model.User
	err := psql.db.QueryRow("SELECT  national_code, name, email, ip, status FROM users WHERE national_code = $1",
		utils.EncodeBase64(nationalCode)).Scan(&user.NationalCode, &user.Name, &user.Email, &user.IPAddress, &user.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	user.NationalCode, err = utils.DecodeBase64(user.NationalCode)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
