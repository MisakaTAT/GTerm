package sftp

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MisakaTAT/GTerm/backend/pkg/exec"
	commonssh "github.com/MisakaTAT/GTerm/backend/pkg/ssh"
	"github.com/MisakaTAT/GTerm/backend/types"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Handler struct {
	SSHClient        *ssh.Client
	SFTPClient       *sftp.Client
	IsConnected      bool
	LastErrorTime    time.Time
	HomeDir          string
	PermissionsCache *PermissionsCache
	execAdapter      *exec.Adapter
}

func NewSFTPHandler() *Handler {
	return &Handler{
		IsConnected:      false,
		PermissionsCache: NewPermissionsCache(),
	}
}

func (h *Handler) Connect(conf *commonssh.Config) error {
	if h.IsConnected {
		return errors.New("already connected to SFTP server")
	}

	client, err := commonssh.NewSSHClient(conf)
	if err != nil {
		return err
	}
	h.SSHClient = client

	h.execAdapter = exec.New(h.SSHClient)
	// 预加载权限信息
	h.PermissionsCache.preloadPermissions(h.SSHClient, h.execAdapter)

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		_ = h.SFTPClient.Close()
		return err
	}

	h.SFTPClient = sftpClient
	h.IsConnected = true

	homeDir, err := h.GetHomeDirectory()
	if err != nil {
		slog.Warn("Failed to get home directory: %v, using / as default", err)
		h.HomeDir = "/"
	} else {
		h.HomeDir = homeDir
	}

	return nil
}

func (h *Handler) Close() {
	if !h.IsConnected {
		return
	}
	if h.SFTPClient != nil {
		_ = h.SFTPClient.Close()
		h.SFTPClient = nil
	}
	if h.SSHClient != nil {
		_ = h.SSHClient.Close()
		h.SSHClient = nil
	}
	h.IsConnected = false
}

func (h *Handler) ListRemoteFiles(path string) ([]*types.FileTransferItemInfo, error) {
	entries, err := h.SFTPClient.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make([]*types.FileTransferItemInfo, 0, len(entries))
	for _, entry := range entries {
		info := &types.FileTransferItemInfo{
			Name:        entry.Name(),
			Size:        entry.Size(),
			IsDir:       entry.IsDir(),
			ModTime:     entry.ModTime().Format(time.RFC3339),
			Permissions: entry.Mode().String(),
			Owner:       "unknown",
			Group:       "unknown",
		}

		if stat, ok := entry.Sys().(*sftp.FileStat); ok {
			info.Owner = h.PermissionsCache.GetUsername(stat.UID)
			info.Group = h.PermissionsCache.GetGroupName(stat.GID)
		}
		files = append(files, info)
	}

	return files, nil
}

func (h *Handler) UploadFile(localPath, remotePath string, progressCallback func(int64, int64)) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer func(localFile *os.File) {
		_ = localFile.Close()
	}(localFile)

	localFileInfo, err := localFile.Stat()
	if err != nil {
		return err
	}

	if localFileInfo.IsDir() {
		return errors.New("directory upload not supported")
	}

	totalSize := localFileInfo.Size()

	remoteDir := filepath.Dir(remotePath)
	if remoteDir != "." && remoteDir != "/" {
		err = h.SFTPClient.MkdirAll(remoteDir)
		if err != nil {
			return err
		}
	}

	remoteFile, err := h.SFTPClient.Create(remotePath)
	if err != nil {
		return err
	}

	defer func(remoteFile *sftp.File) {
		_ = remoteFile.Close()
	}(remoteFile)

	_, err = io.Copy(localFile, &ProgressReader{
		Reader:           remoteFile,
		TotalSize:        totalSize,
		BytesRead:        0,
		ProgressCallback: progressCallback,
	})

	return err
}

func (h *Handler) DownloadFile(remotePath, localPath string, progressCallback func(int64, int64)) error {
	remoteFileInfo, err := h.SFTPClient.Stat(remotePath)
	if err != nil {
		return err
	}

	if remoteFileInfo.IsDir() {
		return errors.New("directory download not supported")
	}

	totalSize := remoteFileInfo.Size()

	remoteFile, err := h.SFTPClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer func(remoteFile *sftp.File) {
		_ = remoteFile.Close()
	}(remoteFile)

	localDir := filepath.Dir(localPath)
	if localDir != "." {
		err = os.MkdirAll(localDir, 0755)
		if err != nil {
			return err
		}
	}

	localFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer func(localFile *os.File) {
		_ = localFile.Close()
	}(localFile)

	_, err = io.Copy(localFile, &ProgressReader{
		Reader:           remoteFile,
		TotalSize:        totalSize,
		BytesRead:        0,
		ProgressCallback: progressCallback,
	})

	return err
}

func (h *Handler) CreateRemoteFolder(path string) error {
	return h.SFTPClient.MkdirAll(path)
}

func (h *Handler) GetRemoteFileSize(path string) (int64, error) {
	fileInfo, err := h.SFTPClient.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func (h *Handler) GetHomeDirectory() (string, error) {
	result := h.execAdapter.Run("pwd")
	if !result.Success() {
		slog.Error("Failed to execute pwd command: %v", result.Error())
	}

	homeDir := result.Unwrap()
	if homeDir == "" {
		result = h.execAdapter.Run("echo $HOME")
		if !result.Success() {
			slog.Error("Failed to get HOME environment variable: %v", result.Error())
		}
		homeDir = result.Unwrap()
	}

	if homeDir == "" {
		return "/", nil
	}

	if !strings.HasPrefix(homeDir, "/") {
		homeDir = "/" + homeDir
	}

	if strings.HasSuffix(homeDir, "/") && len(homeDir) > 1 {
		homeDir = homeDir[:len(homeDir)-1]
	}

	return homeDir, nil
}

func (h *Handler) ProcessPath(path string) (string, error) {
	if path == "" {
		if h.HomeDir != "" {
			return h.HomeDir, nil
		}
		return "/", nil
	}
	if path == "/" {
		return "/", nil
	}
	if path[0] != '/' {
		path = "/" + path
	}

	return path, nil
}

func (h *Handler) JoinRemotePaths(base, relPath string) (string, error) {
	basePath, err := h.ProcessPath(base)
	if err != nil {
		return "", err
	}
	if relPath == "" {
		return basePath, nil
	}
	if relPath[0] == '/' {
		return relPath, nil
	}
	if basePath[len(basePath)-1] != '/' {
		basePath += "/"
	}

	return basePath + relPath, nil
}

type ProgressReader struct {
	Reader           io.Reader
	TotalSize        int64
	BytesRead        int64
	ProgressCallback func(int64, int64)
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if n > 0 {
		pr.BytesRead += int64(n)
		if pr.ProgressCallback != nil {
			pr.ProgressCallback(pr.BytesRead, pr.TotalSize)
		}
	}
	return n, err
}
