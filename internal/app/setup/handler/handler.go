package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/fatih/color"

	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
	"github.com/FlagField/FlagField-Server/internal/app/setup/sections"
)

type Handler struct {
	Sections []*sections.Section
	Mapping  *hooks.Map
	AfterAll hooks.HookFunc
}

func (h *Handler) Section(_name, _description string, _before, _after hooks.HookFunc) *sections.Section {
	section := &sections.Section{
		Name:        _name,
		Description: _description,
		Before:      _before,
		After:       _after,
	}
	h.Sections = append(h.Sections, section)
	return section
}

func (h *Handler) SetAfterAll(_hook hooks.HookFunc) {
	h.AfterAll = _hook
}

func (h *Handler) DefaultSection(_name, _description string) *sections.Section {
	return h.Section(_name, _description, hooks.Default(), hooks.Default())
}

func (h *Handler) Proceed() {
begin:
	itemKey := color.New(color.FgWhite, color.BgBlue)
	itemDescription := color.New(color.FgBlack, color.BgYellow)
	currentSection := 0
	for currentSection < len(h.Sections) {
	beginSection:
		currentItem := 0
		section := h.Sections[currentSection]
		if section.Before != nil {
			beforeSectionSign := section.Before(h.Mapping)
			switch beforeSectionSign {
			case hooks.Normal:
			case hooks.BeginSection:
				goto beginSection
			case hooks.EndSection:
				goto endSection
			case hooks.NextSection:
				currentSection++
				continue
			case hooks.Begin:
				goto begin
			case hooks.End:
				goto end
			}
		}
		fmt.Printf("\nSection %s, %s:\n", section.Name, section.Description)
		for currentItem < len(section.Items) {
		beginItem:
			var (
				reader *bufio.Reader
				buf    []byte
				err    error
			)
			item := section.Items[currentItem]
			if item.Before != nil {
				beforeItemSign := item.Before(h.Mapping)
				switch beforeItemSign {
				case hooks.Normal:
				case hooks.BeginSection:
					goto beginSection
				case hooks.EndSection:
					goto endSection
				case hooks.NextSection:
					currentSection++
					if currentSection < len(h.Sections) {
						goto beginSection
					} else {
						goto end
					}
				case hooks.BeginItem:
					goto beginItem
				case hooks.EndItem:
					goto endItem
				case hooks.NextItem:
					currentItem++
					continue
				case hooks.Begin:
					goto begin
				case hooks.End:
					goto end
				}
			}
			_, _ = itemKey.Printf("%s", item.Key)
			fmt.Print("->")
			_, _ = itemDescription.Printf("%s", item.Description)
			fmt.Print(": ")
			color.Unset()
			reader = bufio.NewReader(os.Stdin)
			buf, _, err = reader.ReadLine()
			if err != nil {
				panic(err)
			}
			buf = bytes.TrimSpace(buf)
			switch item.Type {
			case reflect.String:
				(*h.Mapping)[item.Key] = string(buf)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				str := string(buf)
				num, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Invalid number! Please check and reinput.\n")
					goto beginItem
				}
				(*h.Mapping)[item.Key] = num
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				str := string(buf)
				num, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Invalid number! Please check and reinput.\n")
					goto beginItem
				}
				(*h.Mapping)[item.Key] = num
			case reflect.Float32, reflect.Float64:
				str := string(buf)
				num, err := strconv.ParseFloat(str, 64)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Invalid number! Please check and reinput.\n")
					goto beginItem
				}
				(*h.Mapping)[item.Key] = num
			case reflect.Bool:
				str := string(buf)
				b, err := strconv.ParseBool(str)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Invalid boolean! Please check and reinput.\n")
					goto beginItem
				}
				(*h.Mapping)[item.Key] = b
			}
		endItem:
			if item.After != nil {
				afterItemSign := item.After(h.Mapping)
				switch afterItemSign {
				case hooks.Normal:
					currentItem++
				case hooks.BeginSection:
					goto beginSection
				case hooks.EndSection:
					goto endSection
				case hooks.NextSection:
					currentSection++
					if currentSection < len(h.Sections) {
						goto beginSection
					} else {
						goto end
					}
				case hooks.BeginItem:
					goto beginItem
				case hooks.EndItem:
					goto endItem
				case hooks.NextItem:
					currentItem++
					continue
				case hooks.Begin:
					goto begin
				case hooks.End:
					goto end
				}
			}
		}
	endSection:
		if section.After != nil {
			afterSectionSign := section.After(h.Mapping)
			switch afterSectionSign {
			case hooks.Normal:
				currentSection++
			case hooks.BeginSection:
				goto beginSection
			case hooks.EndSection:
				goto endSection
			case hooks.NextSection:
				currentSection++
				continue
			case hooks.Begin:
				goto begin
			case hooks.End:
				goto end
			}
		}
	}
end:
	if h.AfterAll != nil {
		afterAllSign := h.AfterAll(h.Mapping)
		switch afterAllSign {
		case hooks.Normal:
		case hooks.Begin:
			goto begin
		case hooks.End:
			goto end
		}
	}
}
