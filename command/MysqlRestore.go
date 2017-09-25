package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type MysqlRestore struct {
	Options MysqlCommonOptions `group:"common"`
	Positional struct {
		Schema string `description:"Schema" required:"1"`
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *MysqlRestore) Execute(args []string) error {
	conf.Options.Init()

	defer NewSigIntHandler(func() {
	})()

	conf.Options.ExecMySqlStatement(fmt.Sprintf("DROP DATABASE IF EXISTS %s", mysqlIdentifier(conf.Positional.Schema)))
	conf.Options.ExecMySqlStatement(fmt.Sprintf("CREATE DATABASE %s", mysqlIdentifier(conf.Positional.Schema)))
	cmd := shell.Cmd(fmt.Sprintf("cat %s", shell.Quote(conf.Positional.Filename))).Pipe(conf.Options.MysqlCommandBuilder(conf.Positional.Schema)...)
	cmd.Run()

	return nil
}