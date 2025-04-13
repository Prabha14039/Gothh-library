package helpers

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

