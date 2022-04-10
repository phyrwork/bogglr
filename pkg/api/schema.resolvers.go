package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/phyrwork/bogglr/pkg/api/generated"
	"github.com/phyrwork/bogglr/pkg/api/model"
	"github.com/phyrwork/bogglr/pkg/database"
	"gorm.io/gorm"
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

func (r *mutationResolver) CreatePlayer(ctx context.Context, name string) (*model.Player, error) {
	record := database.Player{Name: name}
	if err := r.DB.WithContext(ctx).Create(&record).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &model.Player{
		ID:   strconv.Itoa(record.ID),
		Name: record.Name,
	}, nil
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

func (r *mutationResolver) CreateWord(ctx context.Context, gameID string, path []model.Point) (*model.Word, error) {
	game, err := r.Query().Game(ctx, gameID)
	if err != nil {
		return nil, err
	}
	var record database.Word
	record.GameID, err = strconv.Atoi(game.ID)
	record.Path = MapOf(path, func(point model.Point) database.Point {
		return database.Point(point)
	})
	if err := r.DB.WithContext(ctx).Create(&record).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &model.Word{
		ID:   strconv.Itoa(record.ID),
		Game: game,
		Path: path,
	}, nil
}

func (r *playerResolver) Words(ctx context.Context, obj *model.Player) ([]*model.Word, error) {
	var (
		record database.Player
		err    error
	)
	record.ID, err = strconv.Atoi(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid word id '%s': %w", obj.ID, err)
	}
	err = r.DB.WithContext(ctx).Preload("Words").Find(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("player '%s' not found", obj.ID)
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return MapPointersOf(record.Words, func(record database.Word) model.Word {
		return model.Word{
			ID: strconv.Itoa(record.ID),
			Path: MapOf(record.Path, func(record database.Point) model.Point {
				return model.Point(record)
			}),
		}
	}), nil
}

func (r *queryResolver) Player(ctx context.Context, id string) (*model.Player, error) {
	var (
		record database.Player
		err    error
	)
	record.ID, err = strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid player id '%s': %w", id, err)
	}
	err = r.DB.WithContext(ctx).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("player '%s' not found", id)
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &model.Player{
		ID:   strconv.Itoa(record.ID),
		Name: record.Name,
	}, nil
}

func (r *queryResolver) Players(ctx context.Context, first *int, after *string) (*model.PlayersConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Game(ctx context.Context, id string) (*model.Game, error) {
	var (
		record database.Game
		err    error
	)
	record.ID, err = strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid game id '%s': %w", id, err)
	}
	err = r.DB.WithContext(ctx).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("game '%s' not found", id)
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &model.Game{
		ID:    strconv.Itoa(record.ID),
		Board: record.Board,
	}, nil
}

func (r *queryResolver) Games(ctx context.Context, first *int, after *string) (*model.GamesConnection, error) {
	qry := r.DB.WithContext(ctx)
	if after != nil {
		startCursor, err := strconv.Atoi(*after)
		if err != nil {
			return nil, fmt.Errorf("invalid start cursor '%s': %w", *after, err)
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

func (r *queryResolver) Words(ctx context.Context, gameID *string, playerID *string, first *int, after *string) (*model.WordsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *wordResolver) Game(ctx context.Context, obj *model.Word) (*model.Game, error) {
	var (
		record database.Word
		err    error
	)
	record.ID, err = strconv.Atoi(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid word id '%s': %w", obj.ID, err)
	}
	err = r.DB.WithContext(ctx).Preload("Game").Find(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("word '%s' not found", obj.ID)
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &model.Game{
		ID:    strconv.Itoa(record.Game.ID),
		Board: record.Game.Board,
	}, nil
}

func (r *wordResolver) Players(ctx context.Context, obj *model.Word) ([]*model.Player, error) {
	var (
		record database.Word
		err    error
	)
	record.ID, err = strconv.Atoi(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid word id '%s': %w", obj.ID, err)
	}
	err = r.DB.WithContext(ctx).Preload("Players").Find(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("word '%s' not found", obj.ID)
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return MapPointersOf(record.Players, func(record database.Player) model.Player {
		return model.Player{
			ID:   strconv.Itoa(record.ID),
			Name: record.Name,
		}
	}), nil
}

// Game returns generated.GameResolver implementation.
func (r *Resolver) Game() generated.GameResolver { return &gameResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Player returns generated.PlayerResolver implementation.
func (r *Resolver) Player() generated.PlayerResolver { return &playerResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Word returns generated.WordResolver implementation.
func (r *Resolver) Word() generated.WordResolver { return &wordResolver{r} }

type gameResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type playerResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type wordResolver struct{ *Resolver }
