package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yildiz-fatih/gopaste/internal/models"
)

func (app *application) getPaste(ctx context.Context, slug string) (*models.Paste, error) {
	cacheHit := false
	var p models.Paste

	// check redis cache first
	cachedPaste, err := app.redisClient.Get(ctx, fmt.Sprintf("paste:%s", slug))
	if err != nil && err != redis.Nil {
		// something went wrong with redis, and it's not a cache miss
		app.logger.Error("redis get error", "error", err)
	} else if err == nil {
		// cache hit!
		err = json.Unmarshal([]byte(cachedPaste), &p)
		if err != nil {
			app.logger.Error("redis unmarshal error", "error", err)
		} else {
			cacheHit = true
		}
	}

	if !cacheHit {
		// get from database
		p, err = app.pasteModel.Get(slug)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return nil, models.ErrNotFound
			} else {
				return nil, err
			}
		}

		// write to redis cache
		pasteJson, err := json.Marshal(p)
		if err != nil {
			app.logger.Error("redis marshal error", "error", err)
		} else {
			err = app.redisClient.Set(ctx, fmt.Sprintf("paste:%s", p.Slug), pasteJson, time.Until(p.Expires))
			if err != nil {
				app.logger.Error("redis set error", "error", err)
			}
		}
	}

	return &p, nil
}
