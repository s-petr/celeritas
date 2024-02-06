package main

func doMigrate(arg2, arg3 string) error {
	dsn := getDSN()

	switch arg2 {
	case "up":
		return cel.MigrateUp(dsn)
	case "down":
		if arg3 == "all" {
			return cel.MigrateDownAll(dsn)
		} else {
			return cel.Steps(-1, dsn)
		}
	case "reset":
		if err := cel.MigrateDownAll(dsn); err != nil {
			return err
		}

		return cel.MigrateUp(dsn)
	default:
		showHelp()
	}
	return nil
}
