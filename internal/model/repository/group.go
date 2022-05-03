package repository

import (
	"log"
	"database/sql"

	"taskmanager/internal/model/db"
	"taskmanager/internal/model/entity"
)


func init(){
	db := db.GetDB()

	cmd := `
		CREATE TABLE IF NOT EXISTS GROUPS (
			GID INTEGER PRIMARY KEY AUTOINCREMENT,
			GROUP_NAME VARCHAR(100) NOT NULL UNIQUE,
			PASSWORD VARCHAR(100) NOT NULL,
			CREATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
			UPDATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
		);

		CREATE TRIGGER IF NOT EXISTS trigger_group_updated_at AFTER UPDATE ON GROUPS
		BEGIN
    		UPDATE GROUPS 
    			SET UPDATE_AT = DATETIME('now', 'localtime') 
    			WHERE rowid == NEW.rowid;
		END;`

	_, err := db.Exec(cmd)

	if err != nil {
		log.Panic(err)
	}
}


type GroupRepository interface {
    SelectByGId(gid int) (entity.Group, error)
    SelectByGroupName(groupName string) (entity.Group, error)
    Insert(groupName, password string) error
}


type groupRepository struct {
	db *sql.DB
}


func NewGroupRepository() GroupRepository {
	db := db.GetDB()
	return &groupRepository{db}
}


func (gr *groupRepository)SelectByGId(gid int) (entity.Group, error){
	var g entity.Group
	err := gr.db.QueryRow(
		`SELECT 
		 	GID, 
		 	GROUP_NAME, 
		 	CREATE_AT, 
		 	UPDATE_AT 
		 FROM GROUPS 
		 WHERE GID = ?`, gid,
	).Scan(
		&g.GId, &g.GroupName, &g.CreateAt, &g.UpdateAt,
	)

	return g, err
}


func (gr *groupRepository)SelectByGroupName(groupName string) (entity.Group, error){
	var g entity.Group
	err := gr.db.QueryRow(
		`SELECT 
			GID, 
			GROUP_NAME, 
			PASSWORD, 
			CREATE_AT, 
			UPDATE_AT 
		 FROM GROUPS 
		 WHERE GROUP_NAME = ?`, groupName,
	).Scan(
		&g.GId, &g.GroupName, &g.Password, &g.CreateAt, &g.UpdateAt,
	)

	return g, err
}


func (gr *groupRepository)Insert(groupName, password string) error {
	_, err := gr.db.Exec(
		`INSERT INTO GROUPS (
			GROUP_NAME, PASSWORD
		 ) VALUES(?, ?)`,
		 groupName, password,
	)
	return err
}

