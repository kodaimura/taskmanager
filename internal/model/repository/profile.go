package repository

import (
	"log"
	"database/sql"

	"taskmanager/internal/dto"
	"taskmanager/internal/model/db"
	"taskmanager/internal/model/entity"
)


func init(){
	db := db.GetDB()

	cmd := `
		CREATE TABLE IF NOT EXISTS PROFILES (
			UID INTEGER NOT NULL UNIQUE,
			GID INTEGER,
			CREATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
			UPDATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
		);

		CREATE TRIGGER IF NOT EXISTS trigger_profiles_updated_at AFTER UPDATE ON PROFILES
		BEGIN
    		UPDATE PROFILES 
    			SET UPDATE_AT = DATETIME('now', 'localtime') 
    			WHERE rowid == NEW.rowid;
		END;`

	_, err := db.Exec(cmd)

	if err != nil {
		log.Panic(err)
	}
}


type ProfileRepository interface {
	SelectByUId(uid int) (entity.Profile, error)
    Upsert(profile *entity.Profile) error
    UpdateByUId(uid int, profile *entity.Profile) error
    GetProfileExp1ByUId(uid int) (dto.ProfileExp1, error)
}


type profileRepository struct {
	db *sql.DB
}


func NewProfileRepository() ProfileRepository {
	db := db.GetDB()
	return &profileRepository{db}
}


func (pr *profileRepository) SelectByUId(uid int) (entity.Profile, error){
	var p entity.Profile
	err := pr.db.QueryRow(
		`SELECT UID, GID, CREATE_AT, UPDATE_AT
		 FROM PROFILES
		 WHERE UID = ?`, uid,
	).Scan(
		&p.UId, &p.GId, &p.CreateAt, &p.UpdateAt,
	)

	return p, err
}


func (pr *profileRepository) Upsert(profile *entity.Profile) error {
	_, err := pr.db.Exec(
		`INSERT OR REPLACE INTO PROFILES (UID, GID)
		 VALUES(?,?)`, profile.UId, profile.GId,
	)
	return err
}


func (pr *profileRepository) UpdateByUId(uid int, profile *entity.Profile) error {
	_, err := pr.db.Exec(
		`UPDATE USERS SET GID = ? WHERE UID = ?`,
		profile.GId, uid,
	)
	return err
}


func (pr *profileRepository) GetProfileExp1ByUId(uid int) (dto.ProfileExp1, error){
	var pe dto.ProfileExp1
	err := pr.db.QueryRow(
		`SELECT 
			P.UID, U.USER_NAME, P.GID, G.GROUP_NAME
		 FROM PROFILES AS P
		 LEFT JOIN GROUPS AS G ON P.GID = G.GID 
		 JOIN USERS AS U ON P.UID = U.UID
		 WHERE P.UID = ?`, uid,
	).Scan(
		&pe.UId, &pe.UserName, &pe.GId, &pe.GroupName,
	)

	return pe, err
}
