package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func sendMessage(conn net.Conn, message string) {
	conn.Write([]byte(message + "\n"))
}

func receiveFileList(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF || strings.TrimSpace(line) == "FinListe" {
			break
		}
		if err != nil {
			fmt.Println("Erreur lors de la lecture de la liste de fichiers:", err)
			return
		}

		fmt.Print(line)
	}
}

func receiveFile(conn net.Conn, fileName string) {
	// Recevoir la taille du fichier
	sizeStr, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la taille du fichier:", err)
		return
	}

	// Convertir la taille du fichier en entier
	sizeStr = strings.TrimSpace(sizeStr)
	fileSize, err := strconv.ParseInt(sizeStr, 10, 64)
	fmt.Println(fileSize)

	file, err := os.Create(fileName)
	defer file.Close()

	buffer := make([]byte, 1024)
	totalReceived := int64(0)

	for totalReceived < fileSize {
		fmt.Println(totalReceived)
		n, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Erreur lors de la lecture depuis la connexion:", err)
			return
		}

		totalReceived += int64(n)
	}

	fmt.Println("Fichier reçu avec succès :", fileName)

	err = unzipFile(fileName)
	if err != nil {
		fmt.Println("Erreur lors de la décompression")
	}
	err = os.Remove(fileName)
	if err != nil {
		fmt.Println("Erreur lors de la suppression du fichier:", err)
		return
	}
	fmt.Println("Décompression du fichier terminée, suppression de l'archive")
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Tapez 1 pour récupérer la liste des célébrités")
		fmt.Println("Tapez 2 pour télécharger les photos d'une célébrité")
		fmt.Println("Tapez 3 pour couper la connexion et fermer le programme.")
		fmt.Print("Votre choix : ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			sendMessage(conn, "1")
			receiveFileList(conn)
		case "2":
			fmt.Print("Lequel : ")
			fileName, _ := reader.ReadString('\n')
			fileName = strings.TrimSpace(fileName)

			sendMessage(conn, "2 "+fileName)
			receiveFile(conn, fileName+".zip")

			fmt.Println("Vous avez téléchargé le fichier ", fileName)
		case "3":
			sendMessage(conn, "3") // Informe le serveur que le client veut couper la connexion
			fmt.Println("Connexion coupée. Fin du programme.")
			time.Sleep(1 * time.Second)
			return
		default:
			fmt.Println("Commande non reconnue. Veuillez réessayer.")
		}
	}
}

func unzipFile(zipFilePath string) error {
	// Vérifie si le fichier passé en paramètre est un fichier zip
	if !strings.HasSuffix(zipFilePath, ".zip") {
		return fmt.Errorf("Le fichier n'est pas une archive .zip")
	}

	// Ouvre le fichier zip
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Obtient le nom de l'archive sans l'extension .zip
	archiveName := strings.TrimSuffix(filepath.Base(zipFilePath), ".zip")

	// Crée le dossier de destination s'il n'existe pas
	destFolder := filepath.Join(filepath.Dir(zipFilePath), archiveName)
	if err := os.MkdirAll(destFolder, 0755); err != nil {
		return err
	}

	// Parcourt tous les fichiers dans le zip et les extrait dans le dossier de destination
	for _, f := range r.File {
		destFilePath := filepath.Join(destFolder, f.Name)

		if f.FileInfo().IsDir() {
			// Crée le dossier s'il n'existe pas
			os.MkdirAll(destFilePath, 0755)
			continue
		}

		// Crée le dossier parent s'il n'existe pas
		if err := os.MkdirAll(filepath.Dir(destFilePath), 0755); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Crée le fichier dans le dossier de destination
		file, err := os.Create(destFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copie le contenu du fichier du zip vers le fichier extrait
		_, err = io.Copy(file, rc)
		if err != nil {
			return err
		}
	}

	return nil
}
