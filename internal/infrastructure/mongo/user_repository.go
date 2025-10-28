package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend-challenge/internal/application"
	"backend-challenge/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const usersCollection = "users"

// UserRepository is a Mongo-backed implementation of application.UserRepository.
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository constructs a repository and sets indices.
func NewUserRepository(db *mongo.Database) (*UserRepository, error) {
	col := db.Collection(usersCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("unique_email"),
	}

	if _, err := col.Indexes().CreateOne(ctx, indexModel); err != nil {
		return nil, fmt.Errorf("create email index: %w", err)
	}

	return &UserRepository{collection: col}, nil
}

type mongoUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
}

func toDomain(mu mongoUser) domain.User {
	return domain.User{
		ID:        mu.ID.Hex(),
		Name:      mu.Name,
		Email:     mu.Email,
		Password:  mu.Password,
		CreatedAt: mu.CreatedAt,
	}
}

func fromDomain(u domain.User) mongoUser {
	var id primitive.ObjectID
	if u.ID != "" {
		if oid, err := primitive.ObjectIDFromHex(u.ID); err == nil {
			id = oid
		}
	}
	return mongoUser{
		ID:        id,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

func parseID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, application.ErrNotFound
	}
	return oid, nil
}

// Create persists a new user document.
func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	doc := fromDomain(user)
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.User{}, application.ErrDuplicateEmail
		}
		return domain.User{}, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return domain.User{}, errors.New("unexpected inserted id type")
	}
	user.ID = id.Hex()
	return user, nil
}

// GetByEmail retrieves a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var mu mongoUser
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&mu)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, application.ErrNotFound
		}
		return domain.User{}, err
	}
	return toDomain(mu), nil
}

// GetByID retrieves a user by id.
func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	oid, err := parseID(id)
	if err != nil {
		return domain.User{}, err
	}

	var mu mongoUser
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&mu)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, application.ErrNotFound
		}
		return domain.User{}, err
	}
	return toDomain(mu), nil
}

// List returns all users sorted by creation time descending.
func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	cursor, err := r.collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []domain.User
	for cursor.Next(ctx) {
		var mu mongoUser
		if err := cursor.Decode(&mu); err != nil {
			return nil, err
		}
		users = append(users, toDomain(mu))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Update modifies email and/or name for a user.
func (r *UserRepository) Update(ctx context.Context, id string, update domain.UpdateUser) (domain.User, error) {
	oid, err := parseID(id)
	if err != nil {
		return domain.User{}, err
	}

	set := bson.M{}
	if update.Name != nil {
		set["name"] = *update.Name
	}
	if update.Email != nil {
		set["email"] = *update.Email
	}

	if len(set) == 0 {
		return domain.User{}, application.ErrNoFieldsToUpdate
	}

	var mu mongoUser
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": set}, opts).Decode(&mu)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, application.ErrNotFound
		}
		if mongo.IsDuplicateKeyError(err) {
			return domain.User{}, application.ErrDuplicateEmail
		}
		return domain.User{}, err
	}

	return toDomain(mu), nil
}

// Delete removes a user.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	oid, err := parseID(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return application.ErrNotFound
	}
	return nil
}

// Count returns the total number of users.
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.D{})
}
