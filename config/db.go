package config

import (
	"log"
	"os"

	freedb "github.com/FreeLeh/GoFreeDB"
	"github.com/FreeLeh/GoFreeDB/google/auth"
)

var DBAuth *auth.Service

func GetSheetDB(spreadsheetID, sheetName string, columns []string) *freedb.GoogleSheetRowStore {
	if DBAuth == nil {
		err := getAuthConfig()

		if err != nil {
			log.Fatalf("Unable to connect Google Spreadsheet: %v", err)
		}
	}

	db := freedb.NewGoogleSheetRowStore(
		DBAuth,
		spreadsheetID,
		sheetName,
		freedb.GoogleSheetRowStoreConfig{
			Columns: columns,
		},
	)

	return db
}

func getAuthConfig() error {
	serviceAccountPath := os.Getenv("SERVICE_ACCOUNT_PATH")

	var err error
	DBAuth, err = auth.NewServiceFromFile(
		serviceAccountPath,
		freedb.FreeDBGoogleAuthScopes,
		auth.ServiceConfig{},
	)
	if err != nil {
		return err
	}

	return nil
}
