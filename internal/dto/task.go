package dto


type TaskInsertDto struct {
	UId int `json:"uid"`
	Task string `json:"task"`
	Memo string `djson:"memo"`
	Percent int `json:"percent"`
	StateId int `json:"stateid"`
	Deadline string `json:"deadline"`
	PriorityId int `json: "priorityid"`
}


type TaskUpdateDto struct {
	Task string `json:"task"`
	Memo string `json:"memo"`
	Percent int `json:"percent"`
	StateId int `json:"stateid"`
	Deadline string `json:"deadline"`
	PriorityId int `json: "priorityid"`
}


type TaskExp1 struct {
	TId int `json:"tid"`
	UId int `json:"uid"`
	Task string `json:"task"`
	Memo string `json:"memo"`
	Percent int `json:"percent"`
	StateId int `json:"stateid"`
	State string `json:"state"`
	Deadline string `json:"deadline"`
	PriorityId int `json: "priorityid"`
	Priority string `json: "priority"`
	CreateAt string `json:"createat"`
	UpdateAt string `json:"updateat"`
}