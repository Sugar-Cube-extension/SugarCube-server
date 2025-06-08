package services

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CleanupNegativeCoupons(db *mongo.Database) error {
	ctx := context.Background()

	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to list collections")
		return err
	}

	for _, colName := range collections {
		col := db.Collection(colName)

		filter := bson.M{"score": bson.M{"$lt": 0}}
		result, err := col.DeleteMany(ctx, filter)
		if err != nil {
			log.Error().
				Err(err).
				Str("collection", colName).
				Msg("Failed to delete negative-score coupons")
			continue
		}

		if result.DeletedCount > 0 {
			log.Info().
				Int64("deleted", result.DeletedCount).
				Str("collection", colName).
				Msg("Removed negative-score coupons")
		}
	}

	return nil
}

func StartNegativeScorePruner(db *mongo.Database) {
	s := gocron.NewScheduler(time.UTC)
	s.Every(6).Hours().Do(func() {
		log.Info().Msg("Starting scheduled cleanup of negative-score coupons")

		if err := CleanupNegativeCoupons(db); err != nil {
			log.Error().Err(err).Msg("Cleanup job failed")
		} else {
			log.Info().Msg("Cleanup job completed successfully")
		}
	})
	s.StartAsync()
}
