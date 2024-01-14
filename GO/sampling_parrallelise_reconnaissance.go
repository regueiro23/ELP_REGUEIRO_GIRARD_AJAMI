package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Kagami/go-face"
)

const dataDirSamples = "testdata/samples"
const dataDirImages = "testdata/images"
const dataDirModels = "testdata/models"
const dataDirResultats = "testdata/resultats"

var photosBase []string
var photosComparees []string
var labels []string

// On définit une structure pour les tâches
type Task struct {
	Index int
	Image string
}

func main() {
	startTime := time.Now()

	photosBase = recupererFichiers(dataDirSamples)
	photosComparees = recupererFichiers(dataDirImages)

	fmt.Println("Reconnaissance 3000")

	// On initalise le modèle de reconnaissance
	rec, err := face.NewRecognizer(dataDirModels)
	if err != nil {
		fmt.Println("Impossible d'initialiser le modèle de reconnaissance faciale")
	}
	defer rec.Close()
	fmt.Println("Modèle de reconnaissance initialisé")

	////////////////////////////////////////////////////////////////////
	////
	////		Analyse des visages samples AVEC GOROUTINES
	////
	////////////////////////////////////////////////////////////////////

	var samples []face.Descriptor
	labels = make([]string, len(photosBase))
	var identifiants []int32

	// Creation du waitgroup et du mutex pour remplir nos listes sans décallage causés par la parralélisation
	var wgParallel sync.WaitGroup
	var mu sync.Mutex

	for indice, image := range photosBase {
		wgParallel.Add(1)
		go func(index int, imageName string) {
			defer wgParallel.Done()

			fmt.Println("Analyse/Sampling du visage présent sur", imageName)
			localVisage := sampleVisage(rec, imageName)

			// On utilise mutex pour remplir les listes sans chevauchement
			mu.Lock()
			samples = append(samples, localVisage.Descriptor)
			labels[index] = strings.TrimSuffix(imageName, ".jpg")
			identifiants = append(identifiants, int32(index))
			mu.Unlock()
		}(indice, image)
	}

	// On attend que toutes les routines se terminent
	wgParallel.Wait()

	// On envoie nos samples au modèle de reconnaissance
	rec.SetSamples(samples, identifiants)

	////////////////////////////////////////////////////////////////////
	////
	////			Analyse des photos à comparer
	////
	////////////////////////////////////////////////////////////////////

	var wg sync.WaitGroup
	numWorkers := 8

	taskChannel := make(chan Task, numWorkers)

	// On initialise les workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, rec, taskChannel, &wg)
	}

	// On remplit les tâches
	for index, image := range photosComparees {
		taskChannel <- Task{Index: index, Image: image}
	}
	close(taskChannel)

	// On attend que tous les workers aient fini
	wg.Wait()

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Temps d'exécution total de la comparaison : %s\n", elapsedTime)

	////////////////////////////////////////////////////////////////////
	////
	////			Démarrage de la session TCP
	////
	////////////////////////////////////////////////////////////////////

	fmt.Println("Démarrage de la session TCP")

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erreur lors de la création du serveur:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Serveur en attente de connexions sur localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur lors de la connexion du client:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func worker(workerID int, rec *face.Recognizer, taskChannel <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChannel {
		fmt.Printf("Worker %d solving problem: %d - %s\n", workerID, task.Index, task.Image)
		// Résoud la tâche
		solveTask(task, rec)
	}
}

func solveTask(task Task, rec *face.Recognizer) {
	image := task.Image

	// On analyse tous les visages présents dans l'image, ils sont stockés dans une liste appelée visagesComparés
	visagesCompares := sampleMultiplesVisages(rec, image)

	if visagesCompares == nil {
		fmt.Println("Aucun visage sur cette image")
	} else {
		// Pour chaque visage présent dans l'image, on le compare avec nos samples avec un seuil de comparaison fixe
		for _, visage := range visagesCompares {
			IDVisage := rec.ClassifyThreshold(visage.Descriptor, 0.15)
			// Si le visage ne correspond à aucun sample, on ne peut pas classifier, sinon on l'enregistre dans le bon dossier, sous le bon nom, etc...
			if IDVisage < 0 {
				fmt.Println("Ne peut pas classifier")
			} else {
				fmt.Println(labels[IDVisage])
				fmt.Println(visage.Rectangle.Min)
				fmt.Println(visage.Rectangle.Max)
				enregistreCopieRectangle(image, visage.Rectangle.Min.X, visage.Rectangle.Min.Y, visage.Rectangle.Max.X, visage.Rectangle.Max.Y, labels[IDVisage], "compare_"+image)
			}
		}
	}
}

func sampleVisage(rec *face.Recognizer, photo string) face.Face {
	fichierImage := filepath.Join(dataDirSamples, photo)
	visage, err := rec.RecognizeSingleFile(fichierImage)
	if err != nil {
		log.Fatalf("Can't recognize: %v", photo)
	}
	return *visage
}

func sampleMultiplesVisages(rec *face.Recognizer, photo string) []face.Face {
	fichierImage := filepath.Join(dataDirImages, photo)
	liste_visages, err := rec.RecognizeFile(fichierImage)
	if err != nil {
		log.Fatalf("Can't recognize: %v", photo)
	}
	return liste_visages
}

