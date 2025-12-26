package wpdatabase

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/nikosmpi/mozaik-cli/wpconfig"
	"golang.org/x/crypto/ssh"
)

type progressWriter struct {
	Total      int64
	Downloaded int64
}

type warningFilter struct {
	Writer io.Writer
}

func (pw *progressWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	pw.Downloaded += int64(n)
	if pw.Total > 0 {
		percentage := float64(pw.Downloaded) / float64(pw.Total) * 100
		fmt.Printf("\rΣυγχρονισμός: %.2f%% (%.2f MB / %.2f MB)", percentage, float64(pw.Downloaded)/(1024*1024), float64(pw.Total)/(1024*1024))
	} else {
		fmt.Printf("\rΣυγχρονισμός: %.2f MB", float64(pw.Downloaded)/(1024*1024))
	}
	return n, nil
}

func (f *warningFilter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !strings.Contains(trimmed, "Using a password on the command line interface can be insecure") {
			_, err := fmt.Fprintln(f.Writer, line)
			if err != nil {
				return 0, err
			}
		}
	}
	return len(p), nil
}

func SyncStagingToLocal(config wpconfig.WPConfig) error {
	localMysqlPath := config.MySQLPath
	if _, err := os.Stat(localMysqlPath); os.IsNotExist(err) {
		localMysqlPath = "mysql"
	}
	key, err := os.ReadFile(config.Staging.SSHKeyPath)
	if err != nil {
		return fmt.Errorf("Δεν βρέθηκε το SSH key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("Λάθος στο SSH key: %v", err)
	}
	sshConfig := &ssh.ClientConfig{
		User:            config.Staging.SSHUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshHost := config.Staging.SSHHost
	if !strings.Contains(sshHost, ":") {
		sshHost += ":22"
	}

	client, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		return fmt.Errorf("Αποτυχία σύνδεσης στον Server: %v", err)
	}
	defer client.Close()

	// Get DB size for progress
	var dbSize int64
	sizeSession, err := client.NewSession()
	if err == nil {
		defer sizeSession.Close()
		sizeSession.Stderr = &warningFilter{Writer: os.Stderr}
		sizeCmd := fmt.Sprintf("mysql -u%s -p%s -e \"SELECT SUM(data_length + index_length) FROM information_schema.TABLES WHERE table_schema = '%s'\" -sN",
			config.Staging.DBUser, config.Staging.DBPass, config.Staging.DBName)
		out, err := sizeSession.Output(sizeCmd)
		if err == nil {
			fmt.Sscanf(string(out), "%d", &dbSize)
		}
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Αποτυχία SSH Session: %v", err)
	}
	defer session.Close()

	remoteCmd := fmt.Sprintf("mysqldump -u%s -p%s --single-transaction --quick --add-drop-table %s",
		config.Staging.DBUser, config.Staging.DBPass, config.Staging.DBName)

	remoteReader, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Δεν μπορώ να πάρω το pipe εξόδου: %v", err)
	}
	session.Stderr = &warningFilter{Writer: os.Stderr}

	var args []string
	args = append(args, fmt.Sprintf("-u%s", config.DBUser))
	if config.DBPass != "" {
		args = append(args, fmt.Sprintf("-p%s", config.DBPass))
	}
	args = append(args, config.DBName)

	localCmd := exec.Command(localMysqlPath, args...)

	pw := &progressWriter{Total: dbSize}
	localCmd.Stdin = io.TeeReader(remoteReader, pw)
	localCmd.Stderr = &warningFilter{Writer: os.Stderr}

	fmt.Println("Ξεκινάει ο συγχρονισμός...")
	if err := localCmd.Start(); err != nil {
		return fmt.Errorf("Δεν μπόρεσε να ξεκινήσει η τοπική MySQL (τσέκαρε το path): %v", err)
	}

	if err := session.Run(remoteCmd); err != nil {
		return fmt.Errorf("Σφάλμα κατά το remote dump: %v", err)
	}

	if err := localCmd.Wait(); err != nil {
		return fmt.Errorf("Σφάλμα κατά την τοπική εγγραφή (Import): %v", err)
	}

	// Force 100% progress display
	if dbSize > 0 {
		fmt.Printf("\rΣυγχρονισμός: 100.00%% (%.2f MB / %.2f MB)", float64(pw.Downloaded)/(1024*1024), float64(pw.Downloaded)/(1024*1024))
	}
	fmt.Println("\nΟ συγχρονισμός ολοκληρώθηκε επιτυχώς!")
	return nil
}
