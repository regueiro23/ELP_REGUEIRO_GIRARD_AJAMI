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
	"strings"
	"sync"
	"time"

	"github.com/Kagami/go-face"
)

// Seuil de tolérance pour la reconnaissance : 0 = très précis, 1 = très imprécis.
// Ajuster en fonction de la cohérence du premier jet (valeur recommandée : 0.40)
var seuil_tolerance_reconnaissance float32 = 0.35

// Utilisation de go-routines pour accélerer le sampling des visages de départ
// Mettre à false pour ne pas parralléliser cette tâche
var sampling_parrallelise bool = false

// Nombre de workers pour l'analyse des images
// Ajuster en fonction du CPU pour obtenir des performances max.
// Mettre à 1 pour qu'il ne pas parralléliser cette tâche (peu recommandé, performances très réduites)
var numWorkers int = 8

// Dossiers où sont stockés les samples, les images à comparer,
// les modèles de reconnaissance et les résultats d'analyse
const dataDirSamples = "testdata/samples"
const dataDirImages = "testdata/images"
const dataDirModels = "testdata/models"
const dataDirResultats = "testdata/resultats"

// Listes où seront stockées les photos des samples (photosBase), les photos à comparer (photosComparees)
// les noms associés à tous les visages samplés (labels) ainsi que leurs identifiants (identifiants)
// et les descripteurs des visages samplés (samples)
var photosBase []string
var photosComparees []string
var labels []string
var identifiants []int32
var samples []face.Descriptor

// Structure pour les tâches
type Task struct {
	Index int
	Image string
}

func main() {

	totalStartTime := time.Now()

	fmt.Println("Reconnaissance 3000")

	// On récupère toutes les photos à analyser et on les stocke dans leurs listes respectives
	photosBase = recupererFichiers(dataDirSamples)
	photosComparees = recupererFichiers(dataDirImages)
	labels = make([]string, len(photosBase))

	// On initalise le modèle de reconnaissance
	rec, err := face.NewRecognizer(dataDirModels)
	if err != nil {
		log.Fatalf("Impossible d'initialiser le modèle de reconnaissance")
	}
	defer rec.Close()
	fmt.Println("Modèle de reconnaissance initialisé")

	////////////////////////////////////////////////////////////////////
	////
	////	Analyse des visages samples AVEC ou SANS goroutines
	////
	////////////////////////////////////////////////////////////////////

	samplingStartTime := time.Now()

	if sampling_parrallelise {
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

	} else {
		for indice, image := range photosBase {
			fmt.Println("Analyse/Sampling du visage présent sur", image)
			localVisage := sampleVisage(rec, image)
			samples = append(samples, localVisage.Descriptor)
			labels[indice] = strings.TrimSuffix(image, ".jpg")
			identifiants = append(identifiants, int32(indice))
		}
	}

	// On envoie nos samples au modèle de reconnaissance
	rec.SetSamples(samples, identifiants)

	samplingEndTime := time.Now()
	samplingElapsedTime := samplingEndTime.Sub(samplingStartTime)

	////////////////////////////////////////////////////////////////////
	////
	////			Analyse des photos à comparer
	////
	////////////////////////////////////////////////////////////////////

	analyseStartTime := time.Now()

	// Création du waitgroup pour distribuer les photos à analyser parmis les workers
	var wg sync.WaitGroup
	// Initialisation du channel des tâches à réaliser
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

	analyseEndTime := time.Now()
	analyseElapsedTime := analyseEndTime.Sub(analyseStartTime)

	totalEndTime := time.Now()
	totalElapsedTime := totalEndTime.Sub(totalStartTime)

	fmt.Printf("\nDurée du sampling des visages: %s\n", samplingElapsedTime)
	fmt.Printf("Durée de l'analyse des photos : %s\n", analyseElapsedTime)
	fmt.Printf("Durée totale (initialisation+sampling+analyse) : %s\n\n", totalElapsedTime)

	fmt.Print("Parrallélisation du sampling : ")
	if sampling_parrallelise {
		fmt.Println("Oui")
	} else {
		fmt.Println("Non")
	}

	fmt.Printf("Nombre de workers pour l'analyse : %d\n\n", numWorkers)

	////////////////////////////////////////////////////////////////////
	////
	////			Démarrage de la session TCP
	////
	////////////////////////////////////////////////////////////////////

	fmt.Println("Démarrage de la session TCP")

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Impossible d'ouvrir le serveur TCP sur le port 8080 : ", err)
	}
	defer listener.Close()

	fmt.Println("Serveur en attente de connexions sur localhost:8080")

	for {
		// On ouvre chaque connexion en parrallèle pour pouvoir en recevoir plusieurs
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Impossible d'accepter la connexion TCP client :", err)
		}
		go handleConnection(conn)
	}
}

