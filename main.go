package main

import (
	"Hotel_Reservation_API/api"
	"Hotel_Reservation_API/db"
	"context"
	"flag"
	"log"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        return c.JSON(map[string]string{"error": err.Error()})
    },
}
func main() {
    listenAddress := flag.String("listenAddress",":5000","This is the port listen address")
    flag.Parse()

    client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(db.DBURI))
    if err != nil {
        log.Fatal(err)
    }

    //handlers initialization
    var (
    userHandler = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
    hotelStore = db.NewMongoHotelStore(client)
    roomStore = db.NewMongoRoomStore(client, hotelStore)
    hotelHandler = api.NewHotelHandler(hotelStore, roomStore)
    app = fiber.New(config)
    apiv1 = app.Group("/api/v1/")
)

    //User handlers
    app.Get("/foo",handleFoo)
    apiv1.Put("/user/:id", userHandler.HandlePutUser)
    apiv1.Post("/user", userHandler.HandlePostUser)
    apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
    apiv1.Get("/user", userHandler.HandleGetUsers)
    apiv1.Get("/user/:id", userHandler.HandleGetUser)

    //Hotel handlers
    apiv1.Get("/hotel",hotelHandler.HandleGetHotels)
   app.Listen(*listenAddress)

}
    func handleFoo(c *fiber.Ctx) error {
        return c.JSON(map[string]string{"msg":"Working just fine"})
    }
