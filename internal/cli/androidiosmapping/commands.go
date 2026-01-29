package androidiosmapping

import "github.com/peterbourgon/ff/v3/ffcli"

// Command returns the android-ios-mapping command group.
func Command() *ffcli.Command {
	return AndroidIosMappingCommand()
}
