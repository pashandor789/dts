package types

type TestType int

const (
	Input  TestType = 1
	Output TestType = 2
)

type Test struct {
	Id     uint16   `json:"id" bson:"id"`
	Type   TestType `json:"type" bson:"type"`
	TaskId uint64   `json:"task_id" bson:"task_id"`
	Data   []byte   `json:"data" bson:"data"`
}

type Build struct {
	Id            uint64 `json:"id" bson:"_id"`
	ExecuteScript []byte `json:"execute_script" bson:"execute_script"`
	InitScript    []byte `json:"init_script" bson:"init_script"`
}

type TaskMeta struct {
	TaskId uint64
}
