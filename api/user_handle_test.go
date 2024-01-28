package api

import (
	"Hotel_Reservation_API/db"
	"Hotel_Reservation_API/types"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
const (
    testdburi = "mongodb://localhost:27017"
    dbname = "hotel-reservation-test"
)
type testdb struct {
    db.UserStorage
}

func (tdb *testdb) teardown(t *testing.T){
    if  err := tdb.UserStorage.Drop(context.TODO()); err != nil {
        t.Fatal(err)
    }
}

func setup(t *testing.T) *testdb {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
    if err != nil {
        log.Fatal(err)
    }
    return &testdb{
        UserStorage:db.NewMongoUserStore(client, dbname),
    }
}
func TestPostUser(t *testing.T) {
    tdb := setup(t)
    defer tdb.teardown(t)

    app := fiber.New()
    userHandler := NewUserHandler(tdb.UserStorage)
    app.Post("/",userHandler.HandlePostUser)

    params := types.CreateUserParams{
        Email: "abcd@abcd.com",
        FirstName: "ABCD",
        LastName: "XYZ",
        Password: "Password",
    }
    b, _ := json.Marshal(params)

    req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
    req.Header.Add("Content-Type", "application/json")
    resp, err := app.Test(req)

    if err != nil {
        t.Error(err)
    }
    var user types.User
    json.NewDecoder(resp.Body).Decode(&user)
    if user.FirstName != params.FirstName {
        t.Errorf("expected FirstName %s but got %s",params.FirstName,user.FirstName)
    }
    if user.LastName != params.LastName {
        t.Errorf("expected LastName %s but got %s",params.LastName,user.LastName)
    }
    if user.Email != params.Email {
        t.Errorf("expected Email  %s but got %s",params.Email,user.Email)
    }
}
