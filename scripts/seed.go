package main

import (
	"Hotel_Reservation_API/db"
	"Hotel_Reservation_API/types"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
    client *mongo.Client
    roomStore db.RoomStore
    hotelStore db.HotelStore
    ctx = context.Background()
)

func seedHotel(name string, location string, rating int) {
     hotel := types.Hotel{
        Name: name,
        Location: location,
        Rooms: []primitive.ObjectID{},
        Rating: rating,
    }
    rooms := []types.Room{
        {
        Type: types.SinglePersonRoomType,
        BasePrice: 99.9,
    },
    {
        Type: types.DeluxeRoomType,
        BasePrice: 199.9,
    },
    {
         Type: types.SeaSideRoomType,
        BasePrice: 159.9,
    },
}
    insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
    if err != nil {
        log.Fatal(err)
    }
    for _, room := range rooms {
    room.HotelID = insertedHotel.ID
    _, err := roomStore.InsertRoom(ctx, &room)
    if err!= nil {
        log.Fatal(err)
    }
}
}

func main() {
    seedHotel("Hotel Palm", "Chennai", 5)
    seedHotel("Hotel Royal Blue", "Kerala", 4)
    seedHotel("Royal Robbers", "London", 3)
}

func init(){
    var err error
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
    if err != nil {
        log.Fatal(err)
    }
    hotelStore = db.NewMongoHotelStore(client)
    roomStore = db.NewMongoRoomStore(client, hotelStore)


}

