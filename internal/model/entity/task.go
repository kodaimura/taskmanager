package entity


type Task struct {
	TId int `db: "TID" json: "tid"`
	UId int `db: "UID" json: "uid"`
	Task string `db: "TASK" json: "task"`
	Memo string `db: "MEMO" json: "memo"`
	Percent int `db: Percent json: "percent"`
	StateId int `db: "STATE_ID" json: "stateid"`
	Deadline string `db: "DEADLINE" json: "deadline"`
	PriorityId int `db: "PRIORITY_ID" json: "priorityid"`
	CreateAt string `db:"CREATE_AT" json:"createat"`
	UpdateAt string `db:"UPDATE_AT" json:"updateat"`
}