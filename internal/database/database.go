// everything about creating databases
// reading from databases
// updating databases
// and deleting info from them

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	processjwt "github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions/internal/service/security/processJWT"
)

type Storage struct {
	DB *sql.DB
}

type (
	// user will input his or her email, name and password and we'll generate and id and salt
	User struct {
		ID       int64  `json:"id"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// id of the expressison and user
	// time of executing power, multiplication, division, addition and substraction
	// result
	Expression struct {
		ID                 int64  `json:"id"`
		Expression         string `json:"expression"`
		UserID             int64  `json:"user_id"`
		TimePower          int64  `json:"time_power"`
		TimeMultiplication int64  `json:"time_mult"`
		TimeDivision       int64  `json:"time_div"`
		TimeAddition       int64  `json:"time_add"`
		TimeSubstraction   int64  `json:"time_sub"`
		Result             string `json:"result"`
	}
)

type (
	ExpressionReq struct {
		ID                 int64  `json:"id"`
		Expression         string `json:"expression"`
		UserID             int64  `json:"user_id"`
		TimePower          int64  `json:"time_power"`
		TimeMultiplication int64  `json:"time_mult"`
		TimeDivision       int64  `json:"time_div"`
		TimeAddition       int64  `json:"time_add"`
		TimeSubstraction   int64  `json:"time_sub"`
	}
)

// Default values for the case of result absence.
// And some additional values for executing proper
// time for each operation in every expression
const (
	null      = "null"
	std_power = 5
	std_mult  = 5
	std_div   = 5
	std_add   = 5
	std_sub   = 5
)

// done
func InitDB() error {
	var err error
	var storage Storage
	// production code shouldn't include context.TODO
	// but for temporary use — connection to database — we create it
	ctx := context.TODO()

	storage.DB, err = sql.Open("sqlite3", "./internal/database/store.db")
	if err != nil {
		return err
	}

	err = storage.DB.PingContext(ctx)
	if err != nil {
		return err
	}

	err = createUsersTable(ctx, storage.DB)
	if err != nil {
		log.Fatal("Error creating 'Users' table:", err)
	}

	err = createExpressionsTable(ctx, storage.DB)
	if err != nil {
		log.Fatal("Error creating 'Expressions' table:", err)
	}

	return nil
}

// done
func (db *Storage) Closee() {
	db.DB.Close()
}

func createUsersTable(ctx context.Context, db *sql.DB) error {
	const (
		usersTable = `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`
	)

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	return nil
}

func createExpressionsTable(ctx context.Context, db *sql.DB) error {
	const (
		expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		id 		   INTEGER PRIMARY KEY AUTOINCREMENT, 
		expression TEXT NOT NULL,
		userid 	   INTEGER,
		timepower  INTEGER,
		timemult   INTEGER,
		timediv    INTEGER,
		timeadd    INTEGER,
		timesub    INTEGER,
		result     TEXT,
		FOREIGN KEY (userid) REFERENCES users (id)
	);`
	)

	if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
		log.Fatal("Error creating table:", err)
		return err
	}

	return nil
}

// done
// puts user into the table
func (db *Storage) RegisterUser(ctx context.Context, user *User) error {
	var query = `
	INSERT INTO users (id, email, name, password) values ($1, $2, $3, $4)
	`
	saltedBytes := []byte(user.Password)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	_, err = db.DB.ExecContext(ctx, query, user.ID, user.Email, user.Name, hashedBytes)
	if err != nil {
		return err
	}

	return nil
}

// done
// selects a user and generates JWT Token
func (db *Storage) LoginUser(ctx context.Context, user1 *User) (string, error) {
	// create a query which we pass to our database
	query := `SELECT id, password FROM users WHERE email=$1 AND name=$2`

	// create a variable for storing data from database with query
	var user User
	// take values
	err := db.DB.QueryRowContext(ctx, query, user1.Email, user1.Name).Scan(&user.ID, &user.Password)
	if err != nil {
		return "", err
	}
	// comparing passwords of actual and passed users
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user1.Password))
	if err != nil {
		return "", err
	}

	// generating 'salt' for our password — hash
	token, err := processjwt.GenerateNewJWTToken(user.ID)
	if err != nil {
		return "", nil
	}

	return token, nil
}

