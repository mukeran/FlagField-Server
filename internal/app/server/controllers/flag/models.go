package flag

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
)

type ReqCreate struct {
	Type     string      `json:"type" validate:"required" binding:"required"`
	Settings interface{} `json:"settings"`
}

type ReqModify struct {
	Type     *string     `json:"type" validate:"omitempty,min=1" binding:"omitempty,min=1"`
	Settings interface{} `json:"settings"`
}

func (r *ReqModify) Bind(f *cpt.Flag) (err error) {
	if r.Type != nil {
		f.Type = *r.Type
	}
	if r.Settings != nil {
		f.Settings, err = cpt.NewFlagSettingsFromMap(cpt.FlagType(f.Type), r.Settings.(map[string]interface{}))
	}
	return
}

type ReqBatchDelete struct {
	Orders []int `json:"orders"`
}

type RespCreate struct {
	Response
	FlagOrder int `json:"flag_order,omitempty"`
}

type RFlag struct {
	Type     string      `json:"type"`
	Settings interface{} `json:"settings"`
}

func BindFlag(flag *cpt.Flag) *RFlag {
	return &RFlag{
		Type:     flag.Type,
		Settings: flag.Settings,
	}
}

type RespList struct {
	Response
	Flags []*RFlag `json:"flags"`
}

func BindList(flags []cpt.Flag) []*RFlag {
	var out []*RFlag
	for _, flag := range flags {
		out = append(out, BindFlag(&flag))
	}
	return out
}

type RespShow struct {
	Response
	Flag *RFlag `json:"flag"`
}

type RespBatchDeleteDetail struct {
	Response
	OldOrder int `json:"old_order"`
}

type RespBatchDelete struct {
	Response
	Details []*RespBatchDeleteDetail `json:"details"`
}
