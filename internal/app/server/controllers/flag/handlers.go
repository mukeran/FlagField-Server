package flag

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	. "github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/gin-gonic/gin"
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
		tx := GetDB(c)
		problem := GetProblem(c)
		flags := problem.GetFlags(tx)
		response.Flags = BindList(flags)
	}
}

func (*Handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqCreate
			response RespCreate
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		problem := GetProblem(c)
		response.FlagOrder, err = problem.AddFlagWithSettingsMap(tx, cpt.FlagType(request.Type), request.Settings.(map[string]interface{}))
	}
}

func (*Handler) BatchDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqBatchDelete
			response RespBatchDelete
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		problem := GetProblem(c)
		flags := problem.GetFlags(tx)
		l := len(flags)
		for _, flagOrder := range request.Orders {
			var (
				r    RespBatchDeleteDetail
				err2 error
			)
			if flagOrder > l || flagOrder <= 0 {
				err2 = ErrOutOfRange
			} else {
				(flags)[flagOrder-1].Delete(tx)
			}
			if err2 != nil {
				err = ErrPartlyFailed
			}
			r.OldOrder = flagOrder
			r.Pack(&err2)
			response.Details = append(response.Details, &r)
		}
	}
}

func (*Handler) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespShow
			err      error
		)
		defer Pack(c, &err, &response)
		flag := GetFlag(c)
		response.Flag = BindFlag(flag)
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
		flag := GetFlag(c)
		err = request.Bind(flag)
		if err != nil {
			return
		}
		flag.Update(tx)
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
		flag := GetFlag(c)
		flag.Delete(tx)
	}
}
