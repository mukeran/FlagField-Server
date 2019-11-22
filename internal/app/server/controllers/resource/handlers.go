package resource

import (
	"fmt"
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
)

type Handler struct {
}

func (*Handler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespList
			err      error
		)
		defer Pack(c, &err, &response)
		_contestID, hasC := c.GetQuery("contest_id")
		problemAlias, hasP := c.GetQuery("problem_alias")
		if hasP && !hasC {
			err = errors.ErrInvalidRequest
			return
		}
		tx := GetDB(c)
		session := GetSession(c)
		var resources []cpt.Resource
		if hasC {
			var contestID uint64
			if _contestID == "__practice__" {
				contestID = uint64(cpt.DefaultContest.ID)
			} else {
				contestID, err = strconv.ParseUint(_contestID, 10, 32)
				if err != nil {
					err = errors.ErrInvalidRequest
					return
				}
			}
			if !session.User.IsAdmin && !cpt.IsContestAdmin(tx, uint(contestID), session.UserID) {
				err = errors.ErrNotFound
				return
			}
			if !hasP {
				resources = cpt.GetContestResources(tx, uint(contestID))
			} else {
				problemID, err := cpt.GetProblemID(tx, uint(contestID), problemAlias)
				if err != nil {
					return
				}
				resources = cpt.GetProblemResources(tx, problemID)
			}
		} else {
			if !session.User.IsAdmin {
				err = errors.ErrNotFound
				return
			}
			resources = cpt.GetResources(tx)
		}
		response.Resources = BindList(resources)
	}
}

func (*Handler) Upload() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespUpload
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		session := GetSession(c)
		config := GetConfig(c)
		_contestID, hasC := c.GetPostForm("contest_id")
		problemAlias, hasP := c.GetPostForm("problem_alias")
		if hasP && !hasC {
			err = errors.ErrInvalidRequest
			return
		}
		var contestID uint64
		if hasC {
			if _contestID == "__practice__" {
				contestID = uint64(cpt.DefaultContest.ID)
			} else {
				contestID, err = strconv.ParseUint(_contestID, 10, 32)
				if err != nil {
					err = errors.ErrInvalidRequest
					return
				}
			}
			if !session.User.IsAdmin && !cpt.IsContestAdmin(tx, uint(contestID), session.UserID) {
				err = errors.ErrNotFound
				return
			}
			if hasP && !cpt.HasProblemAlias(tx, uint(contestID), problemAlias) {
				err = errors.ErrNotFound
				return
			}
		} else if !session.User.IsAdmin {
			err = errors.ErrNotFound
			return
		}
		fh, err := c.FormFile("file")
		if err != nil {
			err = errors.ErrInvalidRequest
			return
		}
		f, err := fh.Open()
		if err != nil {
			return
		}
		defer f.Close()
		res := cpt.NewResourceWithReader(tx, config.Resource.BaseDir, fh.Filename, fh.Header.Get("content-type"), f, session.UserID)
		if hasC {
			res.ContestID = uint(contestID)
			if hasP {
				problemID, _ := cpt.GetProblemID(tx, uint(contestID), problemAlias)
				res.ProblemID = problemID
			}
			res.Update(tx)
		}
		response.UUID = res.UUID
	}
}

func (*Handler) Download() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("type", "file")
		inline := c.DefaultQuery("inline", "0")
		res := GetResource(c)
		_, err := os.Stat(res.Path)
		if err != nil {
			panic(err)
		}
		disposition := "attachment"
		if inline == "1" || inline == "true" {
			disposition = "inline"
		}
		c.Header("Cache-Control", "no-store")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", fmt.Sprintf("%s; filename=%s", disposition, res.Name))
		c.Header("Content-Type", res.ContentType)
		c.File(res.Path)
	}
}

func (*Handler) Modify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqModify
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		session := GetSession(c)
		res := GetResource(c)
		err = request.Bind(tx, session, res)
		if err != nil {
			return
		}
		res.Update(tx)
	}
}

func (*Handler) SingleDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		res := GetResource(c)
		res.Delete(tx)
	}
}
