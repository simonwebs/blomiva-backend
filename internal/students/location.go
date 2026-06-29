package students

import "time"

type LocationInfo struct {
	Type        string     `json:"type,omitempty" bson:"type,omitempty"`
	Coordinates []float64  `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
	Latitude    float64    `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude   float64    `json:"longitude,omitempty" bson:"longitude,omitempty"`
	Accuracy    float64    `json:"accuracy,omitempty" bson:"accuracy,omitempty"`
	Source      string     `json:"source,omitempty" bson:"source,omitempty"`
	Permission  string     `json:"permission,omitempty" bson:"permission,omitempty"`
	CapturedAt  *time.Time `json:"capturedAt,omitempty" bson:"capturedAt,omitempty"`
}
