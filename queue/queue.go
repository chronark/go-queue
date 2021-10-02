package queue

import (
	"context"
	"fmt"

	"github.com/chronark/go-queue/fauna"
	f "github.com/fauna/faunadb-go/v4/faunadb"
)

const (
	TODO        = "todo"
	IN_PROGRESS = "progress"
	DONE        = "done"
)

type Queue struct {
	faunaClient *f.FaunaClient
	ctx         context.Context
}

func NewQueue(faunaClient *f.FaunaClient) (*Queue, error) {

	return &Queue{
		faunaClient: faunaClient,
		ctx:         context.Background(),
	}, nil
}

func (q *Queue) Produce(message Message) error {

	res, err := q.faunaClient.Query(f.Create(f.Collection(fauna.COLLECTION_TODO), f.Obj{"data": f.ToObject(message)}))
	if err != nil {
		fmt.Printf("%+v\n", err)
		return err
	}
	fmt.Printf("%+v\n", res)

	return nil
}

func (q *Queue) Consume(topic string) (message Message, err error) {

	res, err := q.faunaClient.Query(
		f.Let().Bind(
			"message",
			f.Get(
				f.MatchTerm(
					f.Index(fauna.INDEX_TODO_BY_TOPIC),
					topic,
				),
			),
		).In(

			f.Do(
				f.Delete(f.Select("ref", f.Var("message"))),
				f.Create(fauna.COLLECTION_IN_PROGRESS, f.Obj{"data": f.Select("data", f.Var("message"))}),
			),
		),
	)

	if err != nil {
		return Message{}, err
	}

	err = res.At(f.ObjKey("data")).Get(&message)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}

func (q *Queue) Acknowledge(messageId string) error {
	_, err := q.faunaClient.Query(
		f.Let().Bind(
			"message",
			f.Get(
				f.MatchTerm(
					f.Index(fauna.INDEX_IN_PROGRESS_BY_ID),
					messageId,
				),
			),
		).In(

			f.Do(
				f.Delete(f.Select("ref", f.Var("message"))),
				f.Create(fauna.COLLECTION_DONE, f.Obj{"data": f.Select("data", f.Var("message"))}),
			),
		),
	)

	return err
}
