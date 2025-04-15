package helpers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo"
)

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

	log.Printf("Received upload URL: %s", response.Data[0].URL)
	ch <- response
}

func uploadFiles(value []byte, Response UploadThingResponse) {

	Data := Response.Data[0]
	Fields := Data.Fields
	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)

	for key, val := range Fields {
		err := writer.WriteField(key, val)
		if err != nil {
			log.Println("Error adding value : %s with error : %v", key, err)
		}
	}

	filename := filepath.Base(Fields["key"])

	FieldForm, err := writer.CreateFormFile("file", filename)

	FieldForm.Write(value)

	_ = writer.Close()

	req, err := http.NewRequest("POST", Response.Data[0].URL, &buf)
	if err != nil {
		log.Printf("Error creating PUT request: %v", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("File upload failed: %v", err)
		return
	}

	defer res.Body.Close()

	log.Printf("Upload finished. Status: %s", res.Status) // if this prints 204 then taht means it was succesful
}

func updateDatabase(db *sql.DB, name string, url string) {
	query := "INSERT INTO images (name, url) VALUES ($1, $2)"
	_, err := db.Exec(query, name, url)
	if err != nil {
		log.Println("Error inserting image into database: ", err)
		return
	}
}

func Upload(c echo.Context, db *sql.DB) error {

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

	uploadResponse := make(chan UploadThingResponse)

	var uploadBody UploadRequest = UploadRequest{
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

	go prepareUpload(uploadBody, uploadResponse) //Preparing upload

	value := <-uploadResponse

	if len(value.Data) == 0 {
		c.Logger().Error("UploadThing response is empty")
		return echo.NewHTTPError(http.StatusInternalServerError, "upload failed")
	}

	uploadFiles(fileData, value) // uploadingFiles

	updateDatabase(db, fileInfo.Name, value.Data[0].FileUrl)

	log.Printf("Upload completed.")
	return nil
}
