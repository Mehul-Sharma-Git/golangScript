package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"script/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Customers struct {
	CustomerID      primitive.ObjectID `bson:"_id,omitempty" json:"customer_id"`
	Name            string             `bson:"name" json:"name"`
	Email           string             `bson:"email" json:"email"`
	MojoUserId      primitive.ObjectID `bson:"ma_user_id" json:"ma_user_id"`
	Role            string             `bson:"role" json:"role"`
	FirstLogin      bool               `bson:"first_login" json:"first_login"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	Status          string             `bson:"status" json:"status"`
	ExpireAt        time.Time          `bson:"expire_at" json:"expire_at"`
	OnboardingLists []string           `bson:"onboarding_lists" json:"onboarding_lists"`
	OrganizationID  primitive.ObjectID `bson:"organization_id" json:"organization_id"`
}

type Organization struct {
	Id           primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	Name         string             `bson:"name" json:"name"`
	Slug         string             `bson:"slug" json:"slug"`
	Projectid    primitive.ObjectID `bson:"project_id" json:"project_id"`
	IsActive     bool               `bson:"is_active" json:"is_active"`
	CreatedDate  time.Time          `bson:"created_at" json:"created_at"`
	ModifiedDate time.Time          `bson:"modified_at" json:"modified_at"`
}

type Projects struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"project_id"`
	TestProject     primitive.ObjectID     `bson:"test_project,omitempty" json:"test_project,omitempty"`
	Name            string                 `bson:"name" json:"name"`
	Credentials     []Credential           `bson:"credentials" json:"credentials"`
	URL             []string               `bson:"url" json:"url"`
	CustomerID      primitive.ObjectID     `bson:"customer_id" json:"customer_id"`
	Branding        map[string]interface{} `bson:"branding" json:"branding"`
	WebAuthn        map[string]interface{} `bson:"webauthn" json:"webauthn"`
	Integrations    []Integration          `bson:"integrations" json:"integrations"`
	CreatedAt       time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	Status          string                 `bson:"status" json:"status"`
	ExpireAt        *time.Time             `bson:"expire_at,omitempty" json:"expire_at,omitempty"`
	DataStorage     map[string]interface{} `bson:"data_storage" json:"data_storage"`
	JWTconfig       map[string]interface{} `bson:"jwt_config" json:"jwt_config"`
	Settings        map[string]interface{} `bson:"settings" json:"settings"`
	OauthIdentifier string                 `bson:"oauth_identifier" json:"oauth_identifier"`
	ResourceID      primitive.ObjectID     `bson:"resource_id" json:"resource_id"`
}
type Integration struct {
	Name   string                 `bson:"name" json:"name"`
	Slug   string                 `bson:"slug" json:"slug"`
	Status string                 `bson:"status" json:"status"`
	Label  string                 `bson:"label" json:"label"`
	Config map[string]interface{} `bson:"config" json:"config"`
}

type Credential struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"credential_id"`
	Name      string             `bson:"name" json:"name"`
	APIKey    string             `bson:"api_key" json:"api_key"`
	APISecret string             `bson:"api_secret" json:"api_secret"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Resource struct {
	Id             primitive.ObjectID `bson:"resource_id" json:"resource_id"`
	Name           string             `bson:"resource_name" json:"resource_name"`
	Slug           string             `bson:"resource_slug" json:"resource_slug"`
	OrganizationId primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	IsActive       bool               `bson:"is_active" json:"is_active"`
	CreatedDate    time.Time          `bson:"created_at" json:"created_at"`
	ModifiedDate   time.Time          `bson:"modified_at" json:"modified_at"`
}

func main() {

	AddOrganization()
}

