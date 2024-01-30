# Guess It!

## Introduction
Guess It!" est une application web interactive développée en Elm. Le jeu consiste à deviner un mot aléatoire à partir de ses définitions.

## Prérequis
- Serveur local (localhost)
- Elm installé sur votre machine

## Configuration
Suivez ces étapes pour configurer et lancer l'application :

### Étape 1 : Configuration du Serveur
Modifiez l'URL du serveur dans le fichier `Main.elm`. Remplacez l'adresse existante par celle de votre serveur local, en veillant à inclure le port approprié (probablement ligne 67 du code).

```elm
-- Dans Main.elm
, Http.get
    { url = "http://localhost:8000/static/mots.txt"  --<--- Remplacez cette ligne avec votre URL
    , expect = Http.expectString WordsLoaded
```

Assurez-vous d'inclure le chemin `/static/mots.txt` après votre adresse.

### Étape 2 : Compilation Elm
Compilez le fichier `Main.elm` avec la commande suivante :

```bash
elm make Main.elm --output main.js
```

Exécutez cette commande dans le répertoire approprié pour éviter les erreurs de chemin.

### Étape 3 : Lancement du Serveur
Si vous n'êtes pas familier avec le lancement d'un serveur Elm, suivez ces instructions :

   Ouvrez un terminal dans le répertoire `ELP_REGUEIRO_GIRARD_AJAMI\ELM`.
   Exécutez la commande suivante :

```bash
elm reactor
```

**Note :** Il est important de lancer le serveur dans le même dossier que le fichier `index.html`.

### Étape 4 : Jouer au Jeu
Après avoir configuré le serveur, accédez à l'adresse de votre serveur local pour commencer à jouer.

## Fonctionnalités du Jeu
"Guess It!" offre une expérience de jeu dynamique avec les caractéristiques suivantes :

- **Devinettes de Mots :** Les joueurs tentent de deviner des mots à partir de définitions fournies.
- **Interface Utilisateur Intuitive :** Une interface claire et facile à naviguer.
- **Système de Score et Gestion du Temps :** Le jeu intègre un système de score et un chronomètre pour augmenter le défi.
