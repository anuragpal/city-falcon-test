package models

import (
    "time"
    "errors"
    "strings"
    "strconv"
    "database/sql"
    "golang.org/x/crypto/bcrypt"
    "github.com/gofiber/fiber/v2"
    "github.com/go-ozzo/ozzo-validation/v4"
    "github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
    Id                  string              `json:"id,omitempty"`
    FirstName           string              `json:"first_name,omitempty"`
    LastName            string              `json:"last_name,omitempty"`
    MiddleName          string              `json:"middle_name,omitempty"`
    Password            string              `json:"password,omitempty"`
    ConfirmPassword     string              `json:"confirm_password,omitempty"`
    Email               string              `json:"email,omitempty"`
    Status              int                 `json:"status"`
    AddedDate           *time.Time          `json:"added_date,omitempty"`
    ModifiedDate        *time.Time          `json:"modified_date,omitempty"`
}

func (u User) UniqueEmail(value interface{}) error {
    var (
        cnt         int
        sql         string
        params      []interface{}
    )

    un, _ := value.(string)
    if u.Id != "" {
        sql = `SELECT count(id) as cnt FROM users WHERE email = $1 AND id <> $2`
        params = append(params, un, u.Id)
    } else {
        sql = `SELECT count(id) as cnt FROM users WHERE email = $1`
        params = append(params, un)
    }
    stmt, err := DB.Prepare(sql)
    if err != nil {
    }

    err = stmt.QueryRow(params...).Scan(&cnt)
    if err != nil {
    }

    if cnt > 0 {
        return errors.New("Email entered already taken.")
    }
    return nil
}

func (u User) ValidatePassword(value interface{}) error {
    if u.Password != u.ConfirmPassword {
        return errors.New("Password not matched.")
    }
    return nil
}

func (u User) Validate() error {
    return validation.ValidateStruct(&u,
        validation.Field(&u.FirstName, validation.Required.Error("Please provide first name."), validation.Length(3, 20)),
        validation.Field(&u.LastName, validation.Required.Error("Please provide last name."), validation.Length(3, 20)),
        validation.Field(&u.Email, validation.When(u.Email != "", is.Email.Error("Please enter valid email."), validation.By(u.UniqueEmail)).Else(validation.Nil)),
        validation.Field(&u.Password, validation.Required.Error("Please provide password.")),
        validation.Field(&u.ConfirmPassword, validation.Required.Error("Please reenter password."), validation.By(u.ValidatePassword)),
    )
}

func (u User) Add() (fiber.Map, int) {
    err := u.Validate()
    if err != nil {
        return fiber.Map{"errors": err}, fiber.StatusBadRequest
    }

    pwd, errPwd := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if errPwd != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }

    err = DB.QueryRow(`
    INSERT INTO users (
        first_name, last_name, middle_name, email, password, added_date, modified_date
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7
    ) returning id`, u.FirstName, u.LastName, u.MiddleName, u.Email, string(pwd), time.Now().UTC(), time.Now().UTC()).Scan(&u.Id)

    if err != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }
    return fiber.Map{}, fiber.StatusCreated
}

func (u User) ValidateUpdateUser() error {
    return validation.ValidateStruct(&u,
        validation.Field(&u.FirstName, validation.Required.Error("Please provide first name."), validation.Length(3, 20)),
        validation.Field(&u.LastName, validation.Required.Error("Please provide last name."), validation.Length(3, 20)),
        validation.Field(&u.Email, validation.When(u.Email != "", is.Email.Error("Please enter valid email."), validation.By(u.UniqueEmail)).Else(validation.Nil)),
    )
}

func (u User) Update() (fiber.Map, int) {
    sqlQuery := `
    SELECT
        u.id
    FROM users as u
    WHERE
        u.id = $1
    `
    stmt, err := DB.Prepare(sqlQuery)
    if err != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }

    err = stmt.QueryRow(u.Id).Scan(&u.Id)
    if err != nil && err != sql.ErrNoRows {
        return fiber.Map{}, fiber.StatusInternalServerError
    } else if err == sql.ErrNoRows {
        return fiber.Map{}, fiber.StatusNotFound
    }

    err = u.ValidateUpdateUser()
    if err != nil {
        return fiber.Map{"errors": err}, fiber.StatusBadRequest
    }

    updateQuery := `UPDATE users SET first_name = $1, last_name = $2, middle_name = $3, email = $4, modified_date = $5 WHERE id = $6`
    _, errUpdate := DB.Exec(updateQuery, u.FirstName, u.LastName, u.MiddleName, u.Email, time.Now().UTC(), u.Id)
    if errUpdate != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }

    return fiber.Map{}, fiber.StatusOK
}

