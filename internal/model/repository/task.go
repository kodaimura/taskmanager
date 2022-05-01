package repository

import (
	"log"
	"database/sql"

	"taskmanager/internal/model/db"
	"taskmanager/internal/dto"
)


func init(){
	db := db.GetDB()

	cmd := `
		CREATE TABLE IF NOT EXISTS TASKS (
			TID INTEGER PRIMARY KEY AUTOINCREMENT,
			UID INTEGER NOT NULL,
			TASK VARCHAR(100) NOT NULL,
			MEMO TEXT NOT NULL,
			PERCENT INTEGER NOT NULL,
			STATE_ID INTEGER NOT NULL,
			DEADLINE VARCHAR(10),
			CREATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
			UPDATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
		);

		CREATE TRIGGER IF NOT EXISTS trigger_task_updated_at AFTER UPDATE ON TASKS
		BEGIN
    		UPDATE TASKS
    			SET UPDATE_AT = DATETIME('now', 'localtime') 
    			WHERE rowid == NEW.rowid;
		END;`

	_, err := db.Exec(cmd)

	if err != nil {
		log.Panic(err)
	}
}


type TaskRepository interface {
    SelectByUId(uid int, state, deadline string) ([]dto.TaskExp1, error)
    SelectByTId(tid int) (dto.TaskExp1, error)
    Insert(tuDto *dto.TaskInsertDto) error
    Update(tuDto *dto.TaskUpdateDto, tid, uid int) error
}


type taskRepository struct {
	db *sql.DB
}


func NewTaskRepository() TaskRepository {
	db := db.GetDB()
	return &taskRepository{db}
}


func (tr *taskRepository)SelectByUId(uid int, state, deadline string) ([]dto.TaskExp1, error){
	var tasks []dto.TaskExp1
	rows, err := tr.db.Query(
		`SELECT 
			TID, UID, TASK, MEMO, PERCENT, STATE_ID,
			GE.VALUE1 AS STATE, DEADLINE, T.CREATE_AT, T.UPDATE_AT
		 FROM TASKS AS T
		 LEFT JOIN GENERALS AS GE 
		 ON T.STATE_ID = GE.KEY1 AND GE.CLASS = 'task_state'
		 WHERE UID = ?
		 AND (? = "" OR STATE = ?)
		 AND (? = "" OR DEADLINE = ?)`, uid, state, state, deadline, deadline,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := dto.TaskExp1{}
		err = rows.Scan(
			&t.TId, &t.UId, &t.Task, &t.Memo, &t.Percent, &t.StateId,
			&t.State, &t.Deadline, &t.CreateAt, &t.UpdateAt,
		)
		if err != nil {
			break
		}
		tasks = append(tasks, t)
	}

	return tasks, err
}


func (tr *taskRepository)SelectByTId(tid int) (dto.TaskExp1, error){
	var t dto.TaskExp1
	err := tr.db.QueryRow(
		`SELECT
			TID, UID, TASK, MEMO, PERCENT, STATE_ID
			STATE, DEADLINE, T.CREATE_AT, T.UPDATE_AT
		 FROM TASKS AS T
		 RIGHT JOIN GENERALS AS GE 
		 ON GE.KEY1 = T.STATE_ID AND GE.CLASS = 'task_state'
		 WHERE TID = ?`, tid,
	).Scan(
		&t.TId, &t.UId, &t.Task, &t.Memo, &t.Percent, &t.StateId,
		&t.State, &t.Deadline, &t.CreateAt, &t.UpdateAt,
	)

	return t, err
}


func (tr *taskRepository)Insert(tiDto *dto.TaskInsertDto) error {
	_, err := tr.db.Exec(
		`INSERT INTO TASKS (UID, TASK, MEMO, PERCENT, STATE_ID, DEADLINE) 
		 VALUES(?, ?, ?, ?, ?, ?)`,
		tiDto.UId, tiDto.Task, tiDto.Memo, 
		tiDto.Percent, tiDto.StateId, tiDto.Deadline,
	)
	return err
}


func (tr *taskRepository)Update(tuDto *dto.TaskUpdateDto, tid, uid int) error {
	_, err := tr.db.Exec(
		`UPDATE TASKS SET 
			TASK = ?, MEMO = ?, PERCENT = ?, STATE_ID = ?, DEADLINE = ?
		 WHERE TID = ? AND UID = ?`,
		tuDto.Task, tuDto.Memo, tuDto.Percent, tuDto.StateId, tuDto.Deadline,
		tid, uid,
	)
	return err
}

