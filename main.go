package main

func main() {

	// auth, err := auth.NewServiceFromFile(
	// 	"comment-system-service-account.json",
	// 	freedb.FreeDBGoogleAuthScopes,
	// 	auth.ServiceConfig{},
	// )
	// if err != nil {
	// 	log.Fatalf("Unable to connect Google Spreadsheet: %v", err)
	// }

	// spreadsheetId := "1gyZHny4-3H4T9PokiKhUB-_zad-T8Zqw1TWKODzH-ms"
	// sheetName := "post-1"
	// store := freedb.NewGoogleSheetRowStore(
	// 	auth,
	// 	spreadsheetId,
	// 	sheetName,
	// 	freedb.GoogleSheetRowStoreConfig{Columns: comment.Columns},
	// )

	// defer store.Close(context.Background())

	StartServer()
}
