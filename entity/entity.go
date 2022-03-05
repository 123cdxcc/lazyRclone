package entity

type RcloneCopyConfig struct {
	RemoteName        string
	ThreadCount       int
	LogFilePath       string
	UploadingFilePath string
	RemoteFilePath    string
}

type EmailResult struct {
	Status  bool
	Code    int
	Message string
}

type PushResult struct {
	Code  int
	Msg   string
	Data  string
	Count interface{}
}
