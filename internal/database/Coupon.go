package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CouponEntry struct {
	Coupon    string         `bson:"coupon" json:"coupon"`           //Actual coupon
	Score     int            `bson:"score" json:"score"`             //Internal, a score to track how effective the code is
	ExpiresAt time.Time      `bson:"expires_at" json:"expires_at"`   //IMPORTANT: ISO8601-formatted
	Extra     map[string]any `bson:",inline" json:"extra,omitempty"` //Random stuff for other sites
}

type Site struct {
	Name          string        `json:"name"` //URL
	CouponEntries []CouponEntry `json:"coupon_entries"`
}

func GetSiteStruct(siteName string, db *mongo.Database) (*Site, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collections, err := db.ListCollectionNames(ctx, bson.M{"name": siteName})
	if err != nil {
		return nil, fmt.Errorf("error listing collections: %w", err)
	}
	if len(collections) == 0 {
		return nil, fmt.Errorf("site '%s' does not exist", siteName)
	}

	coll := db.Collection(siteName)
	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching coupons from '%s': %w", siteName, err)
	}
	defer cur.Close(ctx)

	var coupons []CouponEntry
	for cur.Next(ctx) {
		var entry CouponEntry
		if err := cur.Decode(&entry); err != nil {
			return nil, fmt.Errorf("decode error: %w", err)
		}
		coupons = append(coupons, entry)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return &Site{
		Name:          siteName,
		CouponEntries: coupons,
	}, nil

}

func AddCouponToExistingSite(siteName string, coupon CouponEntry, db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collections, err := db.ListCollectionNames(ctx, bson.M{"name": siteName})
	if err != nil {
		return fmt.Errorf("error listing collections: %w", err)
	}
	if len(collections) == 0 {
		return fmt.Errorf("site '%s' does not exist", siteName)
	}

	filter := bson.M{"coupon": coupon.Coupon}
	var existing bson.M
	err = db.Collection(siteName).FindOne(ctx, filter).Decode(&existing)
	if err == nil {
		return fmt.Errorf("coupon '%s' already exists", coupon.Coupon)
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("failed to check existing coupon: %w", err)
	}

	_, err = db.Collection(siteName).InsertOne(ctx, coupon)
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}

func AddSite(siteName string, db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collections, err := db.ListCollectionNames(ctx, bson.M{"name": siteName})
	if err != nil {
		return fmt.Errorf("error listing collections: %w", err)
	}
	if len(collections) != 0 {
		return fmt.Errorf("site '%s' does already exists!", siteName)
	}
	e := db.CreateCollection(ctx, siteName)
	if e != nil {
		return fmt.Errorf("site collection for '%s' failed: %w ", siteName, err)

	}
	indexErr := EnsureCouponIndex(db.Collection(siteName))
	if indexErr != nil {
		return fmt.Errorf("collection index adjustment for site '%s' failed: %w", siteName, indexErr)
	}

	return nil

}

func ProcessCallback(db *mongo.Database, siteName string, callbackResults map[string]bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection(siteName)

	for code, worked := range callbackResults {
		if !worked {
			filter := bson.M{"coupon": code}
			update := bson.M{"$inc": bson.M{"score": -1}}
			_, _ = coll.UpdateOne(ctx, filter, update)
		}
	}
}

func EnsureCouponIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "coupon", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetName("coupon_idx"),
	}

	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}
