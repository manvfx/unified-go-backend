package seed

import (
	"context"
	"fmt"
	"log"
	"unified-go-backend/config"
	"unified-go-backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/crypto/bcrypt"
)

func SeedData(cfg *config.Config) {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(cfg.MongoURI).SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.TODO())

	// Get collections
	usersCollection := mongoClient.Database("mdmdb").Collection("users")
	rolesCollection := mongoClient.Database("mdmdb").Collection("roles")
	permissionsCollection := mongoClient.Database("mdmdb").Collection("permissions")
	accessGroupsCollection := mongoClient.Database("mdmdb").Collection("access_groups")

	// Create permissions
	permissions := []models.Permission{
		{Name: "create_user"},
		{Name: "read_user"},
		{Name: "update_user"},
		{Name: "delete_user"},
		{Name: "list_users"},
		{Name: "create_access_group"},
		{Name: "read_access_group"},
		{Name: "update_access_group"},
		{Name: "delete_access_group"},
	}

	for _, permission := range permissions {
		_, err := permissionsCollection.InsertOne(context.TODO(), permission)
		if err != nil {
			log.Fatalf("Failed to insert permission %s: %v", permission.Name, err)
		}
	}

	// Create roles
	roles := []models.Role{
		{
			ID:          primitive.NewObjectID(),
			Name:        "admin",
			Permissions: []string{"create_user", "read_user", "update_user", "delete_user", "list_users", "create_access_group", "read_access_group", "update_access_group", "delete_access_group"},
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "operator",
			Permissions: []string{"read_user", "update_user", "list_users"},
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "user",
			Permissions: []string{"read_user"},
		},
	}

	for _, role := range roles {
		_, err := rolesCollection.InsertOne(context.TODO(), role)
		if err != nil {
			log.Fatalf("Failed to insert role %s: %v", role.Name, err)
		}
	}

	// Create access groups
	accessGroups := []models.AccessGroup{
		{
			ID:          primitive.NewObjectID(),
			Name:        "admin_group",
			Roles:       []string{"admin"},
			Permissions: []string{},
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "operator_group",
			Roles:       []string{"operator"},
			Permissions: []string{},
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "user_group",
			Roles:       []string{"user"},
			Permissions: []string{},
		},
	}

	for _, group := range accessGroups {
		_, err := accessGroupsCollection.InsertOne(context.TODO(), group)
		if err != nil {
			log.Fatalf("Failed to insert access group %s: %v", group.Name, err)
		}
	}

	// Create hashed password for dummy users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create users
	users := []models.User{
		{
			ID:          primitive.NewObjectID(),
			Email:       "admin@example.com",
			Username:    "admin",
			Password:    string(hashedPassword),
			Verified:    true,
			Roles:       []string{"admin"},
			AccessGroup: "admin_group",
		},
		{
			ID:          primitive.NewObjectID(),
			Email:       "operator@example.com",
			Username:    "operator",
			Password:    string(hashedPassword),
			Verified:    true,
			Roles:       []string{"operator"},
			AccessGroup: "operator_group",
		},
		{
			ID:          primitive.NewObjectID(),
			Email:       "user@example.com",
			Username:    "user",
			Password:    string(hashedPassword),
			Verified:    true,
			Roles:       []string{"user"},
			AccessGroup: "user_group",
		},
	}

	for _, user := range users {
		_, err := usersCollection.InsertOne(context.TODO(), user)
		if err != nil {
			log.Fatalf("Failed to insert user %s: %v", user.Username, err)
		}
	}

	fmt.Println("Dummy data inserted successfully!")
}
