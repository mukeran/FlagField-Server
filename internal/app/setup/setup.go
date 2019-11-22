package setup

import (
	"fmt"
	"os"

	"github.com/FlagField/FlagField-Server/internal/app/setup/handler"
	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
	"github.com/FlagField/FlagField-Server/internal/app/setup/sections"
	"github.com/FlagField/FlagField-Server/internal/app/setup/sections/database"
	"github.com/FlagField/FlagField-Server/internal/app/setup/sections/final"
	"github.com/FlagField/FlagField-Server/internal/app/setup/sections/redis"
	"github.com/FlagField/FlagField-Server/internal/app/setup/sections/resource"
	"github.com/FlagField/FlagField-Server/internal/app/setup/sections/server"
)

const (
	motd = `   ___ _               ___ _      _     _       __ by mukeran, am009 & Tinywangxx
  / __\ | __ _  __ _  / __(_) ___| | __| |     / _\ ___ _ ____   _____ _ __ 
 / _\ | |/ _` + "`" + ` |/ _` + "`" + ` |/ _\ | |/ _ \ |/ _` + "`" + ` |_____\ \ / _ \ '__\ \ / / _ \ '__|
/ /   | | (_| | (_| / /   | |  __/ | (_| |_____|\ \  __/ |   \ V /  __/ |   
\/    |_|\__,_|\__, \/    |_|\___|_|\__,_|     \__/\___|_|    \_/ \___|_|   
               |___/
________    _____      by mukeran
__  ___/______  /____  _________ 
_____ \_  _ \  __/  / / /__  __ \
____/ //  __/ /_ / /_/ /__  /_/ /
/____/ \___/\__/ \__,_/ _  .___/ 
                        /_/
Welcome to the FlagField-Server setup!
This setup is for version %v, please make sure you have the right server version installed.
`
	destVersion = "1.0.0"
)

type Setup struct {
	Handler handler.Handler
	Mapping hooks.Map
}

func hasSetupLock() bool {
	_, err := os.Stat("setup-lock")
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func showMotd() {
	fmt.Printf(motd, destVersion)
}

func (s *Setup) initSections() {
	server.Register(&s.Handler)
	database.Register(&s.Handler)
	redis.Register(&s.Handler)
	resource.Register(&s.Handler)
	final.Register(&s.Handler)
	s.Handler.SetAfterAll(sections.AfterAll())
}

func (s *Setup) Run() {
	if hasSetupLock() {
		fmt.Printf("Detected \"setup-lock\" file, which means that the server may have been configured.\n")
		fmt.Printf("If you want to continue, please delete \"setup-lock\" first.\n")
		os.Exit(1)
	}
	showMotd()
	s.Handler.Mapping = &s.Mapping
	s.initSections()
	s.Handler.Proceed()
}

func New() *Setup {
	s := &Setup{}
	s.Mapping = hooks.Map{}
	return s
}
