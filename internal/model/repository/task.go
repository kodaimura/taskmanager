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
			PRIORITY_ID INTEGER NOT NULL,
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
    SelectByUId(uid int, state, deadline, priority string) ([]dto.TaskExp1, error)
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


func (tr *taskRepository)SelectByUId(uid int, state, deadline, priority string) ([]dto.TaskExp1, error){
	var tasks []dto.TaskExp1
	rows, err := tr.db.Query(
		`SELECT 
			TID, 
			UID, 
			TASK, 
			MEMO, 
			PERCENT, 
			STATE_ID, 
			GE.VALUE1 AS STATE, 
			DEADLINE, 
			PRIORITY_ID, 
			GE2.VALUE1 AS PRIORITY,
			T.CREATE_AT, 
			T.UPDATE_AT
		 FROM TASKS AS T
		 LEFT JOIN GENERALS AS GE 
		 	ON T.STATE_ID = GE.KEY1 AND GE.CLASS = 'task_state'
		 LEFT JOIN GENERALS AS GE2 
		 	ON T.PRIORITY_ID = GE2.KEY1 AND GE2.CLASS = 'task_priority'
		 WHERE UID = ?
		 AND (? = "" OR STATE = ?)
		 AND (? = "" OR DEADLINE = ?)
		 AND (? = "" OR PRIORITY = ?)`,
		 uid, state, state, deadline, deadline, priority, priority,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := dto.TaskExp1{}
		err = rows.Scan(
			&t.TId, &t.UId, &t.Task, &t.Memo, &t.Percent, 
			&t.StateId, &t.State, &t.Deadline, 
			&t.PriorityId, &t.Priority, &t.CreateAt, &t.UpdateAt,
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
			TID, 
			UID, 
			TASK, 
			MEMO, 
			PERCENT, 
			STATE_ID, 
			GE.VALUE1 AS STATE, 
			DEADLINE, 
			PRIORITY_ID, 
			GE2.VALUE1 AS PRIORITY , 
			T.CREATE_AT, 
			T.UPDATE_AT
		 FROM TASKS AS T
		 LEFT JOIN GENERALS AS GE 
		 	ON T.STATE_ID = GE.KEY1 AND GE.CLASS = 'task_state'
		 LEFT JOIN GENERALS AS GE2 
		 	ON T.PRIORITY_ID = GE2.KEY1 AND GE2.CLASS = 'task_priority'
		 WHERE TID = ?`, tid,
	).Scan(
		&t.TId, &t.UId, &t.Task, &t.Memo, &t.Percent, 
			&t.StateId, &t.State, &t.Deadline, 
			&t.PriorityId, &t.Priority, &t.CreateAt, &t.UpdateAt,
	)

	return t, err
}


func (tr *taskRepository)Insert(tiDto *dto.TaskInsertDto) error {
	_, err := tr.db.Exec(
		`INSERT INTO TASKS (
			UID, TASK, MEMO, PERCENT, STATE_ID, DEADLINE, PRIORITY_ID
		 ) 
		 VALUES(?, ?, ?, ?, ?, ?, ?)`,
		tiDto.UId, tiDto.Task, tiDto.Memo, tiDto.Percent, 
		tiDto.StateId, tiDto.Deadline, tiDto.PriorityId,
	)
	return err
}


func (tr *taskRepository)Update(tuDto *dto.TaskUpdateDto, tid, uid int) error {
	_, err := tr.db.Exec(
		`UPDATE TASKS SET 
			TASK = ?, 
			MEMO = ?, 
			PERCENT = ?, 
			STATE_ID = ?, 
			DEADLINE = ?, 
			PRIORITY_ID = ?
		 WHERE TID = ? AND UID = ?`,
		tuDto.Task, tuDto.Memo, tuDto.Percent, tuDto.StateId, 
		tuDto.Deadline, tuDto.PriorityId,
		tid, uid,
	)
	return err
}

