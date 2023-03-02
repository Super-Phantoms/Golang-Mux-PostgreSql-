package domain

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/golangdevm/fullstack/dto"
	"github.com/golangdevm/fullstack/errs"
	"github.com/golangdevm/fullstack/logger"
	"github.com/jmoiron/sqlx"

	"github.com/jinzhu/gorm"
)

type AuthRepository interface {
	FindBy(username string, password string) (*Login, *errs.AppError)
	Save(dto.RegisterRequest) (*User, *errs.AppError)
	GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError)
	RefreshTokenExists(refreshToken string) (*RefreshToken, *errs.AppError)

	IsAuthorized(token string, routeName string, vars map[string]string) bool
}

type AuthRepositoryDb struct {
	client    *gorm.DB
	sqlClient *sqlx.DB
}

func (d AuthRepositoryDb) Save(req dto.RegisterRequest) (*User, *errs.AppError) {

	saveErr := req.BeforeSave()
	if saveErr != nil {
		logger.Error("Error while hashing password: " + saveErr.Error())
		return &User{}, errs.NewValidationError("hasing error")
	}

	req.Prepare()

	valErr := req.Validate("register")
	if valErr != nil {
		logger.Error("validation error: " + valErr.Error())
		return &User{}, errs.NewValidationError("validation error")
	}

	newUser := User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
	}

	err := d.client.Debug().Model(&User{}).Create(&newUser).Error
	if err != nil {

		logger.Error("Error while creating new user: " + err.Error())
		return &User{}, errs.NewUnexpectedError("Unexpected error from database")
	} else {
		return &newUser, nil
	}
}
func (d AuthRepositoryDb) RefreshTokenExists(refreshToken string) (*RefreshToken, *errs.AppError) {

	// fmt.Println("#######AccessToken:::::", request.AccessToken)
	r := &RefreshToken{}
	err := d.client.Debug().Model(RefreshToken{}).Where("refresh_token = ?", refreshToken).Take(r).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errs.NewAuthenticationError("refresh token not registered in the store")
		} else {
			logger.Error("Unexpected database error: " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}

	return r, nil
}

func (d AuthRepositoryDb) GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError) {
	// generate the refresh token
	var appErr *errs.AppError
	var refreshToken string
	if refreshToken, appErr = authToken.newRefreshToken(); appErr != nil {
		return "", appErr
	}

	// store it in the store
	var r_token = RefreshToken{
		Refresh_Token: refreshToken,
	}
	err := d.client.Debug().Create(&r_token).Error
	if err != nil {
		logger.Error("unexpected database error: " + err.Error())
		return "", errs.NewUnexpectedError("unexpected database error")
	}
	return refreshToken, nil
}

func (d AuthRepositoryDb) FindBy(username, password string) (*Login, *errs.AppError) {

	var login Login
	sqlVerify := `
		SELECT username, u.customer_id, u.role, ARRAY_TO_STRING(ARRAY_AGG(account_id), ',') as accounts
		FROM users u
		LEFT JOIN accounts a ON a.customer_id = u.customer_id
		WHERE username = ?
		GROUP BY username, u.customer_id, u.role
	`
	err := d.client.Raw(sqlVerify, username).Scan(&login).Error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("invalid credentials")
		} else {
			logger.Error("Error while verifying login request from database: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &login, nil
}

// ////////////banking main
func (d AuthRepositoryDb) IsAuthorized(token string, routeName string, vars map[string]string) bool {

	u := buildVerifyURL(token, routeName, vars)

	if response, err := http.Get(u); err != nil {
		fmt.Println("Error while sending..." + err.Error())
		return false
	} else {
		m := map[string]bool{}
		if err = json.NewDecoder(response.Body).Decode(&m); err != nil {
			logger.Error("Error while decoding response from auth server:" + err.Error())
			return false
		}
		return m["isAuthorized"]
	}
}

func buildVerifyURL(token string, routeName string, vars map[string]string) string {
	u := url.URL{Host: "localhost:8080", Path: "/auth/verify", Scheme: "http"}
	q := u.Query()
	q.Add("token", token)
	q.Add("routeName", routeName)
	for k, v := range vars {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func NewAuthRepository(client *gorm.DB, sqlClient *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client, sqlClient}
}
