package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/labstack/echo"
)

type PresignedResponse struct {
	Data []struct {
		Key string					`json:"key"`
		Fields map[string]string	`json:"fileds"`
		FileUrl string				`json:"fileUrl"`
		URL string					`json:"url"`
	} `json:"data"`
}

type UploadRequest struct {
	Fileinfo []File				`json:"files"`
	Acl string 					`json:"acl"`
	Metadata *string				`json:"metadata"`
	ContentDisposition string	`json:"contentDisposition"`
}

type File struct {
	Name string			`json:"name"`
	Size int			`json:"size"`
	CustomId *string
	Type string			`json:"type"`
}


func uploadFiles(request UploadRequest) error{

	key := FetchEnv()

	url := "https://api.uploadthing.com/v6/uploadFiles"

	payload, err := json.Marshal(request)
	if err!=nil {
		return err
	}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err!= nil {
		return err
	}

	req.Header.Add("X-Uploadthing-Api-Key", key.UploadThing_Key)
	req.Header.Set("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	fmt.Println("Status:", res.Status)
	fmt.Println("Response Body:", string(body))

	return nil

}

func PrepareUpload(c echo.Context) error {


	var Fileinfo File
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	Fileinfo.Name = file.Filename
	Fileinfo.Size = int(file.Size)
	Fileinfo.Type = file.Header.Get("Content-Type")

	uploadbody := UploadRequest {
		Fileinfo: []File {
			{
			Name : Fileinfo.Name,
			Size : Fileinfo.Size,
			CustomId : nil,
			Type : Fileinfo.Type,
			},
		},
		Acl : "public-read",
		Metadata: nil,
		ContentDisposition: "inline",
	}

	return uploadFiles(uploadbody)
}
