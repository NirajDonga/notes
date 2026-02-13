package notes

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	coll *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		coll: db.Collection("notes"),
	}
}

func (r *Repo) Create(ctx context.Context, note Note) (Note, error) {
	childCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.coll.InsertOne(childCtx, note)
	if err != nil {
		return Note{}, fmt.Errorf("Insert note failed: %w", err)
	}

	return note, nil
}

func (r *Repo) List(ctx context.Context) ([]Note, error) {
	childCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{}

	//find return a cursor like iterator over matching elements
	cursor, err := r.coll.Find(childCtx, filter)
	if err != nil {
		return nil, fmt.Errorf("find notes Failed: %w", err)
	}
	defer cursor.Close(childCtx)

	var notes []Note
	if err := cursor.All(childCtx, &notes); err != nil {
		return nil, fmt.Errorf("Decode Nodes Failed: %w", err)
	}

	return notes, nil
}

func (r *Repo) GetByID(ctx context.Context, noteId primitive.ObjectID) (Note, error) {
	childCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": noteId}

	var note Note
	// options give you features for sort, sum, projection of mongodb
	err := r.coll.FindOne(childCtx, filter, options.FindOne()).Decode(&note)
	if err != nil {
		return Note{}, fmt.Errorf("Find note by id failed: %w", err)
	}
	return note, nil
}

func (r *Repo) updateByID(ctx context.Context, noteId primitive.ObjectID, req UpdateNoteRequest) (Note, error) {
	childCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": noteId}

	update := bson.M{
		"$set": bson.M{
			"title":     req.Title,
			"content":   req.Content,
			"pinned":    req.Pinned,
			"updatedAt": time.Now().UTC(),
		},
	}

	after := options.After
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	var updated Note
	err := r.coll.FindOneAndUpdate(childCtx, filter, update, &opts).Decode(&updated)
	if err != nil {
		return Note{}, fmt.Errorf("Update note failed: %w", err)
	}
	return updated, nil

}
