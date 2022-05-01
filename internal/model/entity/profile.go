package entity


type Profile struct {
	UId int `db:"UID" json:"uid"`
	GId int `db "GID" json: "gid"`
	CreateAt string `db:"CREATE_AT" json:"createat"`
	UpdateAt string `db:"UPDATE_AT" json:"updateat"`
}