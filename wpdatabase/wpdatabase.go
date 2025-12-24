package wpdatabase

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/nikosmpi/mozaik-cli/wpconfig"
	"golang.org/x/crypto/ssh"
)

func SyncStagingToLocal(config wpconfig.WPConfig) {
	localMysqlPath := `C:\Programs\wamp64\bin\mysql\mysql9.1.0\bin\mysql.exe`
	key, err := os.ReadFile(config.Staging.SSHKeyPath)
	if err != nil {
		log.Fatalf("Δεν βρέθηκε το SSH key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Λάθος στο SSH key: %v", err)
	}
	sshConfig := &ssh.ClientConfig{
		User:            config.Staging.SSHUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", config.Staging.SSHHost, sshConfig)
	if err != nil {
		log.Fatalf("Αποτυχία σύνδεσης στον Server: %v", err)
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Αποτυχία SSH Session: %v", err)
	}
	defer session.Close()
	remoteCmd := fmt.Sprintf("mysqldump -u%s -p%s --single-transaction --quick --add-drop-table %s",
		config.Staging.DBUser, config.Staging.DBPass, config.Staging.DBName)
	remoteReader, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Δεν μπορώ να πάρω το pipe εξόδου: %v", err)
	}
	session.Stderr = log.Writer()
	var args []string
	args = append(args, fmt.Sprintf("-u%s", config.DBUser))
	if config.DBPass != "" {
		args = append(args, fmt.Sprintf("-p%s", config.DBPass))
	}
	args = append(args, config.DBName)
	localCmd := exec.Command(localMysqlPath, args...)
	localCmd.Stdin = remoteReader
	localCmd.Stderr = log.Writer()
	fmt.Println("Ξεκινάει ο συγχρονισμός...")
	if err := localCmd.Start(); err != nil {
		log.Fatalf("Δεν μπόρεσε να ξεκινήσει η τοπική MySQL (τσέκαρε το path): %v", err)
	}
	if err := session.Run(remoteCmd); err != nil {
		log.Fatalf("Σφάλμα κατά το remote dump: %v", err)
	}
	if err := localCmd.Wait(); err != nil {
		log.Fatalf("Σφάλμα κατά την τοπική εγγραφή (Import): %v", err)
	}
	fmt.Println("Ο συγχρονισμός ολοκληρώθηκε επιτυχώς!")
}
