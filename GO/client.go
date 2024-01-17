package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net"
	"net_test/utils"
	"os"
	"path/filepath"
	"strings"
)

// Ce client permet de récupérer la liste des personnes reconnues par le serveur et d'en télécharger les photos
// On utilise une connexion directe TCP pour échanger avec le serveur.
func main() {
	// Initialisation de la connexion TCP sur le port 8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur:", err)
		return
	}
	defer conn.Close()

	// Buffer pour récupérer les input de l'utilisateur
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
			fmt.Print("Quelle célébrité ?\n")
			fileName, _ := reader.ReadString('\n')
			fileName = strings.TrimSpace(fileName)

			sendMessage(conn, "2 "+fileName)

			response, _ := bufio.NewReader(conn).ReadString('\n')
			response = strings.TrimSpace(response)

			if response == "NotFound" {
				fmt.Println("Le serveur n'a pas de photos de", fileName)
			} else if response == "Sending" {
				receiveFile(conn, fileName+".zip")
				fmt.Println("Vous avez téléchargé les photos de", fileName)
			}
			conn.Close()
			return

		case "3":
			sendMessage(conn, "3")
			fmt.Println("Connexion coupée. Fin du programme.")
			conn.Close()
			return
		default:
			fmt.Println("\nCommande non reconnue. Veuillez réessayer.\n")
		}
	}
}

// receiveFile : Reçoit l'archive zip du serveur et la décompresse
func receiveFile(conn net.Conn, fileName string) {

	fo, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Erreur lors de la création fichier recu", err)
		return
	}
	defer fo.Close()

	_, err = io.Copy(fo, conn)
	utils.CheckError(err)

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

// receiveFileList : Recoit et affiche la liste des personne envoyée par le serveur
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

// unzipFile : Décompresse une archive zip dans un dossier du même nom
func unzipFile(zipFilePath string) error {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer r.Close()
	// On retire le .zip pour obtenir le nom du dossier à créer
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
			// Crée le sous-dossier s'il n'existe pas
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

// sendMessage : envoie un message sur la connexion avec le serveur
func sendMessage(conn net.Conn, message string) {
	conn.Write([]byte(message + "\n"))
}
