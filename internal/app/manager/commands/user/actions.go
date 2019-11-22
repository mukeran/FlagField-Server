package user

import (
	"fmt"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/constants"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/FlagField/FlagField-Server/internal/pkg/random"
	"github.com/FlagField/FlagField-Server/internal/pkg/validator"
	"github.com/bndr/gotabulate"
	"github.com/jinzhu/gorm"
	"github.com/urfave/cli/v2"
	"strconv"
)

func list(db *gorm.DB) cli.ActionFunc {
	return func(c *cli.Context) error {
		offset := c.Uint("offset")
		limit := c.Uint("limit")
		query := c.String("query")
		users := cpt.GetUsersWithQuery(db, offset, limit, query)
		var output [][]interface{}
		for _, user := range users {
			output = append(output, []interface{}{strconv.FormatUint(uint64(user.ID), 10), user.Username, user.Email, user.CreatedAt.String(), user.IsAdmin})
		}
		if len(output) != 0 {
			t := gotabulate.Create(output)
			t.SetHeaders([]string{"ID", "Username", "Email", "Register Time", "Is Admin"})
			t.SetAlign("right")
			fmt.Println(t.Render("grid"))
		}
		fmt.Printf("%v rows\n", len(output))
		return nil
	}
}

func add(db *gorm.DB) cli.ActionFunc {
	return func(c *cli.Context) error {
		fromFile := c.String("from-file")
		var output [][]interface{}
		if fromFile != "" {
			fmt.Println("This function is not implemented yet")
			return nil
		} else {
			type Register struct {
				Username string `json:"username" validate:"max=20,min=1,keyword"`
				Password string `json:"password" validate:"max=50,min=6"`
				Email    string `json:"email" validate:"email"`
			}
			reg := Register{
				Username: c.String("username"),
				Password: c.String("password"),
				Email:    c.String("email"),
			}
			if reg.Password == "" {
				reg.Password = random.RandomString(12, constants.DicLetterNumeric)
			}
			err := validator.Validate(&reg)
			if err != nil {
				fmt.Print(validator.GenerateDetail(err))
				return err
			}
			tx := db.Begin()
			defer func() {
				if err := recover(); err != nil {
					tx.Rollback()
					panic(err)
				}
			}()
			if cpt.HasEmail(tx, reg.Email) {
				tx.Rollback()
				return errors.ErrDuplicatedEmail
			}
			user, err := cpt.NewUser(tx, reg.Username, reg.Password)
			if err != nil {
				tx.Rollback()
				return err
			}
			user.Email = reg.Email
			user.IsAdmin = c.Bool("admin")
			user.Update(tx)
			if v := tx.Commit(); v.Error != nil {
				panic(v.Error)
			}
			output = append(output, []interface{}{strconv.FormatUint(uint64(user.ID), 10), user.Username, reg.Password, user.Email, user.IsAdmin})
		}
		fmt.Print("Success!\n")
		t := gotabulate.Create(output)
		t.SetHeaders([]string{"ID", "Username", "Password", "Email", "Is Admin"})
		t.SetAlign("right")
		fmt.Println(t.Render("grid"))
		return nil
	}
}

func del(db *gorm.DB) cli.ActionFunc {
	return func(c *cli.Context) error {
		fmt.Print("This function is not implemented yet")
		return nil
	}
}
