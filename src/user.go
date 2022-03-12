package soccer

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type UserClient struct {
	Name    string    `json:"name"`
	Expired time.Time `json:"expired"`

	Token        string    `json:"-"`
	RegisterTime time.Time `json:"-"`
}

func (uc UserClient) ToJSON() []byte {
	a, _ := json.Marshal(uc)
	return a
}

func (s *Soccer) GetUserFromToken(ctx context.Context, token string) (u UserClient, err error) {
	response := s.Mongo.Database("users").Collection("user").FindOne(ctx, bson.M{
		"token": token,
	})
	err = response.Decode(&u)
	return
}

func (s *Soccer) SetUser(ctx context.Context, u UserClient) error {
	_, err := s.Mongo.Database("users").Collection("user").InsertOne(ctx, u)
	return err
}