// worker : Tant que le taskChannel n'est pas vide, les workers prennent de nouvelles tâches et les résolvent
func worker(workerID int, rec *face.Recognizer, taskChannel <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChannel {
		// Chaque worker prend une tâche et la résoud jusqu'à ce que le channel soit vide
		fmt.Printf("Worker %d résoud la tâche: %d - %s\n", workerID, task.Index, task.Image)
		solveTask(task.Image, rec)
	}
}

// solveTask : Pour une image donnée, on analyse tous les visages présents et on les compare à nos samples
func solveTask(image string, rec *face.Recognizer) {
	// visagesCompares : Liste des descripteurs des visages trouvés dans l'image
	visagesCompares := sampleMultiplesVisages(rec, image)

	/*  Si la liste n'est pas vide (= des visages sont détectés dans l'image), on compare chaque visage avec les visages samplés
	précédemment, si la différence entre les visages est sous notre seuil de tolérance, alors le visage est considéré comme reconnu*/
	if visagesCompares != nil {
		for _, visage := range visagesCompares {
			IDVisage := rec.ClassifyThreshold(visage.Descriptor, seuil_tolerance_reconnaissance)
			if IDVisage >= 0 {
				fmt.Println(labels[IDVisage], "reconnu sur l'image", image)
				// On récupère les coordonnés du rectangle qui entoure le visage pour le dessiner dans la fonction enregistreCopieRectangle
				enregistreCopieRectangle(image, visage.Rectangle.Min.X, visage.Rectangle.Min.Y, visage.Rectangle.Max.X, visage.Rectangle.Max.Y, labels[IDVisage], "compare_"+image)
			}
		}
	}
}

// sampleVisage : Analyse une photo avec un seul visage (sample) et renvoie son descripteur
func sampleVisage(rec *face.Recognizer, photo string) face.Face {
	fichierImage := filepath.Join(dataDirSamples, photo)
	visage, err := rec.RecognizeSingleFile(fichierImage)
	if err != nil {
		log.Fatalf("Aucun visage reconnu sur la photo: %v", photo)
	}
	return *visage
}

// sampleMultiplesVisages : Analyse une photo avec un/ou plusieurs visages et renvoie une liste
// 							de tous leurs descripteurs. Permet ensuite de comparer chaque
//							descripteur avec nos samples.
func sampleMultiplesVisages(rec *face.Recognizer, photo string) []face.Face {
	fichierImage := filepath.Join(dataDirImages, photo)
	liste_visages, err := rec.RecognizeFile(fichierImage)
	if err != nil {
		fmt.Println("Erreur lors de la reconnaissance sur la photo: %v", photo)
	}
	return liste_visages
}

