package controllers

import (
	"github.com/mukeran/email"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	validatorV8 "gopkg.in/go-playground/validator.v8"

	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/FlagField/FlagField-Server/internal/pkg/validator"
)

type ResponseInterface interface {
	GetHttpStatus() int
	Pack(*error)
}

type Response struct {
	Status     int                  `json:"status"`
	Message    string               `json:"message"`
	Validate   *map[string][]string `json:"validate,omitempty"`
	httpStatus int                  `json:"-"`
}

/* Universal responses definition */
var (
	RespSuccess                 = Response{0, "success", nil, http.StatusOK}
	RespBadRequest              = Response{1, "bad request", nil, http.StatusBadRequest}
	RespNotFound                = Response{2, "not found", nil, http.StatusNotFound}
	RespLoggedIn                = Response{10, "have been logged in", nil, http.StatusBadRequest}
	RespNotLoggedIn             = Response{11, "haven't been logged in", nil, http.StatusBadRequest}
	RespDuplicatedName          = Response{12, "duplicated name", nil, http.StatusBadRequest}
	RespDuplicatedUsername      = Response{100, "duplicated username", nil, http.StatusBadRequest}
	RespDuplicatedEmail         = Response{101, "duplicated email", nil, http.StatusBadRequest}
	RespValidationFailed        = Response{102, "validate failed", nil, http.StatusBadRequest}
	RespWrongUsernameOrPassword = Response{201, "wrong username or password", nil, http.StatusNotFound}
	RespDuplicatedAlias         = Response{300, "duplicated alias", nil, http.StatusBadRequest}
	RespNotEnoughPoints         = Response{400, "not enough points", nil, http.StatusBadRequest}
	RespMemberParticipated      = Response{500, "some of the members you selected have participated the contest in another team", nil, http.StatusBadRequest}
	RespTeamParticipated        = Response{501, "your team has participated the contest", nil, http.StatusBadRequest}
	RespNotMemberOfTeam         = Response{502, "some of the members you selected are not a part of your team", nil, http.StatusBadRequest}
	RespNotAdminOfTeam          = Response{503, "the user is not the team's admin", nil, http.StatusBadRequest}
	RespProcessed               = Response{600, "the resource has been processed", nil, http.StatusBadRequest}
	RespTimeout                 = Response{601, "the resource has expired", nil, http.StatusBadRequest}
	RespAlreadyJoinedInTeam     = Response{602, "the user has joined the team", nil, http.StatusBadRequest}
	RespInvalidInvitationToken  = Response{603, "invalid invitation token", nil, http.StatusBadRequest}
	RespNotInvited              = Response{604, "the user has not been invited", nil, http.StatusBadRequest}
	RespAlreadyInvited          = Response{605, "the user has been invited", nil, http.StatusBadRequest}
	RespNoSuchApplication       = Response{606, "no such application", nil, http.StatusBadRequest}
	RespAlreadyApplied          = Response{607, "the user has applied", nil, http.StatusBadRequest}
	RespContestPending          = Response{700, "the contest is pending", nil, http.StatusBadRequest}
	RespContestStarted          = Response{701, "the contest is started", nil, http.StatusBadRequest}
	RespContestEnded            = Response{702, "the contest is ended", nil, http.StatusBadRequest}
	RespInvalidCaptchaFor       = Response{800, "invalid captcha for", nil, http.StatusBadRequest}
	RespServerError             = Response{-1, "internal server error", nil, http.StatusInternalServerError}
	RespSessionExpired          = Response{-2, "session is expired", nil, http.StatusUnauthorized}
	RespPermissionDenied        = Response{-3, "permission denied", nil, http.StatusUnauthorized}
	RespPartlyFailed            = Response{-4, "request partly failed", nil, http.StatusMultiStatus}
	RespWriteConfigFile         = Response{-5, "failed to write config file", nil, http.StatusInternalServerError}
	RespCoolDown                = Response{-6, "cooling down", nil, http.StatusBadRequest}
	RespInvalidCaptcha          = Response{-7, "invalid captcha", nil, http.StatusBadRequest}
	RespNotInWhitelist          = Response{-8, "not in whitelist", nil, http.StatusForbidden}
	RespDangerousOperation      = Response{-9, "dangerous operation", nil, http.StatusForbidden}
	RespFailedToSendEmail       = Response{-10, "failed to send email", nil, http.StatusInternalServerError}
)

