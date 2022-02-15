package todos

import "github.com/kamva/mgm/v3"

type Todo struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	Id               string `form:"id" json:"id" bson:"todoId"`
	Text             string `form:"text" json:"text" binding:"required" bson:"text"`
	Done             string `form:"done" json:"done" bson:"done"`
	Deleted          string `form:"deleted" json:"deleted" bson:"deleted"`
}