// enregistreCopieRectangle : Enregistre les images analysées en y ajoutant un rectangle rouge autour du visage
// 							  de la personne reconnue. L'image modifiée est stockée dans un dossier portant le
//							  nom de la personne reconnue.
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
	red := color.RGBA{255, 0, 0, 255}
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
			log.Fatalf("Erreur création du dossier: %s", outputDir)
		}
		fmt.Printf("Dossier '%s' créé.\n", cheminDossier)
	} else if err != nil {
		// Une erreur s'est produite lors de la vérification
		log.Fatalf("Erreur lors de la vérification du dossier pour l'image: %s", inputImageName)
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

// recupererFichiers : Récupère tous les fichiers présents dans un dossier et les renvoie sous la forme d'une liste
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

// recupererFichiers : Récupère tous les sous-dossiers présents dans un dossier et les renvoie sous la forme d'une liste
func recupererSousDossiers(dossierSource string) []string {
	var dirNames []string

	// Fonction à appeler pour chaque élément dans le dossier.
	visit := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && path != dossierSource { // Vérifier si c'est un sous-dossier et non le dossier racine.
			dirName := filepath.Base(path) // Obtenir uniquement le nom du sous-dossier.
			dirNames = append(dirNames, dirName)
		}
		return nil
	}
	// Parcourir le dossier.
	filepath.Walk(dossierSource, visit) // Les erreurs sont ignorées.
	return dirNames
}

// handleConnection : Gère une connexion client. Si le client envoie "1", on lui renvoie
//					  la liste des personnes reconnues lors de la dernière analyse. S'il
//					  envoie "2" suivi du prénom d'une personne reconnue, on crée une
//					  archive zip contenant toutes les photos de la personne en question.
//					  S'il envoie "3", on coupe la connexion.
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

			sendFileList(conn, recupererSousDossiers(dataDirResultats))
			conn.Write([]byte("FinListe\n"))

		case strings.HasPrefix(message, "2"):
			liste_personnes := recupererSousDossiers(dataDirResultats)
			personne_demandee := strings.TrimSpace(strings.TrimPrefix(message, "2 "))
			// On vérifie qu'on a des photos de la personne demandée
			for _, personne := range liste_personnes {
				if personne == personne_demandee {
					conn.Write([]byte("Sending\n"))
					sendFolder(conn, personne)
					conn.Close()
					return
				}
			}
			conn.Write([]byte("NotFound\n"))
			conn.Close()
			return

		case message == "3":
			conn.Close()
			return
		default:
			fmt.Println("Commande non reconnue:", message)
		}
	}
}

// sendFileListe : Envoie la liste des personnes reconnues sur la connexion avec le client
func sendFileList(conn net.Conn, liste []string) {
	fileList := ""

	for _, file := range liste {
		fileList += file + "\n"
	}

	conn.Write([]byte(fileList))
}

// sendFolder : Compresse le dossier souhaité au format zip et l'envoie sur la connexion avec le client
func sendFolder(conn net.Conn, foldername string) {

	chemin := filepath.Join(dataDirResultats, foldername)
	filename, err := compressFolderToZip(chemin)
	if err != nil {
		fmt.Println("Erreur de compression zip")
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		conn.Write([]byte("Erreur lors de l'ouverture du fichier\n"))
		return
	}
	defer file.Close()

	// Envoi du fichier
	_, err = io.Copy(conn, file)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi", err)
		return
	}

	// Suppression de l'archive
	err = os.Remove(filename)
	if err != nil {
		fmt.Println("Erreur lors de la suppression du fichier:", err)
		return
	}

	fmt.Printf("Les photos de %s ont été envoyé avec succès !\n", foldername)
}

// compressFolderToZip compresse un dossier en zip et renvoie le nom du fichier compressé
func compressFolderToZip(source string) (string, error) {
	target := source + ".zip"

	// Création du fichier cible et d'un writer zip
	f, err := os.Create(target)
	if err != nil {
		return "", err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// On passe à travers tous les fichiers du dossier source
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Création d'un header local
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate

		// On définit le nom du header avec le chemin relatif du fichier à la racine
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// On crée un writer pour le header et on stocke le contenu du fichier
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(headerWriter, file)
		return err
	})

	if err != nil {
		return "", err
	}

	return target, nil
}
