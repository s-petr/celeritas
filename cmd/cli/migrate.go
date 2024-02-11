package main

func doMigrate(arg2, arg3 string) error {
	// dsn := getDSN()
	checkForDB()

	tx, err := cel.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer tx.Close()

	switch arg2 {
	case "up":
		// return cel.MigrateUp(dsn)
		return cel.RunPopMigrations(tx)
	case "down":
		if arg3 == "all" {
			return cel.PopMigrateDown(tx, -1)
		} else {
			return cel.PopMigrateDown(tx, 1)
		}
		// if arg3 == "all" {
		// 	return cel.MigrateDownAll(dsn)
		// } else {
		// 	return cel.Steps(-1, dsn)
		// }
	case "reset":
		return cel.PopMigrateReset(tx)
		// if err := cel.MigrateDownAll(dsn); err != nil {
		// 	return err
		// }
		// return cel.MigrateUp(dsn)
	default:
		showHelp()
	}
	return nil
}
