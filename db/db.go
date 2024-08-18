package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rqlite/gorqlite"
)

func New(conn *gorqlite.Connection) *Queries {
	return &Queries{conn: conn}
}

type Queries struct {
	conn *gorqlite.Connection
}

type DocumentsInsertArgs struct {
	Name      string
	Content   string
	Embedding []float64
}

func (q *Queries) DocumentsInsert(ctx context.Context, args DocumentsInsertArgs) (id int64, err error) {
	embeddingJSON, err := json.Marshal(args.Embedding)
	if err != nil {
		return id, fmt.Errorf("failed to marshal embedding: %w", err)
	}
	statements := []gorqlite.ParameterizedStatement{
		{
			Query:     "INSERT INTO documents (name, content) VALUES (?, ?)",
			Arguments: []any{args.Name, args.Content},
		},
		{
			Query:     "INSERT INTO embeddings (id, embedding) VALUES (last_insert_rowid(), ?)",
			Arguments: []any{string(embeddingJSON)},
		},
	}
	if err := q.conn.SetExecutionWithTransaction(true); err != nil {
		return id, err
	}
	results, err := q.conn.WriteParameterizedContext(ctx, statements)
	if err != nil {
		return id, err
	}
	return results[0].LastInsertID, nil
}

type Document struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Content   string    `db:"content"`
	Embedding []float64 `db:"embedding"`
}

func (q *Queries) DocumentsGet(ctx context.Context, id int64) (doc Document, err error) {
	query := gorqlite.ParameterizedStatement{
		Query: `SELECT
							d.id, d.name, d.content, e.embedding
						FROM
							documents d
								INNER JOIN 
									embeddings e ON d.id = e.id
						WHERE d.id = ?`,
		Arguments: []any{id},
	}
	result, err := q.conn.QueryOneParameterized(query)
	if err != nil {
		return doc, err
	}
	for result.Next() {
		if err = result.Scan(&doc.ID, &doc.Name, &doc.Content, &doc.Embedding); err != nil {
			return doc, err
		}
	}
	return doc, nil
}

func (q *Queries) DocumentsSelect(ctx context.Context) (docs []Document, err error) {
	query := `SELECT
							d.id, d.name, d.content, e.embedding
						FROM
							documents d
							INNER JOIN 
								embeddings e ON d.id = e.id`
	result, err := q.conn.QueryOneContext(ctx, query)
	if err != nil {
		return docs, err
	}
	for result.Next() {
		var doc Document
		if err = result.Scan(&doc.ID, &doc.Name, &doc.Content, &doc.Embedding); err != nil {
			return docs, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

type DocumentSelectNearestResult struct {
	Document
	Distance float64 `db:"distance"`
}

func (q *Queries) DocumentsSelectNearest(ctx context.Context) (docs []DocumentSelectNearestResult, err error) {
	query := `WITH vec_results as (
							SELECT
								e.id, e.embedding, distance
							FROM
								embeddings e
							WHERE
								e.embedding match '[0.890, 0.544, 0.825, 0.961, 0.358, 0.0196, 0.521, 0.175]' order by distance limit 2)
						SELECT
							d.id, d.name, d.content, vr.embedding, vr.distance
						FROM vec_results vr
							INNER JOIN
								documents d ON vr.id = d.id;`
	result, err := q.conn.QueryOneContext(ctx, query)
	if err != nil {
		return docs, err
	}
	for result.Next() {
		var doc DocumentSelectNearestResult
		var embeddingJSON string
		if err = result.Scan(&doc.ID, &doc.Name, &doc.Content, &embeddingJSON, &doc.Distance); err != nil {
			return docs, err
		}
		fmt.Printf("Embedding len: %d\n", len(embeddingJSON))
		for i, r := range embeddingJSON {
			fmt.Printf("Rune [%d]: %v %s\n", i, r, string(r))
		}

		docs = append(docs, doc)
	}
	return docs, nil
}
