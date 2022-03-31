package db

import (
	"container/list"
	"qbCopyProject/entity"
)

var Email string

var Push string

var UploadingTasks = list.New()

var TaskQueue = make(chan entity.RcloneCopyConfig, 14)

var ErrorTasks = list.New()
