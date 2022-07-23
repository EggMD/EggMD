// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var (
	// ErrTeamNotFound is returned when a team is not found.
	ErrTeamNotFound = errors.New("team not found")
	// ErrTeamNameAlreadyExists is returned when a team name already exists.
	ErrTeamNameAlreadyExists = errors.New("team name already exists")
)

// TeamsStore is the persistent interface for teams.
type TeamsStore interface {
	// Create creates a new team.
	Create(ctx context.Context, opts CreateTeamOptions) (*Team, error)
	// Update updates the team with given id.
	Update(ctx context.Context, id uint, opts UpdateTeamOptions) error
	// GetByID returns the team with given id.
	GetByID(ctx context.Context, id uint) (*Team, error)
	// GetByUID returns the team with given uid.
	GetByUID(ctx context.Context, uid string) (*Team, error)
	// GetByName returns the team with given name.
	GetByName(ctx context.Context, name string) (*Team, error)
	// DeleteByUID deletes the team with given uid.
	DeleteByUID(ctx context.Context, uid string) error
}

func NewTeamsStore(db *gorm.DB) TeamsStore {
	return &teams{db}
}

var Teams TeamsStore

var _ TeamsStore = (*teams)(nil)

type teams struct {
	*gorm.DB
}

// Team represents the team.
type Team struct {
	gorm.Model

	UID     string `gorm:"UNIQUE"`
	Name    string `gorm:"UNIQUE"`
	OwnerID uint   `gorm:"NOT NULL"`
	Owner   *User  `gorm:"-"`
}

type CreateTeamOptions struct {
	Name    string
	OwnerID uint
}

func (db *teams) Create(ctx context.Context, opts CreateTeamOptions) (*Team, error) {
	_, err := db.GetByName(ctx, opts.Name)
	if err != nil {
		if !errors.Is(err, ErrTeamNotFound) {
			return nil, errors.Wrap(err, "get team by name")
		}
	} else {
		return nil, ErrTeamNameAlreadyExists
	}

	uid := uuid.NewV4().String()
	team := &Team{
		UID:     uid,
		Name:    opts.Name,
		OwnerID: opts.OwnerID,
	}
	if err := db.WithContext(ctx).Create(team).Error; err != nil {
		return nil, errors.Wrap(err, "create team")
	}

	teams, err := db.loadAttributes(ctx, team)
	if err != nil {
		return nil, errors.Wrap(err, "load attributes")
	}
	if len(teams) == 0 {
		return nil, errors.New("teams is empty after load attributes")
	}
	return teams[0], nil
}

type UpdateTeamOptions struct {
	Name string
}

func (db *teams) Update(ctx context.Context, id uint, opts UpdateTeamOptions) error {
	_, err := db.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get team by ID")
	}

	t, err := db.GetByName(ctx, opts.Name)
	if err == nil {
		if t.ID != id {
			return ErrTeamNameAlreadyExists
		}
	} else if err != ErrTeamNotFound {
		return errors.Wrap(err, "get team by team name")
	}

	if err := db.WithContext(ctx).Model(&Team{}).
		Where("id = ?", id).
		Updates(&Team{
			Name: opts.Name,
		}).Error; err != nil {
		return errors.Wrap(err, "update team")
	}
	return nil
}

func (db *teams) getBy(ctx context.Context, whereQuery interface{}, whereArgs ...interface{}) (*Team, error) {
	var team Team
	if err := db.WithContext(ctx).Model(&Team{}).Where(whereQuery, whereArgs...).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}

	teams, err := db.loadAttributes(ctx, &team)
	if err != nil {
		return nil, errors.Wrap(err, "load attributes")
	}
	if len(teams) == 0 {
		return nil, errors.New("teams is empty after load attributes")
	}
	return teams[0], nil
}

func (db *teams) GetByID(ctx context.Context, id uint) (*Team, error) {
	return db.getBy(ctx, "id = ?", id)
}

func (db *teams) GetByUID(ctx context.Context, uid string) (*Team, error) {
	return db.getBy(ctx, "uid = ?", uid)
}

func (db *teams) GetByName(ctx context.Context, name string) (*Team, error) {
	return db.getBy(ctx, "name = ?", name)
}

func (db *teams) DeleteByUID(ctx context.Context, uid string) error {
	team, err := db.GetByUID(ctx, uid)
	if err != nil {
		return errors.Wrap(err, "get team by uid")
	}

	if err := db.WithContext(ctx).Model(&Team{}).Delete(&Team{}, "id = ?", team.ID).Error; err != nil {
		return errors.Wrap(err, "delete team")
	}
	return nil
}

func (db *teams) loadAttributes(ctx context.Context, teams ...*Team) ([]*Team, error) {
	if len(teams) == 0 {
		return teams, nil
	}

	userIDSet := make(map[uint]struct{}, len(teams))
	for _, team := range teams {
		userIDSet[team.OwnerID] = struct{}{}
	}
	userIDs := make([]uint, 0, len(userIDSet))
	for userID := range userIDSet {
		userIDs = append(userIDs, userID)
	}

	usersStore := NewUsersStore(db.DB)
	users, err := usersStore.GetByIDs(ctx, userIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "get users by IDs")
	}

	userSet := make(map[uint]*User)
	for _, user := range users {
		user := user
		userSet[user.ID] = user
	}

	for _, team := range teams {
		team.Owner = userSet[team.OwnerID]
	}
	return teams, nil
}
