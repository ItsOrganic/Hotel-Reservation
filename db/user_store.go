package db

import (
	"Hotel_Reservation_API/types"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
const userColl = "users"
type UserStorage interface {
    GetUserById(context.Context, string) (*types.User, error)
    GetUsers(context.Context) ([]*types.User, error)
    CreateUser(context.Context, *types.User) (*types.User, error)
    DeleteUser(context.Context, string) error
    UpdateUser(ctx context.Context, filter, update bson.M) error
    Drop(context.Context) error
}

type MongoUserStore struct {
    client *mongo.Client
    coll *mongo.Collection
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter, values bson.M) error{
    update := bson.D {
        {
            "$set", values,
        },
    }
    _, err := s.coll.UpdateOne(ctx , filter, update)
    if err != nil {
        return err
    }
    return nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    _, err = s.coll.DeleteOne(ctx, bson.M{"_id":oid})
    if err != nil {
        return err
    }
    return nil
}

func NewMongoUserStore(client *mongo.Client, dbname string) *MongoUserStore {
    return &MongoUserStore{
        client: client,
        coll: client.Database(dbname).Collection(userColl),
    }
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
    fmt.Println("----------------Dropping user Collection-------------")
    return s.coll.Drop(ctx)
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
    res, err := s.coll.InsertOne(ctx, user)
    if err != nil {
        return nil, err
    }
    user.ID = res.InsertedID.(primitive.ObjectID)
    return user, nil
}


func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
    cur, err := s.coll.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    var users []*types.User
    if err := cur.All(ctx, &users); err != nil {
        return []*types.User{}, nil
    }
    return users, nil
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil{
        return nil, err
    }
    var user types.User
    if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
                return nil, err
    }
    return &user, nil
}
