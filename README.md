<details>
  <summary>Voir GO</summary>
# Reconnaissance et isolation d'individus sur une large banque d'images

L'objectif de ce projet est de mettre en place un système de reconnaissance faciale déployable sur une très grande banque d'images.

En donnant en entrée des photos d'individus que l'on souhaite reconnaitre, on peut ensuite les retrouver dans une grande banque d'images. Le programme détoure et trie les photos où les individus ont été reconnus.

<div style="text-align: center;">
    <img width="100%" src="https://image.noelshack.com/fichiers/2024/03/3/1705527356-faces.jpg">
</div>

##  Références externes
Nous utilisons la librairie de reconnaissance faciale go-face développé par Kagami ainsi que les modèles qu'il a entrainés en utilisant dlib, il est possible de le retrouver ici :
 - [GitHub : Kagami/go-face](https://github.com/Kagami/go-face)


## Installation

Voir les requirements de la librairie go-face ci-dessus.

Tous les autres paquets utilisés sont inclus nativement

## Répertoires

Dans testdata/samples : mettre les samples des personnes à analyser. Une photo par personne au format "nom.jpg".

Dans testdata/images : mettres toutes les photos de la banque d'images à analyser.

ATTENTION : Toutes les images doivent être au format ".jpg". Possibilité d'utiliser un convertisseur si les images ne sont pas au bon format ([Exemple de convertisseur](https://convertio.co/fr/image-converter/))

Les résultats seront stockés dans testdata/resultats

## Test

Des samples et images de test sont fournies. Pour faire tourner le modèle, lancer simplement le main.go :

```bash
  >>> go run main.go
```
A la fin de l'analyse, le programme ouvre un serveur TCP local sur le port 8080. Le client permet ainsi de récupérer les photos analysées en les échangeant via la communication TCP. Pour cela, initialiser le client :

```bash
  >>> go run client.go
Tapez 1 pour récupérer la liste des célébrités
Tapez 2 pour télécharger les photos d'une célébrité
Tapez 3 pour couper la connexion et fermer le programme.
Votre choix : 
```
A partir de là, amusez-vous ;)
## Paramètres

Dans le main.go, quelques paramètres permettent de gérer la parllélisation du programme. Voir notamment :

```go
//Utilisation de go-routines pour accélerer le sampling des visages de départ
//Mettre à false pour ne pas parralléliser cette tâche
var sampling_parrallelise bool = false
```
et
```go
// Nombre de workers pour l'analyse des images
// Ajuster en fonction du CPU pour obtenir des performances max.
// Mettre à 1 pour qu'il ne pas parralléliser cette tâche (peu recommandé, performances très réduites)
var numWorkers int = 8
```
Le réglage du seuil de tolérance du modèle de reconnaissance lui se fait via ce paramètre :
```go
// Seuil de tolérance pour la reconnaissance : 0 = très précis, 1 = très imprécis.
// Ajuster en fonction de la cohérence du premier jet (valeur recommandée : 0.35)
var seuil_tolerance_reconnaissance float32 = 0.35
```
## Parallélisation

Nous avons conduit des tests pour tester les effets de la parrallélisation sur la rapidité d'analyse du programme sur une banque d'images fixées.

- ### Sur les samples :
L'instauration de go-routines sur le sampling initial des visages permet en moyenne de faire gagner entre 20% et 30% de rapidité sur l'étape de sampling initiale. Cette augmentation se faire plus sentir quand le nombre de visages à sampler augmente

- ### Sur l'analyse en elle-même : 

L'endroit où la parallélisation peut avoir le plus gros impact est sur l'analyse en elle-même puisque cette étape peut se voir être répétée sur des milliers d'images. On a voulu mesurer la durée moyenne de l'analyse ( sur dix lancers à chaque fois ) en fonction du nombre de workers crées. Précisons que l'on travaille sur une machine à processeur 8 coeurs.
<div style="text-align: center;">
    <img width="50%" src="https://image.noelshack.com/fichiers/2024/03/3/1705529818-tests.png">
</div>

*Evolution de la durée d'analyse en fonction du nombre de workers*

On remarque une nette diminution du temps d'execution lorsque le nombre de workers se rapproche du nombre de coeurs du processeur !

  <summary>Voir ELM</summary>
