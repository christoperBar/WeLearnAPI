package auth

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/christoperBar/WeLearnAPI/models"
	"google.golang.org/api/option"
)

func InitAtuh() *firebase.App {
	opt := option.WithCredentialsFile("/Users/alizaenalabidin/Downloads/welearngamacitra-firebase-adminsdk-ty673-7cfe5e29a3.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("error initializing app: %v\n")
	}
	if app == nil {
		panic("firebase admin init failed")
	}
	fmt.Println("sukses init fb")
	fmt.Println(app)
	return app
}

func CreateUser(
	ctx context.Context,
	client *auth.Client,
	newUser models.Student,
	userEmail string,
	userPassword string,
	userName string,
	isInstructor bool,
) *auth.UserRecord {
	// [START create_user_golang]
	params := (&auth.UserToCreate{}).
		Email(userEmail).
		EmailVerified(false).
		PhoneNumber(newUser.Phone).
		Password(userPassword).
		DisplayName(userName).
		PhotoURL(newUser.Image_url).
		Disabled(false)
	user, err := client.CreateUser(ctx, params)
	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
	}
	log.Printf("Successfully created user: %v\n", user)
	client.SetCustomUserClaims(ctx, user.UID, map[string]interface{}{"isInstructor": isInstructor})
	// [END create_user_golang]
	return user
}