/* Mapping errors to response */
var errorMappings = map[error]*Response{
	nil:                              &RespSuccess,
	errors.ErrDuplicatedUsername:     &RespDuplicatedUsername,
	errors.ErrDuplicatedEmail:        &RespDuplicatedEmail,
	errors.ErrDuplicatedAlias:        &RespDuplicatedAlias,
	errors.ErrDatabaseConnection:     &RespServerError,
	errors.ErrNotFound:               &RespNotFound,
	errors.ErrUserNotFound:           &RespWrongUsernameOrPassword,
	errors.ErrWrongPassword:          &RespWrongUsernameOrPassword,
	errors.ErrInvalidRequest:         &RespBadRequest,
	errors.ErrInvalidFlagType:        &RespBadRequest,
	errors.ErrInvalidFlagSettings:    &RespBadRequest,
	errors.ErrSessionExpired:         &RespSessionExpired,
	errors.ErrRedisConnection:        &RespServerError,
	errors.ErrLoggedIn:               &RespLoggedIn,
	errors.ErrNotLoggedIn:            &RespNotLoggedIn,
	errors.ErrPermissionDenied:       &RespPermissionDenied,
	errors.ErrOutOfRange:             &RespBadRequest,
	errors.ErrPartlyFailed:           &RespPartlyFailed,
	errors.ErrNotEnoughPoints:        &RespNotEnoughPoints,
	errors.ErrWriteConfigFile:        &RespWriteConfigFile,
	errors.ErrMemberParticipated:     &RespMemberParticipated,
	errors.ErrTeamParticipated:       &RespTeamParticipated,
	errors.ErrNotMemberOfTeam:        &RespNotMemberOfTeam,
	errors.ErrNotAdminOfTeam:         &RespNotAdminOfTeam,
	errors.ErrDuplicatedName:         &RespDuplicatedName,
	errors.ErrProcessed:              &RespProcessed,
	errors.ErrTimeout:                &RespTimeout,
	errors.ErrAlreadyJoinedInTeam:    &RespAlreadyJoinedInTeam,
	errors.ErrInvalidInvitationToken: &RespInvalidInvitationToken,
	errors.ErrNotInvited:             &RespNotInvited,
	errors.ErrAlreadyInvited:         &RespAlreadyInvited,
	errors.ErrNoSuchApplication:      &RespNoSuchApplication,
	errors.ErrAlreadyApplied:         &RespAlreadyApplied,
	errors.ErrContestPending:         &RespContestPending,
	errors.ErrContestStarted:         &RespContestStarted,
	errors.ErrContestEnded:           &RespContestEnded,
	errors.ErrInvalidCaptchaFor:      &RespInvalidCaptchaFor,
	errors.ErrCoolDown:               &RespCoolDown,
	errors.ErrInvalidCaptcha:         &RespInvalidCaptcha,
	errors.ErrNotInWhitelist:         &RespNotInWhitelist,
	errors.ErrDangerousOperation:     &RespDangerousOperation,
	gorm.ErrRecordNotFound:           &RespNotFound,
	gorm.ErrCantStartTransaction:     &RespServerError,
	gorm.ErrInvalidSQL:               &RespServerError,
	gorm.ErrInvalidTransaction:       &RespServerError,
	gorm.ErrUnaddressable:            &RespServerError,
	os.ErrNotExist:                   &RespNotFound,
	email.ErrClosed:                  &RespFailedToSendEmail,
	email.ErrTimeout:                 &RespFailedToSendEmail,
}

func (resp *Response) GetHttpStatus() int {
	return resp.httpStatus
}

func (resp *Response) Pack(err *error) {
	switch (*err).(type) {
	case validatorV8.ValidationErrors:
		resp.Status = RespValidationFailed.Status
		resp.Message = RespValidationFailed.Message
		resp.Validate = validator.GenerateDetail(*err)
		resp.httpStatus = RespValidationFailed.httpStatus
	default:
		_resp, ok := errorMappings[*err]
		if ok {
			resp.Status = _resp.Status
			resp.Message = _resp.Message
			resp.httpStatus = _resp.httpStatus
		} else {
			resp.Status = RespServerError.Status
			resp.Message = RespServerError.Message
			resp.httpStatus = RespServerError.httpStatus
		}
	}
}

// func desensitize(obj interface{}, perms *map[string]bool, parameter *map[string]string) {
// 	t := reflect.TypeOf(obj)
// 	v := reflect.ValueOf(obj)
// 	for t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 	}
// 	for v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	switch v.Kind() {
// 	case reflect.Struct:
// 		for i := 0; i < v.NumField(); i++ {
// 			requiredPerms := strings.Split(t.Field(i).Tag.Get("desensitize"), ";")
// 			re := regexp.MustCompile(`\(\w+\)`)
// 			keep := false
// 			empty := true
// 			for _, rp := range requiredPerms {
// 				rp = strings.TrimSpace(rp)
// 				if rp == "" {
// 					continue
// 				}
// 				empty = false
// 				keys := re.FindAllString(rp, -1)
// 				for _, key := range keys {
// 					_, ok := (*parameter)[key]
// 					if !ok {
// 						panic("Invalid desensitize parameter: no such key")
// 					}
// 					rp = strings.Replace(rp, key, (*parameter)[key], -1)
// 				}
// 				_, ok := (*perms)[rp]
// 				if ok {
// 					keep = true
// 					break
// 				}
// 			}
// 			if !empty && !keep {
// 				v.Field(i).Set(reflect.New(t.Field(i).Type).Elem())
// 			} else {
// 				desensitize(v.Field(i).Addr().Interface(), perms, parameter)
// 			}
// 		}
// 	case reflect.Slice, reflect.Array:
// 		for i := 0; i < v.Len(); i++ {
// 			desensitize(v.Index(i).Addr().Interface(), perms, parameter)
// 		}
// 	}
// }
//
// func Desensitize(obj interface{}, perms *map[string]bool) {
// 	t := reflect.TypeOf(obj)
// 	v := reflect.ValueOf(obj)
// 	for t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 	}
// 	for v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	p := map[string]string{}
// 	switch v.Kind() {
// 	case reflect.Struct:
// 		for i := 0; i < v.NumField(); i++ {
// 			switch v.Field(i).Kind() {
// 			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
// 				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
// 				reflect.String:
// 				p[fmt.Sprintf("(%s)", t.Field(i).Name)] = fmt.Sprintf("%v", v.Field(i).Interface())
// 			}
// 		}
// 		desensitize(obj, perms, &p)
// 	case reflect.Slice, reflect.Array:
// 		for i := 0; i < v.Len(); i++ {
// 			Desensitize(v.Index(i).Addr().Interface(), perms)
// 		}
// 	}
// }
