package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"testing"
	"time"

	"github.com/EC-SEAL/perseal/dto"
	"github.com/EC-SEAL/perseal/sm"
	"github.com/EC-SEAL/perseal/utils"
)

var (
	id string
)

func InitIntegration(platform string) dto.PersistenceDTO {
	tokenResp, _ := utils.StartSession()
	id = tokenResp.Payload
	smResp, _ := utils.GenerateTokenAPI(platform, id)
	msToken := smResp.Payload

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(msToken)
	http.PostForm("http://localhost:8082/per/store", url.Values{"msToken": {msToken}})

	//simulate google login redirect
	sm.UpdateSessionData(id, "store", "CurrentMethod")
	variables := map[string]string{
		"PDS":       platform,
		"dataStore": "{\"id\":\"DS_3a342b23-8b46-44ec-bb06-a03042135a5e\",\"encryptedData\":null,\"signature\":\"this is the signature\",\"signatureAlgorithm\":\"this is the signature algorithm\",\"encryptionAlgorithm\":\"this is the encryption algorithm\",\"clearData\":null}",
	}

	session := sm.SessionMngrResponse{}
	sessionData := session.SessionData
	sessionData.SessionID = id
	sessionData.SessionVariables = variables
	session.SessionData = sessionData

	obj, _ := dto.PersistenceBuilder(id, session)

	url, _ := GetRedirectURL(obj)
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()

	//wait to select account and store token in session
	time.Sleep(15 * time.Second)
	return obj
}

func TestGoogleService(t *testing.T) {

	obj := InitIntegration("googleDrive")

	// Test Incorrect GoogleDrive Store
	obj.Password = "qwerty"

	ds, err := StoreCloudData(obj, "datastore.seal")
	log.Println(ds)
	log.Println(err)
	if err == nil {
		t.Error("Should have thrown error")
	}

	// Test Correct GoogleDrive Store
	sessionData, _ := sm.GetSessionData(id, "")
	log.Println("sessionData", sessionData)
	obj.GoogleAccessCreds = sessionData.SessionData.SessionVariables["GoogleDriveAccessCreds"]
	obj.SMResp.SessionData.SessionVariables["dataStore"] = "{\"id\":\"DS_3a342b23-8b46-44ec-bb06-a03042135a5e\",\"encryptedData\":null,\"signature\":\"this is the signature\",\"signatureAlgorithm\":\"this is the signature algorithm\",\"encryptionAlgorithm\":\"this is the encryption algorithm\",\"clearData\":null}"
	ds, err = StoreCloudData(obj, "datastore.seal")
	log.Println(ds)
	if err != nil {
		t.Error("Thrown error, got: ", err)
	}

	// Test Correct Load GoogleDrive Store
	ds, err = FetchCloudDataStore(obj, "datastore.seal")
	if err != nil {
		t.Error("Thrown error, got: ", err)
	}
	log.Println(ds)

	// Test Incorrect Load GoogleDrive Store
	ds, err = FetchCloudDataStore(obj, "datastorewrong.seal")
	if err == nil {
		t.Error("Should have thrown error")
	}

	// Test Get Cloud Files
	files, err := GetCloudFileNames(obj)
	if err != nil {
		t.Error("Thrown error, got: ", err)
	}
	if len(files) == 0 {
		t.Error("no files found")
	}

	// Test Get Cloud Files No GoogleCreds
	obj.GoogleAccessCreds = ""
	files, err = GetCloudFileNames(obj)
	if err == nil {
		t.Error("Should have thrown error")
	}

}