func init_DB() *mongo.Client {

	var ctx = context.TODO()

	// Client - mongodb client
	var Client *mongo.Client
	var err error

	DBConf := config.App
	clientOptions := options.Client().ApplyURI(DBConf.DB_URI)
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return Client
}
func AddOrganization() {

	DBConf := config.App
	var Client *mongo.Client
	Client = init_DB()
	db := Client.Database(DBConf.DB_Name).Collection("ma_customers")

	//Define an array in which you can store the decoded documents
	var results []Customers

	//Passing the bson.D{{}} as the filter matches  documents in the collection
	{

	}
	cur, err := db.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	//Finding multiple documents returns a cursor
	//Iterate through the cursor allows us to decode documents one at a time
	var count int
	count = 0
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem Customers
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		if elem.OrganizationID.Hex() == "000000000000000000000000" {

			email := map[string]string{"name": elem.Email}
			jsonReq, err := json.Marshal(email)
			client := &http.Client{}
			req, err := http.NewRequest("POST", "ORGResuest", bytes.NewBuffer(jsonReq))
			req.SetBasicAuth(config.App.OrgApiKey, config.App.OrgSecret)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()
			bodyBytes, _ := ioutil.ReadAll(resp.Body)

			// Convert response body to string
			bodyString := string(bodyBytes)
			fmt.Println(bodyString)
			var OrgResponse Organization
			json.Unmarshal(bodyBytes, &OrgResponse)
			fmt.Printf("%+v\n", OrgResponse)
			update := bson.M{"$set": bson.M{
				"organization_id": OrgResponse.Id,
			}}
			result, err := db.UpdateOne(
				bson.M{"_id": elem.CustomerID},
				update,
			)
			if err != nil {
				log.Fatal(err)
			}
			// Convert response body to Todo struct

		}

		results = append(results, elem)
		fmt.Println(count)
		fmt.Println(elem.Email)
		fmt.Println(elem.OrganizationID)
		count = count + 1
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	// fmt.Printf("Found multiple documents: %+v\n", results)

}

func AddResources() {
	DBConf := config.App
	var Client *mongo.Client
	Client = init_DB()
	db := Client.Database(DBConf.DB_Name).Collection("ma_projects")

	var results []Projects

	//Passing the bson.D{{}} as the filter matches  documents in the collection
	{

	}
	cur, err := db.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	//Finding multiple documents returns a cursor
	//Iterate through the cursor allows us to decode documents one at a time
	var count int
	count = 0
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var projects Projects
		err := cur.Decode(&projects)
		if err != nil {
			log.Fatal(err)
		}

		if projects.ResourceID.Hex() == "000000000000000000000000" {
			db_org := Client.Database(DBConf.DB_Name).Collection("ma_customers")

			var customer Customers

			cursor, err := db_org.FindOne(context.TODO(), bson.D{"_id": projects.CustomerID})
			if err != nil {
				log.Fatal(err)
			}

			err = cursor.Decode(&customer)
			if err != nil {
				log.Fatal(err)
			}

			project_Name := map[string]string{"name": projects.Name}
			jsonReq, err := json.Marshal(project_Name)
			client := &http.Client{}
			req, err := http.NewRequest("POST", fmt.Sprintf("https://url/organisation/%s/resources", customer.OrganizationID), bytes.NewBuffer(jsonReq))
			req.SetBasicAuth(config.App.OrgApiKey, config.App.OrgSecret)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()
			bodyBytes, _ := ioutil.ReadAll(resp.Body)

			// Convert response body to string
			bodyString := string(bodyBytes)
			fmt.Println(bodyString)
			var ResourceResponse Resource
			json.Unmarshal(bodyBytes, &ResourceResponse)
			fmt.Printf("%+v\n", ResourceResponse)
			update := bson.M{"$set": bson.M{
				"resource_id":     ResourceResponse.Id,
				"organization_id": customer.OrganizationID,
			}}
			result, err := db.UpdateOne(
				bson.M{"_id": projects.ID},
				update,
			)
			if err != nil {
				log.Fatal(err)
			}
			// Convert response body to Todo struct

		}

		results = append(results, projects)
		fmt.Println(count)
		fmt.Println(projects.Name)
		fmt.Println(projects.ResourceID)
		count = count + 1
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	// fmt.Printf("Found multiple documents: %+v\n", results)

}
