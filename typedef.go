package main

import "github.com/google/uuid"

type AuthKey struct {
	ID  uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Key string    `gorm:"not null"`
	Val string    `gorm:"not null"`
}

type Tree struct {
	ID       uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string     `gorm:"type:varchar(255);not null"`
	ParentID *uuid.UUID `gorm:"type:uuid;index"`
	SHA      string     `gorm:"type:varchar(255);not null"`
	URL      string     `gorm:"type:varchar(255);not null"`
	Parent   *Tree      `gorm:"foreignkey:ParentID"`
}

type Content struct {
	ID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Content string    `gorm:"not null"`
	NodeID  uuid.UUID `gorm:"type:uuid;not null"`
}
