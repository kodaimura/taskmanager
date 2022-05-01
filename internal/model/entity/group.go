package entity


type Group struct {
	GId int `db: "GID" json: "gid"`
	GroupName string `db: "GROUP_NAME" json: "groupname"`
	Password string `db: "PASSWORD" json: "password"`
	CreateAt string `db:"CREATE_AT" json:"createat"`
	UpdateAt string `db:"UPDATE_AT" json:"updateat"`
}