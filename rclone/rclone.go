package rclone

import (
	"os/exec"
)

func Uploading(remoteName, threadCount, logFilePath, uploadingFilePath, remoteFilePath string) error {
	cmd := exec.Command("rclone", "-v", "copy",
		"--transfers", threadCount, "--log-file", logFilePath,
		uploadingFilePath, remoteName+":"+remoteFilePath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
