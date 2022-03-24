package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/phyrwork/bogglr/pkg/api/generated"
	"github.com/phyrwork/bogglr/pkg/api/model"
	"github.com/phyrwork/bogglr/pkg/database"
)

func (r *gameResolver) Board(ctx context.Context, obj *model.Game) ([]string, error) {
	if obj.Board != nil {
		return obj.Board, nil
	}
	if err := r.DB.WithContext(ctx).Find(obj).Error; err != nil {
		return nil, err
	}
	return obj.Board, nil
}

func (r *mutationResolver) CreateGame(ctx context.Context, board []string) (*model.Game, error) {
	tiles := model.Board(board).Dump()
	if !tiles.IsRect() {
		w, h := tiles.Dims()
		return nil, fmt.Errorf("board must be rectangular: is %d x %v", w, h)
	}
	var record database.Game
	record.Board = board
	if err := r.DB.WithContext(ctx).Create(&record).Error; err != nil {
		switch {
		case strings.Contains(err.Error(), "value too long"): // TODO: be more specific.
			w, h := tiles.Dims()
			return nil, fmt.Errorf("board is too wide: is %d x %v", w, h)
		case strings.Contains(err.Error(), "violates check"): // TODO: be more specific.
			w, h := tiles.Dims()
			return nil, fmt.Errorf("board is too tall: is %d x %v", w, h)
		default:
			return nil, fmt.Errorf("database error: %w", err)
		}
	}
	obj := model.Game{
		ID:    strconv.Itoa(int(record.ID)),
		Board: board,
	}
	return &obj, nil
}

func (r *queryResolver) Games(ctx context.Context, first *int, after *string) (*model.GamesConnection, error) {
	qry := r.DB.WithContext(ctx)
	if after != nil {
		startCursor, err := strconv.Atoi(*after)
		if err != nil {
			return nil, fmt.Errorf("invalid start cursor %s: %w", *after, err)
		}
		qry = qry.Where("id > ?", startCursor)
	}
	if first != nil {
		qry = qry.Limit(*first)
	}
	qry = qry.Order("id asc")

	var records []database.Game
	if err := qry.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	var edges []*model.Game = MapPointersOf(records, func(record database.Game) model.Game {
		return model.Game{
			ID:    strconv.Itoa(int(record.ID)),
			Board: record.Board,
		}
	})
	if len(edges) == 0 {
		return nil, nil
	}

	pageInfo := model.PageInfo{
		StartCursor: edges[0].ID,
		EndCursor:   edges[len(edges)-1].ID,
	}
	qry = r.DB.WithContext(ctx).
		Where("id > ?", records[len(records)-1].ID).
		Limit(1)
	if err := qry.Find(&records).Error; err != nil {
		log.Print(fmt.Errorf("database error: %w", err))
	} else {
		hasNextPage := len(records) > 0
		pageInfo.HasNextPage = &hasNextPage
	}

	return &model.GamesConnection{
		Edges:    edges,
		PageInfo: &pageInfo,
	}, nil
}

// Game returns generated.GameResolver implementation.
func (r *Resolver) Game() generated.GameResolver { return &gameResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type gameResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
