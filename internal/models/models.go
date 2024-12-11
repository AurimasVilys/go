package models

import "time"

type User struct {
	Identifier string `db:"identifier,pk"`
}

type Scooter struct {
	Identifier             string  `db:"identifier,pk"`
	OccupiedUserIdentifier *string `db:"occupied_user_identifier" json:"occupied_user_identifier"`
	LastConfirmedLatitude  float32 `db:"last_confirmed_latitude" json:"last_confirmed_latitude"`
	LastConfirmedLongitude float32 `db:"last_confirmed_longitude" json:"last_confirmed_longitude"`
}

type ScooterEvent struct {
	Identifier        string    `db:"identifier,pk"`
	ScooterIdentifier string    `db:"scooter_identifier"`
	Event             string    `db:"event"`
	Timestamp         time.Time `db:"timestamp"`
	Latitude          float32   `db:"latitude"`
	Longitude         float32   `db:"longitude"`
}
