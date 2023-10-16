package datasource

import (
	"Projeect/internal/model"
	"Projeect/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type PSQL struct {
	db *sql.DB
}

var (
	host     = os.Getenv("PS_HOST")
	port     = os.Getenv("PS_PORT")
	user     = os.Getenv("PS_USER")
	password = os.Getenv("PS_PASSWORD")
	dbname   = os.Getenv("PS_NAME")
)

func InitializePSQL() (psql *PSQL, err error) {
	connConf := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connConf)
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

func (psql *PSQL) SaveUser(user model.User) (id string, err error) {
	_, err = psql.db.Exec("INSERT INTO users (national_code, name, email, ip, status) VALUES ($1, $2, $3, $4, $5)",
		user.NationalCode, user.Name, user.Email, user.IPAddress, "pending")
	return user.NationalCode, err
}

func (psql *PSQL) FetchUser(nationalID string) (*model.User, error) {
	var user model.User
	err := psql.db.QueryRow("SELECT  national_code, name, email, ip, status FROM users WHERE national_code = $1",
		utils.EncodeBase64(nationalID)).Scan(&user.NationalCode, &user.Name, &user.Email, &user.IPAddress, &user.Status)
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
