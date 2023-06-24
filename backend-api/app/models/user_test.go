package models

import (
    "log"
    "testing"
    "github.com/google/uuid"
    "github.com/gofiber/fiber/v2"
) 

func TestUniqueEmail(t *testing.T) {
    u := User{}
    id := uuid.New()
    u.Email = id.String() + "@gmail.com"
    u.FirstName = "anurag"
    u.LastName = "pal"
    u.Password = "12345678"
    u.ConfirmPassword = "12345678"
    err := u.UniqueEmail(u.Email)
    if err != nil {
        t.Errorf("Expected email should not exist %s", u.Email)
    }
    u.Add()

    err = u.UniqueEmail(u.Email)
    if err.Error() != "Email entered already taken." {
        t.Errorf("Expected email should exist %s", u.Email)
    }

    u1 := User{}
    id1 := uuid.New()
    u1.Email = id1.String() + "@gmail.com"
    u1.FirstName = "anurag"
    u1.LastName = "pal"
    u1.Password = "12345678"
    u1.ConfirmPassword = "12345678"
    i1, _ := u1.Add()
    u1.Id = i1["id"].(string)
    u1.Email = u.Email
    err = u1.UniqueEmail(u1.Email)
    if err.Error() != "Email entered already taken." {
        t.Errorf("Expected email should exist %s", u.Email)
    }

}

func TestValidatePassword(t *testing.T) {
    u := User{}
    u.Password = "12345678"
    u.ConfirmPassword = "123456789"

    err := u.ValidatePassword(u.Password)
    if err.Error() != "Password not matched."  {
        t.Errorf("Expected password and confirm password should not match.")
    }

    u.ConfirmPassword = "12345678"
    err = u.ValidatePassword(u.Password)
    if err != nil {
        t.Errorf("Expected password and confirm password should match.")
    }
}

func TestList(t *testing.T) {
    u := User{}
    params := ListParams{
        RecordPerPage: 10,
        Page: 1,
        Query: "test-not-exists",
    }

    result, status := u.List(params)
    if status != fiber.StatusOK {
        t.Errorf("Expected status code 200, got %d", status)
    }

    data := result["users"].([]User)
    if len(data) != 0 {
        t.Errorf("Expected 0 row in result, got %d", len(data))
    }

    params = ListParams{
        RecordPerPage: 0,
        Page: 0,
        Query: "test-not-exists",
    }

    result, status = u.List(params)
    if status != fiber.StatusOK {
        t.Errorf("Expected status code 200, got %d", status)
    }

    data = result["users"].([]User)
    if len(data) != 0 {
        t.Errorf("Expected 0 row in result, got %d", len(data))
    }

    params = ListParams{
        RecordPerPage: 0,
        Page: 0,
        Query: "",
    }

    result, status = u.List(params)
    if status != fiber.StatusOK {
        t.Errorf("Expected status code 200, got %d", status)
    }

    data = result["users"].([]User)
    if len(data) == 0 {
        t.Errorf("")
    }
}

func TestUpdate(t *testing.T) {
    u := User{}
    id := uuid.New()
    u.Email = id.String() + "@gmail.com"
    u.FirstName = "anurag"
    u.LastName = "pal"
    u.Password = "12345678"
    u.ConfirmPassword = "12345678"
    err := u.UniqueEmail(u.Email)
    if err != nil {
        t.Errorf("Expected email should not exist %s", u.Email)
    }
    i, iStatus := u.Add()
    if iStatus != 201 {
        t.Errorf("Expected status 201, received %d", iStatus)
    }

    userId := i["id"].(string)
    u.Id = userId
    ud := User{Id: userId}
    d, status := ud.Details()
    if status != 200 {
        t.Errorf("Expected status 200, received %d", status)
    }
    uObj := d["user"].(User)
    if uObj.Email != u.Email {
        t.Errorf("Expected email %s, received %s", u.Email, uObj.Email)
    }

    u.FirstName = "anu"
    u.Update()

    d, status = ud.Details()
    if status != 200 {
        t.Errorf("Expected status 200, received %d", status)
    }
    uObj = d["user"].(User)
    if uObj.FirstName != u.FirstName {
        t.Errorf("Expected email %s, received %s", u.FirstName, uObj.FirstName)
    }
}

func TestRemove(t *testing.T) {
    uc := User{}
    l := ListParams{}
    listing, status := uc.List(l)

    if status != 200 {
        t.Errorf("Expected status 200, received %d", status)
    }
    var total int
    total = listing["total"].(int)
    log.Println(total)
    log.Println("******")

    u := User{}
    id := uuid.New()
    u.Email = id.String() + "@gmail.com"
    u.FirstName = "anurag"
    u.LastName = "pal"
    u.Password = "12345678"
    u.ConfirmPassword = "12345678"
    err := u.UniqueEmail(u.Email)
    if err != nil {
        t.Errorf("Expected email should not exist %s", u.Email)
    }
    addedUserDetails, iStatus := u.Add()
    if iStatus != 201 {
        t.Errorf("Expected status 201, received %d", iStatus)
    }

    uc = User{}
    listing, status = uc.List(l)

    if status != 200 {
        t.Errorf("Expected status 200, received %d", status)
    }
    
    totalWithUpdateData := listing["total"].(int)
    if total != (totalWithUpdateData - 1) {
        t.Errorf("Expected records %d, but got %d", totalWithUpdateData, total)
    }

    ur := User{Id: addedUserDetails["id"].(string)}
    _, status = ur.Remove()
    if status != 200 {
        t.Errorf("Expected status 200, received %d", status)
    }

    listing, status = uc.List(l)

    if status != 200 {
        t.Errorf("Expected status 200, received %d", status)
    }

    if total != listing["total"].(int) {
        t.Errorf("Expected records %d, but got %d", total, listing["total"].(int))
    }
}

func TestInValidUser(t *testing.T) {
    id := uuid.New()
    u := User{}
    isExist := u.IsValidUser(id.String())
    if isExist {
        t.Errorf("Expected user should be false, but received true")
    }
}

func TestShouldNotUpdateInValidUser(t *testing.T) {
    id := uuid.New()
    u := User{Id: id.String(), FirstName: "anurag", LastName: "pal", Email: id.String() + "@gmail.com"}
    _, status := u.Update()
    if status != 404 {
        t.Errorf("Expected status 404, whereas returned %d", status)
    }
}

func TestShouldNotTryToRemoveInValidUser(t *testing.T) {
    id := uuid.New()
    u := User{Id: id.String()}
    _, status := u.Update()
    if status != 404 {
        t.Errorf("Expected status 404, whereas returned %d", status)
    }
}