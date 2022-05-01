package controller

import (
    "github.com/gin-gonic/gin"
    
    "taskmanager/internal/pkg/jwtauth"
    "taskmanager/internal/constants"
    "taskmanager/internal/model/repository"
    "taskmanager/internal/model/entity"
)


type GroupController interface {
    GroupPage(c *gin.Context)
    Group(c *gin.Context)
    BelongToGroup(c *gin.Context)
}


type groupController struct {
    gr repository.GroupRepository
}


func NewGroupController() GroupController {
    gr := repository.NewGroupRepository()
    return groupController{gr}
}


//GET /group
func (gc groupController) GroupPage(c *gin.Context) {
    c.HTML(200, "group.html", gin.H{
        "appname": constants.AppName,
    })
}


//POST /group
func (gc groupController) Group(c *gin.Context) {
    uid, err := jwtauth.ExtractUId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    groupName := c.PostForm("groupname")
    password := c.PostForm("password")

    if _, err := gc.gr.SelectByGroupName(groupName); err == nil {
        c.HTML(409, "group.html", gin.H{
            "appname": constants.AppName,
            "error": "Groupnameが既に使われています。",
        })
        c.Abort()
        return
    }

    if gc.gr.Insert(groupName, password) != nil {
        c.HTML(500, "group.html", gin.H{
            "appname": constants.AppName,
            "error": "登録に失敗しました。",
        })
        c.Abort()
        return
    }

    g, err := gc.gr.SelectByGroupName(groupName);

    if  err != nil {
        c.HTML(500, "group.html", gin.H{
            "appname": constants.AppName,
            "error": "予期せぬエラーが発生しました。",
        })
        c.Abort()
        return
    }

    belongToGroup(uid, g.GId)
    c.Redirect(303, "/logout")
}


//POST /belong
func (gc groupController) BelongToGroup(c *gin.Context) {
    uid, err := jwtauth.ExtractUId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    groupName := c.PostForm("groupname")
    password := c.PostForm("password")

    group, err := gc.gr.SelectByGroupName(groupName)

    if err != nil || password != group.Password {
        c.HTML(401, "group.html", gin.H{
            "appname":constants.AppName,
            "error": "GroupnameまたはPasswordが異なります。",
        })
        c.Abort()
        return
    }

    belongToGroup(uid, group.GId)
    c.Redirect(303, "/logout")
}


func belongToGroup(uid, gid int) error {
    p := &entity.Profile{}
    p.UId = uid
    p.GId = gid
    pr := repository.NewProfileRepository()
    return pr.Upsert(p)
}