package session

import (
	"cobackend/internal/shared"
	"cobackend/internal/db"

	"context"
)

// service.go

func StartPracticeSessionService(
	ctx context.Context,
	userID string,
	input StartPracticeSessionInput,
) (*PracticeSessionResponse, error) {

	// ----------------------------------------------------------
	// Validate Input
	// ----------------------------------------------------------

	if input.AcademyBuildingLaneID <= 0 {
		return nil, shared.ErrInvalidLaneID
	}

	if input.ShootingEventID <= 0 {
		return nil, shared.ErrInvalidShootingEventID
	}

	// ----------------------------------------------------------
	// Begin Transaction
	// ----------------------------------------------------------

	tx, err := db.DB.Begin(
		ctx,
	)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(
		ctx,
	)

	// ----------------------------------------------------------
	// Create Practice Session
	// ----------------------------------------------------------

	session, err := StartPracticeSessionRepository(
		ctx,
		tx,
		userID,
		input,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Commit Transaction
	// ----------------------------------------------------------

	err = tx.Commit(
		ctx,
	)

	if err != nil {
		return nil, err
	}

	return session, nil
}