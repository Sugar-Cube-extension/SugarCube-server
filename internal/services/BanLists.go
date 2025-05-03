package services

import (
	"bufio"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func StartBlocklistUpdater(db *mongo.Database) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(12).Hours().Do(func() {
		log.Info().Msg("Updating IP blocklist...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		UpdateIPBlocklist(ctx, db)
	})

	s.StartAsync()
}

func UpdateIPBlocklist(ctx context.Context, db *mongo.Database) {
	sources := []string{
		//"https://raw.githubusercontent.com/firehol/blocklist-ipsets/master/firehol_level1.netset",
		"https://lists.blocklist.de/lists/all.txt",
	}
	for _, url := range sources {
		resp, err := http.Get(url)
		if err != nil {
			log.Warn().Str("source", url).Msg("Failed to fetch blocklist")
			continue
		}
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		var ips []any
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				ips = append(ips, bson.M{"ip": line})
			}
		}

		result, err := db.Collection("ip_bans").InsertMany(ctx, ips, options.InsertMany().SetOrdered(false))
		if err != nil {
			if bwe, ok := err.(mongo.BulkWriteException); ok {
				for _, we := range bwe.WriteErrors {
					if we.Code != 11000 {
						log.Error().Err(err).Msg("Non-duplicate write error in blocklist")
						break
					}
				}
				log.Info().Int("inserted", len(result.InsertedIDs)).Msg("Blocklist update completed with duplicates")
			} else {
				log.Error().Err(err).Msg("Failed to insert blocklist")
			}
		} else {
			log.Info().Int("inserted", len(result.InsertedIDs)).Msg("Successfully inserted IPs into blocklist")
		}
	}
}

func InitBanLists(db *mongo.Database, ctx context.Context) {
	db.Collection("ip_bans").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "ip", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	StartBlocklistUpdater(db)
}
