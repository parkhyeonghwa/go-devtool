package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type MysqlBackup struct {
	Options MysqlCommonOptions `group:"common"`
	Positional struct {
		Schema string `description:"Schema" required:"1"`
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *MysqlBackup) Execute(args []string) error {
	conf.Options.Init()

	defer NewSigIntHandler(func() {
	})()

	conf.Options.dumpCompression = GetCompressionByFilename(conf.Positional.Filename)
	if (conf.Options.dumpCompression != "") {
		fmt.Println(fmt.Sprintf(" - Using %s compression", conf.Options.dumpCompression))
	}

	cmd := shell.Cmd(conf.Options.MysqlDumpCommandBuilder(conf.Positional.Schema)...).Pipe(fmt.Sprintf("cat > %s", shell.Quote(conf.Positional.Filename)))
	cmd.Run()

	return nil
}