func enregistreCopieRectangle(inputImageName string, x1, y1, x2, y2 int, outputDir string, outputImageName string) {
	// Ouvrir le fichier image
	inputImagePath := filepath.Join(dataDirImages, inputImageName)
	inputImageFile, err := os.Open(inputImagePath)
	if err != nil {
		log.Fatalf("Erreur ouverture fichier pour modification : %s", inputImageName)
	}
	defer inputImageFile.Close()

	// Décoder le fichier image
	img, _, err := image.Decode(inputImageFile)
	if err != nil {
		log.Fatalf("Erreur décodage fichier pour modification : %s", inputImageName)
	}

	// Créer un nouvel image RGBA pour dessiner le rectangle rouge
	bounds := img.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, img, image.Point{}, draw.Over)

	// Dessiner les contours du rectangle rouge
	red := color.RGBA{255, 0, 0, 255} // Rouge pur, sans mélange
	draw.Draw(rgbaImg, image.Rect(x1, y1, x2, y1+1), &image.Uniform{red}, image.Point{}, draw.Over)
	draw.Draw(rgbaImg, image.Rect(x1, y1, x1+1, y2), &image.Uniform{red}, image.Point{}, draw.Over)
	draw.Draw(rgbaImg, image.Rect(x2-1, y1, x2, y2), &image.Uniform{red}, image.Point{}, draw.Over)
	draw.Draw(rgbaImg, image.Rect(x1, y2-1, x2, y2), &image.Uniform{red}, image.Point{}, draw.Over)

	// Créer le dossier de sortie :
	cheminDossier := filepath.Join(dataDirResultats, outputDir)
	_, err = os.Stat(cheminDossier)

	if os.IsNotExist(err) {
		// Le dossier n'existe pas, le créer
		err := os.MkdirAll(cheminDossier, os.ModePerm)
		if err != nil {
			log.Fatalf("Erreur création de dossier pour l'image: %s", inputImageName)
		}
		fmt.Printf("Dossier '%s' créé.\n", cheminDossier)
	} else if err != nil {
		// Une erreur s'est produite lors de la vérification
		log.Fatalf("Erreur lors de la vérification du dossier pour l'image: %s", inputImageName)
	} else {
		// Le dossier existe déjà
		fmt.Printf("Le dossier '%s' existe déjà.\n", cheminDossier)
	}

	// Créer le fichier de sortie
	outputImagePath := filepath.Join(dataDirResultats, outputDir, outputImageName)
	outputImageFile, err := os.Create(outputImagePath)
	if err != nil {
		log.Fatalf("Erreur lors de la création du fichier de sortie pour l'image: %s", inputImageName)
	}
	defer outputImageFile.Close()

	// Encoder l'image résultante au format JPEG
	err = jpeg.Encode(outputImageFile, rgbaImg, nil)
	if err != nil {
		log.Fatalf("Erreur d'encodage du fichier pour l'image: %s", inputImageName)
	}
}

func recupererFichiers(dossierSource string) []string {
	var listeFichiers []string

	// Lire le contenu du dossier
	contenuDossier, err := ioutil.ReadDir(dossierSource)
	if err != nil {
		log.Fatalf("Soucis collecte fichiers")
	}

	// Parcourir les fichiers du dossier
	for _, fichier := range contenuDossier {
		// Vérifier si le fichier a l'extension .jpg
		if fichier.IsDir() == false && filepath.Ext(fichier.Name()) == ".jpg" {
			listeFichiers = append(listeFichiers, fichier.Name())
		}
	}

	return listeFichiers
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Erreur lors de la lecture du message:", err)
			return
		}

		message = strings.TrimSpace(message)

		switch {
		case message == "1":
			sendFileList(conn, labels)
			conn.Write([]byte("FinListe\n"))
		case strings.HasPrefix(message, "2"):
			folderName := strings.TrimSpace(strings.TrimPrefix(message, "2 "))
			sendFolder(conn, folderName)
		case message == "3":
			conn.Close()
			return
		default:
			fmt.Println("Commande non reconnue:", message)
		}
	}
}

func sendFileList(conn net.Conn, liste []string) {
	fileList := ""

	for _, file := range liste {
		fileList += file + "\n"
	}

	conn.Write([]byte(fileList))
}

func sendFolder(conn net.Conn, foldername string) {
	chemin := filepath.Join(dataDirResultats, foldername)
	err := compressFolderToZip(chemin)
	if err != nil {
		fmt.Println("Erreur de compression zip")
	}
	filename := chemin + ".zip"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		conn.Write([]byte("Erreur lors de l'ouverture du fichier\n"))
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	fileSizeStr := strconv.FormatInt(fileSize, 10)
	fmt.Println(fileSizeStr)

	conn.Write([]byte(fileSizeStr + "\n"))

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Erreur lors de la lecture du fichier:", err)
			return
		}
		conn.Write(buffer[:n])
	}

	err = os.Remove(filename)
	if err != nil {
		fmt.Println("Erreur lors de la suppression du fichier:", err)
		return
	}

	fmt.Println("Fichier envoyé avec succès !")
}

func compressFolderToZip(source string) error {
	target := source + ".zip"
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name, excluding the root folder
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}
