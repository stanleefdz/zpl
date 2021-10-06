package models

type Player struct {
	ID           int32  `bson:"id"`
	Name         string `bson:"name"`
	Age          int32  `bson:"age"`
	Role         string `bson:"role"`
	Country      string `bson:"country"`
	BattingStyle string `bson:"battingstyle"`
	BowlingStyle string `bson:"bowlingstyle"`
}

type Team struct {
	ID         int32  `bson:"id"`
	Name       string `bson:"name"`
	Owner      string `bson:"owner"`
	HomeGround string `bson:"homeground"`
}

type ErrorMessage struct {
	Description string `bson:"description"`
	Message     string `bson:"message"`
	StatusCode  string `bson:"statuscode"`
}
