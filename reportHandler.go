package consulKvMigrator

import "fmt"

type reportHandler struct {
	chainableDifferenceHandler
	dryRun bool
}

// Reports on changes. If no there are no changes then return 'nil', breaking the handling chain.
func (handler reportHandler) handle(differences differences) error {
	if len(differences) > 0 {
		if handler.dryRun {
			fmt.Print("Running in dry-run mode. ")
		}

		fmt.Printf("Changes found: %v\n", len(differences))
		return handler.next(differences)
	}
	fmt.Println("No changes to migrate.")
	return nil
}
