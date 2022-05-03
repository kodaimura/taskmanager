package controller

import (
    "strconv"
    "github.com/gin-gonic/gin"
    
    "taskmanager/internal/pkg/jwtauth"
    "taskmanager/internal/model/repository"
    "taskmanager/internal/constants"
    "taskmanager/internal/dto"
)


type IndexController interface {
    IndexPage(c *gin.Context)
    MemberPage(c *gin.Context)
    Task(c *gin.Context)
    UpdateTask(c *gin.Context)
}


type indexController struct {
    ur repository.UserRepository
    tr repository.TaskRepository
    pr repository.ProfileRepository
    ger repository.GeneralRepository
}


func NewIndexController() IndexController {
    ur := repository.NewUserRepository()
    tr := repository.NewTaskRepository()
    pr := repository.NewProfileRepository()
    ger := repository.NewGeneralRepository()
    return indexController{ur, tr, pr, ger}
}


//GET /
func (ic indexController) IndexPage(c *gin.Context) {
    uid, err := jwtauth.ExtractUId(c)
    username, err := jwtauth.ExtractUserName(c)
    gid, err := jwtauth.ExtractGId(c)
    groupname, err := jwtauth.ExtractGroupName(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    state := c.Query("state")
    deadline := c.Query("deadline")
    priority := c.Query("priority")

    members, err := ic.ur.SelectByGId(gid)
    tasks, err := ic.tr.SelectByUId(uid, state, deadline, priority)
    status, _ := ic.ger.SelectByClass("task_state")
    priorities, _ := ic.ger.SelectByClass("task_priority")

    c.HTML(200, "index.html", gin.H{
        "appname": constants.AppName,
        "username": username,
        "groupname": groupname,
        "members": members,
        "tasks": tasks,
        "status": status,
        "priorities": priorities,
    })
}


//GET /members/:uid
func (ic indexController) MemberPage(c *gin.Context) {
    loginUId, err := jwtauth.ExtractUId(c)
    gid, err := jwtauth.ExtractGId(c)
    groupname, err := jwtauth.ExtractGroupName(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    uid, err := strconv.Atoi(c.Param("uid"))
    if err != nil {
        c.Redirect(303, "/logout")
        return
    }

    if uid == loginUId {
        c.Redirect(303, "/")
        return
    }

    pe, err := ic.pr.GetProfileExp1ByUId(uid)
    if err != nil || pe.GId != gid {
        c.Redirect(303, "/logout")
        return
    }

    state := c.Query("state")
    deadline := c.Query("deadline")
    priority := c.Query("priority")

    members, _ := ic.ur.SelectByGId(gid)
    tasks, _ := ic.tr.SelectByUId(uid, state, deadline, priority)
    status, _ := ic.ger.SelectByClass("task_state")
    priorities, _ := ic.ger.SelectByClass("task_priority")

    c.HTML(200, "tasks.html", gin.H{
        "appname": constants.AppName,
        "username": pe.UserName,
        "groupname": groupname,
        "members": members,
        "tasks": tasks,
        "status": status,
        "priorities": priorities,
    })
}


//POST /task
func (ic indexController) Task(c *gin.Context) {
    uid, err := jwtauth.ExtractUId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    td := &dto.TaskInsertDto{}
    td.UId = uid
    td.Task = c.PostForm("task")
    td.Memo = c.PostForm("memo")
    td.Deadline = c.PostForm("deadline")

    per, err := strconv.Atoi(c.PostForm("percent"))
    if err != nil {
        per = 0
    }

    stid, err := strconv.Atoi(c.PostForm("stateid"))
    if err != nil {
        stid = 1
    }

    prid, err := strconv.Atoi(c.PostForm("priorityid"))
    if err != nil {
        prid = 1
    }
    td.Percent = per
    td.StateId = stid
    td.PriorityId = prid

    err = ic.tr.Insert(td)

    if err != nil {
        c.Redirect(303, "/logout")
        return
    }

    c.Redirect(303, "/")
}


//POST /task/:tid
func (ic indexController) UpdateTask(c *gin.Context) {
    uid, err := jwtauth.ExtractUId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    tid, err := strconv.Atoi(c.Param("tid"))
    if err != nil {
        c.Redirect(303, "/")
        return
    }

    td := &dto.TaskUpdateDto{}
    td.Task = c.PostForm("task")
    td.Memo = c.PostForm("memo")
    td.Deadline = c.PostForm("deadline")

    per, err := strconv.Atoi(c.PostForm("percent"))
    if err != nil {
        per = 0
    }

    stid, err := strconv.Atoi(c.PostForm("stateid"))
    if err != nil {
        stid = 1
    }

    prid, err := strconv.Atoi(c.PostForm("priorityid"))
    if err != nil {
        prid = 1
    }

    td.Percent = per
    td.StateId = stid
    td.PriorityId = prid

    err = ic.tr.Update(td, tid, uid)

    if err != nil {
        c.Redirect(303, "/logout")
        return
    }

    c.Redirect(303, "/")
}
