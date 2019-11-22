package config

import (
	"fmt"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/bndr/gotabulate"
	"github.com/jinzhu/gorm"
	"github.com/urfave/cli/v2"
)

func list(db *gorm.DB) cli.ActionFunc {
	return func(c *cli.Context) error {
		configs := cpt.GetConfigs(db)
		var output [][]interface{}
		for _, config := range configs {
			output = append(output, []interface{}{config.Key, config.Value})
		}
		if len(output) != 0 {
			t := gotabulate.Create(output)
			t.SetHeaders([]string{"Key", "Value"})
			t.SetAlign("right")
			fmt.Println(t.Render("grid"))
		}
		fmt.Printf("%v rows\n", len(output))
		return nil
	}
}

func set(db *gorm.DB) cli.ActionFunc {
	return func(c *cli.Context) error {
		fromFile := c.String("from-file")
		var output [][]interface{}
		if fromFile != "" {
			fmt.Println("This function is not implemented yet")
			return nil
		} else {
			key := c.String("key")
			value := c.String("value")
			if key == "" {
				return errors.ErrEmptyConfigKey
			}
			tx := db.Begin()
			defer func() {
				if err := recover(); err != nil {
					tx.Rollback()
					panic(err)
				}
			}()
			cpt.SetConfig(tx, key, value)
			if v := tx.Commit(); v.Error != nil {
				panic(v.Error)
			}
			output = append(output, []interface{}{key, value})
		}
		fmt.Print("Success!\n")
		t := gotabulate.Create(output)
		t.SetHeaders([]string{"Key", "Value"})
		t.SetAlign("right")
		fmt.Println(t.Render("grid"))
		return nil
	}
}

func get(db *gorm.DB) cli.ActionFunc {
	return func(c *cli.Context) error {
		key := c.String("key")
		value := cpt.GetConfig(db, key)
		t := gotabulate.Create([][]interface{}{{key, value}})
		t.SetHeaders([]string{"Key", "Value"})
		t.SetAlign("right")
		fmt.Println(t.Render("grid"))
		return nil
	}
}
