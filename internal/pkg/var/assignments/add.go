package assigns

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Add adds an assignment into the collection with specified username and activity name.
// It returns the ID of the new assignment and the channel that receives a result of the assignment.
func Add() (invkID string, onResolved chan interface{}) {
	var (
		err error
		id  uuid.UUID
	)

	for {
		if id, err = uuid.NewRandom(); err != nil {
			err = errors.Wrapf(err, "Failed to create random UUID")
			panic(err)
		}

		invkID = id.String()

		if _, ok := assignments[invkID]; !ok {
			break
		}
	}

	onResolved = make(chan interface{}, 1)
	assignments[invkID] = assign{
		id:     invkID,
		holder: onResolved,
	}

	return
}
