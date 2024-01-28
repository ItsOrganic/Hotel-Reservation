package api

import (
	"Hotel_Reservation_API/db"
	"Hotel_Reservation_API/types"
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
    userStore db.UserStorage
}

func NewUserHandler (userStore db.UserStorage) *UserHandler {
    return &UserHandler{
        userStore: userStore,
    }
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
    userID := c.Params("id")
    if err := h.userStore.DeleteUser(c.Context(), userID);  err != nil {
        return err
    }
    return c.JSON(map[string]string{"deleted": userID})
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
    var (
        values bson.M
        userID = c.Params("id")
    )
    oid, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return err
    }
    if err := c.BodyParser(&values); err != nil {
        return err
    }
    filter := bson.M{"_id": oid}
    if err := h.userStore.UpdateUser(c.Context(), filter, values); err != nil {
        return err
    }
    return c.JSON(map[string]string{"updated": userID})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
    var params types.CreateUserParams
    if err := c.BodyParser(&params); err != nil {
        return err
    }
    if errors := params.Validate(); len(errors) > 0 {
        return c.JSON(errors)
    }
    user, err := types.NewUserFromParams(params)
    if err != nil {
        return err
    }
    insertedUser, err := h.userStore.CreateUser(c.Context(), user)
    if err != nil{
         return err
}
    return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
    var (
        id = c.Params("id")
        )
    user, err := h.userStore.GetUserById(c.Context(),id)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.JSON(map[string]string{"error": "not found or already deleted"})
        }
        return err
    }
    return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
        users, err := h.userStore.GetUsers(c.Context())
        if err != nil {
            return err
        }
        return c.JSON(users)
}