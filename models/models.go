package models

type Player struct {
	ID           int64  `bson:"_id"`
	Name         string `bson:"name"`
	Age          int32  `bson:"age"`
	Role         string `bson:"role"`
	Country      string `bson:"country"`
	BattingStyle string `bson:"batting_style"`
	BowlingStyle string `bson:"bowling_style"`
}

type Team struct {
	ID         int64  `bson:"_id"`
	Name       string `bson:"name"`
	Owner      string `bson:"owner"`
	HomeGround string `bson:"homeground"`
}

type ErrorMessage struct {
	Description string `json:"description"`
	Message     string `json:"message"`
}

type SuccessMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
