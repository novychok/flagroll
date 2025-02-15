package entity

import "time"

type FeatureFlagID string

func (id FeatureFlagID) String() string {
	return string(id)
}

type FeatureFlag struct {
	ID          FeatureFlagID
	OwnerID     UserID
	Name        string
	Active      bool
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FeatureFlagCreate struct {
	OwnerID     UserID
	Active      bool
	Name        string
	Description string
}

type FeatureFlagUpdate struct {
	Active      bool
	Name        string
	Description string
}

type FeatureFlagToggleRequest struct {
	Active bool
}
