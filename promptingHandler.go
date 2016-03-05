package consulKvMigrator

import (
	"bufio"
	"fmt"
	"github.com/apcera/termtables"
	"os"
)

type promptingHandler struct {
	chainableDifferenceHandler
	dryRun bool
}

func (handler promptingHandler) handle(differences differences) error {
	table := termtables.CreateTable()

	table.AddHeaders("Key", "Old Value", "New Value")
	for idx := range differences {
		diff := differences[idx]
		table.AddRow(diff.key, diff.targetValue, diff.sourceValue)
	}
	fmt.Println(table.Render())

	if handler.dryRun {
		return handler.next(differences)
	}

	fmt.Print("Are you sure you want to apply these changes? [Y/n]")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadByte()
	if err != nil {
		return err
	}
	fmt.Print("\n")

	switch input {
	case 'Y', 'y', '\n':
		return handler.next(differences)
	}
	fmt.Println("Aborting migration.")
	return nil
}
