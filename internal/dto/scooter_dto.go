package dto

type CreateScooterDTO struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type UpdateScooterDTO struct {
	Identifier             string
	Latitude               *float32 `json:"latitude"`
	Longitude              *float32 `json:"longitude"`
	OccupiedUserIdentifier *string  `json:"occupied_user_identifier" validate:"scooter_dto:check_occupied"`
}
