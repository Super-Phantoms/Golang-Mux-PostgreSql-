# banking-auth

##### Run `./start.sh` to download the dependencies and run the the application

To run the application, you have to define the environment variables, default values of the variables are defined inside `start.sh`

- SERVER_ADDRESS    `[IP Address of the machine]`
- SERVER_PORT       `[Port of the machine]`
- DB_USER           `[Database username]`
- DB_PASSWD         `[Database password]`
- DB_ADDR           `[IP address of the database]`
- DB_PORT           `[Port of the database]`
- DB_NAME           `[Name of the database]`

# create .gitignore file
```sh
git init
# create a github repository  github.com/new

git remote add origin https://github.com/Honest67924/bank-lib.git
# add all files inside this project
git add .
git commit -m "extracted the error and logger from banking api"
#adding tag
git tag v1.0.0
git push origin master --tags

# search key words: using go moudles
# change git repository
git status
git add .gitignore ReadMe.md
git commit -m "added zap dependency"
git tag v1.0.1
git push origin master --tags
git push -f origin master
```


{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl9pZCI6IiIsImFjY291bnRzIjpudWxsLCJ1c2VybmFtZSI6InRpZ2VyIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNjg0OTM0NzgyfQ.Mtzi9STsP3Mb4p0bm-k4lz_CbfkJbLVfOoJGsRDjqOI",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaF90b2tlbiIsImNpZCI6IiIsImFjY291bnRzIjpudWxsLCJ1biI6InRpZ2VyIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNjg3NTI2NzUyfQ.MrbLbvws7NI7Octyn22Jm_6OWg0721BvK2BlIWzOmOY"
}

# 
```sh
import (
	"gorm.io/driver/<your_database_driver>"
	"gorm.io/gorm"
)

func main() {
	// Initialize GORM DB connection
	dsn := "<your_database_connection_string>"
	db, err := gorm.Open(<your_database_driver>.Open(dsn), &gorm.Config{})
	if err != nil {
		// Handle error
	}

	// Define the result struct
	type Result struct {
		Username       string `gorm:"column:username"`
		CustomerID     int    `gorm:"column:customer_id"`
		Role           string `gorm:"column:role"`
		AccountNumbers string `gorm:"column:account_numbers"`
	}

	// Execute the query
	var result Result
	err = db.Raw(`
		SELECT username, u.customer_id, u.role, ARRAY_TO_STRING(ARRAY_AGG(account_id), ',') as account_numbers
		FROM users u
		LEFT JOIN accounts a ON a.customer_id = u.customer_id
		WHERE username = 'Steven' AND password = '$2a$12$zZxjUraYPXc8z7vRUznbf.RQWUspBXkCFwPvm8A9p76VGTqck.JJm'
		GROUP BY username, u.customer_id, u.role
	`).Scan(&result).Error
	if err != nil {
		// Handle error
	}

	// Access the result
	fmt.Println("Username:", result.Username)
	fmt.Println("Customer ID:", result.CustomerID)
	fmt.Println("Role:", result.Role)
	fmt.Println("Account Numbers:", result.AccountNumbers)
}


import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/<your_database_driver>"
)

type User struct {
	Username   string `db:"username"`
	CustomerID int    `db:"customer_id"`
	Role       string `db:"role"`
}

func main() {
	// Initialize SQL database connection
	db, err := sqlx.Connect("<your_database_driver>", "<your_database_connection_string>")
	if err != nil {
		// Handle error
	}

	// Define the result struct
	type Result struct {
		Username       string `db:"username"`
		CustomerID     int    `db:"customer_id"`
		Role           string `db:"role"`
		AccountNumbers string `db:"account_numbers"`
	}

	// Execute the query
	var result Result
	err = db.Get(&result, `
		SELECT username, u.customer_id, u.role, ARRAY_TO_STRING(ARRAY_AGG(account_id), ',') as account_numbers
		FROM users u
		LEFT JOIN accounts a ON a.customer_id = u.customer_id
		WHERE username = 'Steven' AND password = '$2a$12$zZxjUraYPXc8z7vRUznbf.RQWUspBXkCFwPvm8A9p76VGTqck.JJm'
		GROUP BY username, u.customer_id, u.role
	`)
	if err != nil {
		// Handle error
	}

	// Access the result
	fmt.Println("Username:", result.Username)
	fmt.Println("Customer ID:", result.CustomerID)
	fmt.Println("Role:", result.Role)
	fmt.Println("Account Numbers:", result.AccountNumbers)
}

```