// done
// POST method
func (db *Storage) InsertExpression(ctx context.Context, expression *Expression) (int64, error) {
	// create a query for inserting expression which user have passed to our program using REST API (delivery)
	var query = `INSERT INTO expressions (expression, user_id, time_power, time_mult, time_div, time_add, time_sub, result) values ($1, $2, $3, $4, $5, $6, $7, $8)`
	result, err := db.DB.ExecContext(ctx, query, expression.ID, expression.Expression, expression.UserID,
		expression.TimePower, expression.TimeMultiplication, expression.TimeDivision,
		expression.TimeAddition, expression.TimeSubstraction, expression.Result)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// done
// GET method
func (db *Storage) FetchAllExpressions(ctx context.Context) ([]Expression, error) {
	var expReturn []Expression
	query := `SELECT id, expression, time_power, time_mult, time_div, time_add, time_sub, result FROM expressions`
	// expects to return several values after parsing a query
	rows, err := db.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	// define closing request to database
	defer rows.Close()

	for rows.Next() {
		var Result sql.NullInt64
		exp := Expression{}
		err := rows.Scan(&exp.ID, &exp.Expression, &exp.TimePower,
			&exp.TimeMultiplication, &exp.TimeDivision, &exp.TimeAddition, &exp.TimeSubstraction, &Result)

		if err != nil {
			return nil, err
		}

		if Result.Valid {
			exp.Result = strconv.Itoa(int(Result.Int64))
		} else if !Result.Valid {
			exp.Result = "0"
		}

		expReturn = append(expReturn, exp)
	}

	return expReturn, nil
}

// done
func (db *Storage) FetchExpressionByUserID(ctx context.Context, userID int64) ([]Expression, error) {
	var Result sql.NullInt64
	var expReturn []Expression
	query := `SELECT id, expression, time_power, time_mult, time_div, time_add, time_sub, result FROM expressions WHERE user_id=$1`
	rows, err := db.DB.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		exp := Expression{}
		err := rows.Scan(&exp.ID, &exp.Expression, &exp.TimePower,
			&exp.TimeMultiplication, &exp.TimeDivision, &exp.TimeAddition, &exp.TimeSubstraction, &Result)

		if err != nil {
			return nil, err
		}

		if Result.Valid {
			exp.Result = strconv.Itoa(int(Result.Int64))
		} else if !Result.Valid {
			exp.Result = "0"
		}

		expReturn = append(expReturn, exp)
	}

	return expReturn, nil
}

// done
func (db *Storage) FetchExpressionByID(ctx context.Context, ID int64) (Expression, error) {
	var Result sql.NullInt64
	expReturn := Expression{}
	query := `SELECT user_id, expression, time_power, time_mult, time_div, time_add, time_sub, result FROM expressions WHERE id=$1`
	err := db.DB.QueryRowContext(ctx, query, ID).Scan(&expReturn.UserID, &expReturn.Expression,
		&expReturn.TimePower, &expReturn.TimeMultiplication,
		&expReturn.TimeDivision, &expReturn.TimeAddition,
		&expReturn.TimeSubstraction, &Result)

	if err != nil {
		return expReturn, err
	}

	if Result.Valid {
		expReturn.Result = strconv.Itoa(int(Result.Int64))
	} else if !Result.Valid {
		expReturn.Result = "0"
	}

	return expReturn, nil
}

// done
// PUT method
func (db *Storage) UpdateExpressionByID(ctx context.Context, result string, ID int64) error {
	query := `UPDATE expressions SET result=$1 WHERE id=$2`

	_, err := db.DB.ExecContext(ctx, query, result, ID)

	if err != nil {
		return err
	}

	return nil
}

// DELETE method
func (db *Storage) DeleteExpressionByID(ctx context.Context, ID int64) error {
	query := `DELETE FROM expressions WHERE id = ?`

	_, err := db.DB.ExecContext(ctx, query, ID)

	if err != nil {
		return err
	}

	return nil
}

// REST API level
// decided to put it here, not delivery\http\v1\handlers
// because it's more convinient to use
// for further scaling we'll need to put it in different files, even directories
// but now it's here
// for help — read comments
// wrapping connections to database
// that connection — for adding users

// RegisterUser
func RegisterHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			log.Printf("method not allowed: %d", http.StatusMethodNotAllowed)
		}

		w.Header().Set("Content-Type", "application/json")
		var dbCopy *Storage

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("error: %v", err)
			return
		}

		var userStruct = User{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Password: user.Password,
		}

		err := dbCopy.RegisterUser(ctx, &userStruct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "This username is already registered", http.StatusUnauthorized)
			log.Printf("error: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Print("200 OK. Token Created. RegisterHandler")
	}
}

