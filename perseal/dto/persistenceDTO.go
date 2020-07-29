package dto

import (
	"encoding/json"
	"log"

	"github.com/EC-SEAL/perseal/model"
	"github.com/EC-SEAL/perseal/sm"
	"golang.org/x/oauth2"
)

type PersistenceDTO struct {
	//The SessionID
	ID string
	//The Microservice Token
	MSToken string
	//The PDS Location(googleDrive, oneDrive, Browser, Mobile)
	PDS string
	//The operation (Store or Load)
	Method string
	//The URL that the persistence redirects to when finishing its processes
	ClientCallbackAddr string
	//The Custom URL, used in Mobile implementation
	CustomURL string

	//The password to encrypt or decrypt the DataStore
	Password string
	//The byte array of the selected local file, used in Browser Implementation
	LocalFileBytes []byte

	//The OAuth Tokens to access the cloud services
	GoogleAccessCreds oauth2.Token
	OneDriveToken     oauth2.Token

	//The Response that is written in the HTML page
	Response model.HTMLResponse

	//Option to be used in the Persistence Menus
	MenuOption string
	//Current User Device (Desktop or Mobile)
	UserDevice string
	//The QRcode to be shown in the HTML, used in the Mobile implementation
	Image string
}

// Builds Standard Persistence DTO
func PersistenceBuilder(id string, smResp sm.SessionMngrResponse, method ...string) (dto PersistenceDTO, err *model.HTMLResponse) {

	client := smResp.SessionData.SessionVariables[model.EnvVariables.SessionVariables.ClientCallbackAddr]

	log.Println(client)
	if client == "" && model.Test {
		client = model.EnvVariables.TestURLs.MockRedirectDashboard
	}
	pds := smResp.SessionData.SessionVariables[model.EnvVariables.SessionVariables.PDS]
	userDevice := smResp.SessionData.SessionVariables[model.EnvVariables.SessionVariables.OneDriveToken]

	dto = PersistenceDTO{
		ID:                 id,
		PDS:                pds,
		ClientCallbackAddr: client,
		UserDevice:         userDevice,
	}
	googleTokenBytes, oneDriveTokenBytes, err := getGoogleAndOneDriveTokens(dto, smResp)
	if err != nil {
		return
	}

	json.Unmarshal(googleTokenBytes, &dto.GoogleAccessCreds)
	json.Unmarshal(oneDriveTokenBytes, &dto.OneDriveToken)

	if len(method) > 0 || method != nil {
		dto.Method = method[0]
	} else {
		dto.Method = smResp.SessionData.SessionVariables[model.EnvVariables.SessionVariables.CurrentMethod]
	}

	dto.UserDevice = "Mobile"
	return
}

// Builds Persistence DTO With Password
func PersistenceWithPasswordBuilder(id, password string, smResp sm.SessionMngrResponse, method ...string) (dto PersistenceDTO, err *model.HTMLResponse) {

	client := smResp.SessionData.SessionVariables[model.EnvVariables.SessionVariables.ClientCallbackAddr]

	log.Println(client)
	if client == "" && model.Test {
		client = model.EnvVariables.TestURLs.MockRedirectDashboard
	}

	pds := smResp.SessionData.SessionVariables[model.EnvVariables.SessionVariables.PDS]

	dto = PersistenceDTO{
		ID:                 id,
		PDS:                pds,
		ClientCallbackAddr: client,
		Password:           password,
	}
	googleTokenBytes, oneDriveTokenBytes, err := getGoogleAndOneDriveTokens(dto, smResp)
	if err != nil {
		return
	}
	if len(method) > 0 || method != nil {
		dto.Method = method[0]
	} else {
		dto.Method = smResp.SessionData.SessionVariables[model.EnvVariables.SessionVariables.CurrentMethod]
	}

	json.Unmarshal(googleTokenBytes, &dto.GoogleAccessCreds)
	json.Unmarshal(oneDriveTokenBytes, &dto.OneDriveToken)
	return
}

func getGoogleAndOneDriveTokens(dto PersistenceDTO, smResp sm.SessionMngrResponse) (googleTokenBytes, oneDriveTokenBytes []byte, err *model.HTMLResponse) {

	oneDriveToken := smResp.SessionData.SessionVariables[model.EnvVariables.SessionVariables.OneDriveToken]
	var token1 interface{}
	json.Unmarshal([]byte(oneDriveToken), &token1)
	oneDriveTokenBytes, erro := json.Marshal(token1)
	if erro != nil {
		err = &model.HTMLResponse{
			Code:         400,
			Message:      "Could Not Marshall One Drive Token",
			ErrorMessage: erro.Error(),
		}
		return
	}

	googleDrive := smResp.SessionData.SessionVariables[model.EnvVariables.SessionVariables.GoogleDriveToken]
	var token2 interface{}
	json.Unmarshal([]byte(googleDrive), &token2)
	googleTokenBytes, erro = json.Marshal(token2)
	if erro != nil {
		err = &model.HTMLResponse{
			Code:         400,
			Message:      "Could Not Marshall One Drive Token",
			ErrorMessage: erro.Error(),
		}
	}
	return
}
