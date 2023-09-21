package main

import (
	"bufio"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/pflag"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

func main() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n $ minebak [options] world_name\n")
		fmt.Fprintf(os.Stderr, "Example:\n $ minebak MyWorld --addr 127.0.0.1 --port 21\n\n")
		pflag.PrintDefaults()
	}

	var addr, port, user string
	pflag.StringVar(&addr, "addr", "", "The address of the MultiCraft's FTP server (Refer to your hosting service documentation)")
	pflag.StringVar(&port, "port", "", "The port of the MultiCraft's FTP server (Refer to your hosting service documentation)")
	pflag.StringVar(&user, "user", "", "The username used for the login into your MultiCraft's FTP server.\nGenerally this username is not the one you use to login into your MultiCraft account (Refer to your hosting service documentation)")

	var passwordFilePath string
	var password string
	pflag.StringVar(&passwordFilePath, "password-file", "", "The path to a file containing exactly one line: the password for your MultiCraft's FTP server")

	var outputFilePath string
	pflag.StringVar(&outputFilePath, "output", "", "The file name for the backup, by default a timestamp is appended at the end of the backup name (see --with-timestamp)")

	var withTimestamp bool
	pflag.BoolVar(&withTimestamp, "with-timestamp", false, "Whether or not to append a timestamp in the name of the backup")

	var noInput, quiet bool
	pflag.BoolVar(&noInput, "no-input", false, "If set it prevents the program from asking any input interactively")
	pflag.BoolVar(&quiet, "quiet", false, "if set disables the output of human readable information like the download size and the progress downloadBar")

	pflag.Parse()

	if pflag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "You must specify your minecraft world name\n")
		pflag.Usage()
		os.Exit(1)
	} else if pflag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Too many arguments, refer to \"minebak --help\" for some examples\n")
		pflag.Usage()
		os.Exit(1)
	}

	var reader *bufio.Reader
	if !noInput {
		reader = bufio.NewReader(os.Stdin)
	}

	if addr == "" {
		if noInput {
			fmt.Fprintf(os.Stderr, "You must specify the address of your MultiCraft's FTP server\n")
			pflag.Usage()
			os.Exit(2)
		} else {
			fmt.Print("MultiCraft FTP server IP address: ")
			text, _ := reader.ReadString('\n')
			addr = strings.TrimRight(text, "\r\n")
		}
	}
	if port == "" {
		if noInput {
			fmt.Fprintf(os.Stderr, "You must specify the port of your MultiCraft's FTP server\n")
			pflag.Usage()
			os.Exit(2)
		} else {
			fmt.Print("MultiCraft FTP server port: ")
			text, _ := reader.ReadString('\n')
			port = strings.TrimRight(text, "\r\n")
		}
	}
	if user == "" {
		if noInput {
			fmt.Fprintf(os.Stderr, "You must specify the user name to login into your MultiCraft's FTP server\n")
			pflag.Usage()
			os.Exit(2)
		} else {
			fmt.Print("MultiCraft FTP server user name: ")
			text, _ := reader.ReadString('\n')
			user = strings.TrimRight(text, "\r\n")
		}
	}
	if passwordFilePath == "" {
		if noInput {
			fmt.Fprintf(os.Stderr, "You must specify the password file to login into your MultiCraft's FTP server\n")
			pflag.Usage()
			os.Exit(2)
		} else {
			fmt.Print("MultiCraft FTP server password: ")
			text, _ := reader.ReadString('\n')
			password = strings.TrimRight(text, "\r\n")
		}
	} else {
		file, err := os.Open(passwordFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "The password file was specified but it cannot be opened: %s\n", passwordFilePath)
			log.Fatal(err)
		}

		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "The password file was specified but the password cannot be read: %s\n", passwordFilePath)
			log.Fatal(err)
		}

		password = string(content)
		file.Close()
	}

	worldName := pflag.Arg(0)

	conn, err := ftp.Dial(addr+":"+port, ftp.DialWithTimeout(15*time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "FTP Connection: Failed to create a connection to %s. Please retry again\n", addr+":"+port)
		log.Fatal(err)
	}

	err = conn.Login(user, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FTP Connection: Login failed. Please retry again\n")
		log.Fatal(err)
	}

	if !quiet {
		fmt.Println("FTP Connection: Connection established")
		fmt.Println("Calculating download size...\nThis may take from a few seconds to 3-4 minutes depending on your connection and your world size")
	}

	fileList, err := conn.NameList("/")
	if err != nil {
		fmt.Fprintf(os.Stderr, "FTP Connection: Impossible to establish if your world: %s is present in the MultiCraft FTP server\n", worldName)
		log.Fatal(err)
	}
	if !slices.Contains(fileList, worldName) {
		fmt.Fprintf(os.Stderr, "FTP Connection: Your world \"$s\" is not present in the MultiCraft FTP server\n", worldName)
		log.Fatal("World directory not found in the MultiCraft FTP server")
	}

	var SizeBar *progressbar.ProgressBar
	if !quiet {
		SizeBar = progressbar.DefaultBytes(-1, "Calculating...")
	}

	walker := conn.Walk(worldName)
	sizeAccumulator := uint64(0)
	for walker.Next() {
		if walker.Stat().Type == 0 {
			sizeAccumulator += walker.Stat().Size
		}
		if !quiet {
			SizeBar.Add64(int64(walker.Stat().Size))
		}
	}

	SizeBar.Close()

	backupSizeBytes := sizeAccumulator
	backupSizeGiB := float64(sizeAccumulator) / float64(1073741824)

	var outputFileName string
	if outputFilePath != "" {
		outputFileName = outputFilePath
	} else {
		outputFileName = worldName
	}

	if withTimestamp {
		outputFileName = outputFileName + time.Now().Format("20060102")
	}

	err = os.MkdirAll(outputFileName, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create a local folder to contain your world backup at %s", outputFileName)
		log.Fatal()
	}
	if !quiet {
		fmt.Printf("The download size is %.3f GiB\n", backupSizeGiB)
		fmt.Println("Starting the download")
	}

	var downloadBar *progressbar.ProgressBar
	if !quiet {
		downloadBar = progressbar.DefaultBytes(int64(backupSizeBytes), "Downloading...")
	}

	walker = conn.Walk(worldName)
	for walker.Next() {
		if walker.Stat().Type == 0 {
			path := outputFileName + "/" + strings.Join(strings.Split(walker.Path(), "/")[1:], "/")
			err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			f, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}
			r, err := conn.Retr(walker.Path())
			if err != nil {
				log.Fatal(err)
			}
			content, err := io.ReadAll(r)
			if err != nil {
				log.Fatal(err)
			}
			if err := r.Close(); err != nil {
				log.Fatal()
			}
			if _, err := f.Write(content); err != nil {
				log.Fatal()
			}
			if err := f.Close(); err != nil {
				log.Fatal()
			}
		}
		if !quiet {
			downloadBar.Add64(int64(walker.Stat().Size))
		}
	}
	downloadBar.Finish()

	if !quiet {
		fmt.Println("Download terminated.")
	}
}
