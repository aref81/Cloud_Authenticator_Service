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
	PS_URI = os.Getenv("PS_URI")
)

func InitializePSQL() (psql *PSQL, err error) {
	//connConf := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	uri := "postgres://avnadmin:AVNS_r5bk62SSBQn99QgS0cX@pg-sahabi-boboy1390-02ae.a.aivencloud.com:21538/defaultdb?sslmode=require"

	db, err := sql.Open("postgres", uri)
	if err != nil {
		logrus.Errorf("connection to database failed: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logrus.Errorf("database ping has failde: %v", err)
		return nil, err
	}

	return &PSQL{db: db}, nil
}

func (psql *PSQL) FetchUser(encodedNationalCode string) (*model.User, error) {
	var user model.User
	err := psql.db.QueryRow("SELECT  national_code, name, email, ip, status FROM users WHERE national_code = $1",
		encodedNationalCode).Scan(&user.NationalCode, &user.Name, &user.Email, &user.IPAddress, &user.Status)
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

func (psql *PSQL) UpdateStatus(nationalCode, newStatus string) error {
	_, err := psql.db.Exec("UPDATE users SET status = $1 WHERE national_code = $2", newStatus, utils.EncodeBase64(nationalCode))
	if err != nil {
		return err
	}

	return nil
}
