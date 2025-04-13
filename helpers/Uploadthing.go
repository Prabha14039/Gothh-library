package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"github.com/labstack/echo"
)

type UploadThingResponse struct {
	Data []struct {
		Key               string            `json:"key"`
		FileName          string            `json:"fileName"`
		FileType          string            `json:"fileType"`
		FileUrl           string            `json:"fileUrl"`
		ContentDisposition string           `json:"contentDisposition"`
		PollingJwt        string            `json:"pollingJwt"`
		PollingUrl        string            `json:"pollingUrl"`
		CustomId          string            `json:"customId"`
		URL               string            `json:"url"`
		Fields			  map[string]string `json:"fields"`
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


func prepareUpload(request UploadRequest, ch chan UploadThingResponse) {
	key := FetchEnv()
	url := "https://api.uploadthing.com/v6/uploadFiles"

	payload, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Add("X-Uploadthing-Api-Key", key.UploadThing_Key)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	var response UploadThingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return
	}

	log.Printf(string(body))

	log.Printf("Received upload URL: %s", response.Data[0].URL)
	ch <- response
}

func uploadFiles(value []byte, url string, content_type string) {

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(value))
	if err != nil {
		log.Printf("Error creating PUT request: %v", err)
		return
	}

	req.Header.Set("Content-Type",content_type)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("File upload failed: %v", err)
		return
	}
	defer res.Body.Close()

	log.Printf("File uploaded successfully. Status: %s", res.Status)
}


func Upload(c echo.Context) error {
	var fileInfo File
	file, err := c.FormFile("file")
	if err != nil {
		c.Logger().Error("File not found in form data: ", err)
		return err
	}

	src, err := file.Open()
	if err != nil {
		c.Logger().Error("Failed to open file: ", err)
		return err
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		c.Logger().Error("Failed to read file: ", err)
		return err
	}

	fileInfo.Name = file.Filename
	fileInfo.Size = int(file.Size)
	fileInfo.Type = file.Header.Get("Content-Type")

	uploadBody := UploadRequest{
		Fileinfo: []File{
			{
				Name:     fileInfo.Name,
				Size:     fileInfo.Size,
				CustomId: nil,
				Type:     fileInfo.Type,
			},
		},
		Acl:                "public-read",
		Metadata:           nil,
		ContentDisposition: "inline",
	}

	uploadURL := make(chan UploadThingResponse)
	go prepareUpload(uploadBody, uploadURL)

	value := <-uploadURL
	if len(value.Data) == 0 {
		c.Logger().Error("UploadThing response is empty")
		return echo.NewHTTPError(http.StatusInternalServerError, "upload failed")
	}

	url := value.Data[0].URL
	log.Printf("Uploading file to %s", url)

	content_type := value.Data[0].Fields["Content-Type"]

	log.Printf("%s", content_type)
	uploadFiles(fileData, url, content_type)

	log.Printf("Upload completed.")
	return nil
}