func (u User) Remove() (fiber.Map, int) {
    sqlQuery := `
    SELECT
        u.id
    FROM users as u
    WHERE
        u.id = $1 AND
        u.status = 1
    `
    stmt, err := DB.Prepare(sqlQuery)
    if err != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }

    err = stmt.QueryRow(u.Id).Scan(&u.Id)
    if err != nil && err != sql.ErrNoRows {
        return fiber.Map{}, fiber.StatusInternalServerError
    } else if err == sql.ErrNoRows {
        return fiber.Map{}, fiber.StatusNotFound
    }

    updateQuery := `UPDATE users SET status = $1, modified_date = $2 WHERE id = $3`
    _, errUpdate := DB.Exec(updateQuery, 2, time.Now().UTC(), u.Id)
	if errUpdate != nil {
		return fiber.Map{}, fiber.StatusInternalServerError
	}
    return fiber.Map{}, fiber.StatusOK
}

func (u User) List(l ListParams) (fiber.Map, int) {
    var (
        condition_count             int
        params                      []interface{}
        conditions                  []string
        total                        int
    )

    conditions = []string{}
    condition_count++
    params = append(params, 1)
    conditions = append(conditions, " u.status = $"+strconv.Itoa(condition_count))

    if l.Query != "" {
        condition_count = condition_count + 2
        params = append(params, "%"+l.Query+"%", "%"+l.Query+"%")
        conditions = append(conditions, " (u.first_name LIKE $2 OR u.last_name LIKE $3)")
    }

    users := []User{}
    if l.RecordPerPage <= 0 {
        l.RecordPerPage = 10
    }

    if l.Page <= 0  {
        l.Page = 1
    }

    offset := (l.Page - 1) * l.RecordPerPage

    stmt, err := DB.Prepare(`SELECT count(id) as cnt FROM users as u WHERE ` + strings.Join(conditions, " AND "))
    if err != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }

    err = stmt.QueryRow(params...).Scan(&total)
    if err != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }

    limit_param_count := condition_count + 1
    params = append(params, l.RecordPerPage)

    offset_param_count := limit_param_count + 1
    params = append(params, offset)

    sql := `
    SELECT
        u.id,
        u.first_name,
        u.last_name,
        COALESCE(u.middle_name, '') as middle_name,
        COALESCE(u.email, '') as email,
        u.status
    FROM users as u
    WHERE ` + strings.Join(conditions, " AND ") + " LIMIT $"+strconv.Itoa(limit_param_count) + " OFFSET $" + strconv.Itoa(offset_param_count)
    rows, errRows := DB.Query(sql, params...)

    if errRows != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }

    for rows.Next() {
        user := User{}
        errScan := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.MiddleName, &user.Email, &user.Status)
        if errScan != nil {
            return fiber.Map{}, fiber.StatusInternalServerError
        }
        users = append(users, user)
    }
    
    return fiber.Map{
        "users": users,
        "rpp": l.RecordPerPage,
        "sort_by": l.SortBy,
        "order_by": l.OrderBy,
        "page": l.Page,
        "total": total,
    }, fiber.StatusOK
}

func (u User) Details() (fiber.Map, int) {
    query := `
    SELECT
        u.id,
        u.first_name,
        u.last_name,
        COALESCE(u.middle_name, '') as middle_name,
        COALESCE(u.email, '') as email,
        u.status
    FROM users as u
    WHERE
        u.id = $1 AND
        u.status = 1
    `
    
    stmt, err := DB.Prepare(query)

    if err != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }

    err = stmt.QueryRow(u.Id).Scan(&u.Id, &u.FirstName, &u.LastName, &u.MiddleName, &u.Email, &u.Status)
    if err != nil  && err != sql.ErrNoRows {
        return fiber.Map{}, fiber.StatusInternalServerError
    } else if err == sql.ErrNoRows {
        return fiber.Map{}, fiber.StatusNotFound
    }
    return fiber.Map{"user": u}, fiber.StatusOK
}