package cliftp

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"github.com/jlaffaye/ftp"
	"github.com/kataras/golog"
)

type FTPOptions struct {
	User string
	Word string
	Server string
	Port string
	Conn *ftp.ServerConn
}

type FileData struct {
	Name string
	Size uint64
	Type string
}

// The first step in using the roftp package is to get a new logged in connection,
// and cache it locally
func NewFTPConn(opts *FTPOptions) (error) {
	conn, err := ftp.Connect(opts.Server + ":" + opts.Port)
	if err != nil {
		golog.Println(err, "Error connecting to FTP Server")
		return err
	}

	opts.Conn = conn

	err = login(opts)
	if err != nil {
		golog.Println(err, "Error logging in to ftp server")
		return err
	}
	return nil
}

// login on the supplied basic connection
func login(opts *FTPOptions) error {
	return opts.Conn.Login(opts.User, opts.Word)
}

// Change to the serverPath directory and List files
// Provide an already logged in connection
// ListFiles will change directory to the listed directory
func ListFiles(conn *ftp.ServerConn, serverPath string) (filesData []FileData, err error) {
	currPath, err := ChDir(conn, serverPath)
	if err != nil {
		golog.Println(err, "Unable to obtain current directory")
		return filesData, err
	}

	entries, err := conn.List("")
	if  err != nil {
		golog.Println(err, "Error listing files", "currentDir", currPath)
		return nil, err
	}
	golog.Println(len(entries), "file(s) found at", currPath)
	for _, entry := range entries {
		fileType := "other"
		switch entry.Type {
		case ftp.EntryTypeFile:
			fileType = "file"
		case ftp.EntryTypeFolder:
			fileType = "directory"
		}
		filesData = append(filesData, FileData{ Name: entry.Name, Size: entry.Size, Type: fileType })
	}
	return
}

// Change directory and return the new path or err
func ChDir(conn *ftp.ServerConn, serverPath string) (currPath string, err error) {
	if err = conn.ChangeDir(serverPath); err != nil {
		golog.Println(err, "Error changing directory on ftp server", "requestedPath", serverPath)
		return "", err
	}
	currPath, err = conn.CurrentDir()
	if err != nil {
		golog.Println(err, "Unable to obtain server's current directory after changing directory")
		return "", err
	}
	return
}

// Upload file to the server
// conn should be already logged in and current directory changed to desired dir on server
// Server path is dest path (without filename) on server
func UploadFile(conn *ftp.ServerConn, srcFullPath, serverPath string, destFilename ...string) error {
	file, err := os.Open(srcFullPath)
	if err != nil {
		golog.Println(err, "Unable to open file for upload")
		return err
	}
	defer file.Close()

	if len(destFilename) > 0 {
		serverPath = filepath.Join(serverPath, destFilename[0])
	}
	// Upload
	golog.Println("Uploading sermon:", srcFullPath)
	err = conn.Stor(serverPath, file)
	if err != nil {
		golog.Println(err, "Error uploading file", "actual_server_path", serverPath)
		return err
	}
	return err
}

// Download and write file from server
func DownloadFiles(conn *ftp.ServerConn, serverPath string) (err error) {

	items, err := ListFiles(conn, serverPath)
	if err != nil { 
		golog.Println(err, "Could not obtain dir entries by name") 
		return err
	}

	for _, item := range items {
		if item.Type != "file" || item.Name == "." || item.Name == ".." { continue }
		golog.Println("Downloading", item.Name)
		err = DownloadFile(conn, serverPath, item.Name)
		if err != nil {	
			golog.Println(err.Error())
			return err
		}
	}
	return nil
}

func DownloadFile(conn *ftp.ServerConn, serverPath, destName string) error {
	data, err := DownloadFileBuffer(conn, filepath.Join(serverPath, destName))
	if err != nil {	return err	}
	ioutil.WriteFile(destName, data, 0664)
	return nil
}

// Download file from server as []byte
func DownloadFileBuffer(conn *ftp.ServerConn, serverPath string) (fileData []byte, err error) {
	golog.Println("Downloading", serverPath, "to buffer")
	resp, err := conn.Retr(serverPath)
	if err != nil {
		golog.Println(err, "Error downloading file from server", "remote_file", serverPath)
		return nil, err
	}
	//resp.SetDeadline(time.Now().Add(time.Minute * 15))
	fileData, err = ioutil.ReadAll(resp)
	if err != nil {
		golog.Println(err, "Error reading file from server", "remote_file", serverPath)
		return nil, err
	}
	golog.Println(len(fileData), "bytes read")
	resp.Close()
	return
}

func FTPQuit(conn *ftp.ServerConn) error {
	return conn.Quit()
}
