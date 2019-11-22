package final

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"

	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
	"github.com/FlagField/FlagField-Server/internal/app/setup/sections/database"
	"github.com/FlagField/FlagField-Server/internal/app/setup/sections/redis"
	"github.com/FlagField/FlagField-Server/internal/app/setup/sections/resource"
	"github.com/FlagField/FlagField-Server/internal/app/setup/sections/server"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/config"
)

type Hooks struct {
}

func (*Hooks) Before() hooks.HookFunc {
	return func(m *hooks.Map) uint {
		fmt.Println()
		conf := config.DefaultConfig
		conf.Server.Port = uint16((*m)[server.KeyPort].(uint64))
		conf.Database.Type = (*m)[database.KeyType].(string)
		conf.Database.Parameter = (*m)[database.KeyParameter].(string)
		conf.Redis.URI = (*m)[redis.KeyURI].(string)
		conf.Redis.Password = (*m)[redis.KeyPassword].(string)
		conf.Redis.Db = int((*m)[redis.KeyDB].(uint64))
		conf.Resource.BaseDir = (*m)[resource.KeyBaseDir].(string)
		color.Yellow("Saving configure to config.json...\n")
		err := conf.ToFile("./config.json")
		if err != nil {
			red := color.Set(color.FgHiRed)
			_, _ = red.Fprintf(os.Stderr, "Failed to create config.json! The setup will end.\n")
			os.Exit(2)
		}
		color.Yellow("Running migrator...\n")
		cmd := exec.Command("./dist/migrator", "-template", "initial")
		var out, serr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &serr
		err = cmd.Run()
		if err != nil {
			red := color.Set(color.FgHiRed)
			_, _ = red.Fprintf(os.Stderr, "Failed to run migrator.\nStdout: (%s)\nStderr: (%s)\nError message: (%v)\nThe setup will end.\n", out.String(), serr.String(), err)
			os.Exit(3)
		}
		return hooks.Normal
	}
}

func (*Hooks) After() hooks.HookFunc {
	return func(m *hooks.Map) uint {
		color.Yellow("Creating super user...\n")
		db, _ := gorm.Open((*m)[database.KeyType].(string), (*m)[database.KeyParameter].(string))
		defer db.Close()
		db.SingularTable(true)
		if v := db.Where("name = ?", cpt.DefaultContest.Name).First(&cpt.DefaultContest); v.Error != nil {
			panic(v.Error)
		}
		if v := db.Where("username = ?", cpt.DefaultUser.Username).First(&cpt.DefaultUser); v.Error != nil {
			panic(v.Error)
		}
		username := (*m)[KeyUsername].(string)
		password := (*m)[KeyPassword].(string)
		email := (*m)[KeyEmail].(string)
		tx := db.Begin()
		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
				red := color.Set(color.FgHiRed)
				_, _ = red.Fprintf(os.Stderr, "Operation failed!\nError message: (%v)\nThe setup will end.\n", err)
				os.Exit(4)
			}
		}()
		user, err := cpt.NewUser(tx, username, password)
		if err != nil {
			tx.Rollback()
			red := color.Set(color.FgHiRed)
			_, _ = red.Fprintf(os.Stderr, "Failed to create admin user.\nError message: (%v)\nThe setup will end.\n", err)
			os.Exit(4)
		}
		user.Email = email
		user.IsAdmin = true
		user.Update(tx)
		tx.Commit()
		color.Green("Successfully created super user!\n")
		color.Yellow("Set configs...\n")
		tx = db.Begin()
		cpt.SetConfig(tx, "system.setup_time", time.Now().Format(time.RFC3339))
		cpt.SetConfig(tx, "index.notification", "FlagField-Server has been configured! Explore this fantasy software!")
		tx.Commit()
		color.Green("Successfully set configs!\n")
		return hooks.Normal
	}
}
