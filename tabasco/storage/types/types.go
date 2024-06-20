package types

type TestType int

const (
	input  TestType = 1
	output TestType = 2
)

type Test struct {
	Id     uint16   `json:"id" bson:"id"`
	Type   TestType `json:"type" bson:"type"`
	TaskId uint64   `json:"task_id" bson:"task_id"`
	Data   string   `json:"data" bson:"data"`
}

type Build struct {
	Id            uint64 `json:"id" bson:"_id"`
	ExecuteScript string `json:"execute_script" bson:"execute_script"`
	InitScript    string `json:"init_script" bson:"init_script"`
}