// LoginUser
func LoginHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

		w.Header().Set("Content-Type", "application/json")
		var dbCopy *Storage

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("error: %v", err)
			return
		}

		var userStruct = User{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Password: user.Password,
		}

		token, err := dbCopy.LoginUser(ctx, &userStruct)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			log.Printf("error: %v", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth token",
			Value:    token,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})

		//http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
		w.WriteHeader(http.StatusAccepted)
		log.Print("200 OK. Token Created. LoginHandler")
	}
}

// InsertExpression
func CreateExpressionHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

		w.Header().Set("Content-Type", "application/json")
		var dbCopy *Storage
		var expression ExpressionReq

		err := json.NewDecoder(r.Body).Decode(&expression)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		var expStruct = Expression{
			Expression:         expression.Expression,
			UserID:             expression.UserID,
			TimePower:          expression.TimePower,
			TimeMultiplication: expression.TimeMultiplication,
			TimeDivision:       expression.TimeDivision,
			TimeAddition:       expression.TimeAddition,
			TimeSubstraction:   expression.TimeSubstraction,
			Result:             null,
		}

		id, err := dbCopy.InsertExpression(ctx, &expStruct)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("accepted: %v", id)
		w.WriteHeader(http.StatusAccepted)
	}
}

// FetchExpressionByID
func GetExpressionByIDHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

		w.Header().Set("Content-Type", "application/json")
		var dbCopy *Storage

		var ID int64
		intID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "unable to receive ID from request", http.StatusBadRequest)
		}
		ID = int64(intID)

		expression, err := dbCopy.FetchExpressionByID(ctx, ID)

		if err != nil {
			w.WriteHeader(http.StatusInsufficientStorage)
			http.Error(w, "unable to fetch all expressions from database", http.StatusInsufficientStorage)
			return
		}

		json.NewEncoder(w).Encode(expression)
	}
}

// FetchExpressionByUserID
func GetExpressionByUserIDHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

		w.Header().Set("Content-Type", "application/json")
		var dbCopy *Storage

		var userID int64
		intUserID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "unable to receive ID from request", http.StatusBadRequest)
		}
		userID = int64(intUserID)

		allExpressions, err := dbCopy.FetchExpressionByUserID(ctx, userID)

		if err != nil {
			w.WriteHeader(http.StatusInsufficientStorage)
			http.Error(w, "unable to fetch all expressions from database", http.StatusInsufficientStorage)
			return
		}

		var answer []Expression
		for _, value := range allExpressions {
			response := Expression{
				ID:                 value.ID,
				Expression:         value.Expression,
				UserID:             userID,
				TimePower:          std_power,
				TimeMultiplication: std_mult,
				TimeDivision:       std_div,
				TimeAddition:       std_add,
				TimeSubstraction:   std_sub,
				Result:             value.Result,
			}

			answer = append(answer, response)
		}

		json.NewEncoder(w).Encode(answer)
	}
}

// FetchAllExpressions
func GetAllExpressionsHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

		w.Header().Set("Content-Type", "application/json")
		var dbCopy *Storage

		allExpressions, err := dbCopy.FetchAllExpressions(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInsufficientStorage)
			http.Error(w, "unable to fetch all expressions from database", http.StatusInsufficientStorage)
			return
		}

		var answer []Expression
		for _, value := range allExpressions {
			response := Expression{
				ID:                 value.ID,
				Expression:         value.Expression,
				UserID:             value.UserID,
				TimePower:          std_power,
				TimeMultiplication: std_mult,
				TimeDivision:       std_div,
				TimeAddition:       std_add,
				TimeSubstraction:   std_sub,
				Result:             value.Result,
			}

			answer = append(answer, response)
		}

		json.NewEncoder(w).Encode(answer)
	}
}

// DeleteExpressionByID
func DeleteExpressionHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

		w.Header().Set("Content-Type", "application/json")
		var dbCopy *Storage

		var ID int64
		intID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "unable to receive ID from request", http.StatusBadRequest)
		}
		ID = int64(intID)

		err = dbCopy.DeleteExpressionByID(ctx, int64(ID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusAccepted)
		log.Print("success deleting!!!")
	}
}
