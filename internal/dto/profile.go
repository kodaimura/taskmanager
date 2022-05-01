package dto


type ProfileExp1 struct {
	UId int `db:"UID" json:"uid"`
	UserName string `db:"USER_NAME" json:"username"`
	GId int `db:"GID" json:"gid"`
	GroupName string `db:"GROUP_NAME" json:"groupname"`
}