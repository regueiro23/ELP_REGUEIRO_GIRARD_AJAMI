# Jarnac! 
<div style="text-align: center;">
    <img width="100%" src="https://i.ibb.co/Qp4dq11/A-convertir-online-video-cutter-com.gif">
</div>

## Prérequis
- Un terminal
- Node pour utiliser JS
Nous allons supposer que vous disposez du premier. Si vous disposez pas du deuxième, lancez la commande suivante dans votre terminal:

```bash
sudo apt update
sudo apt install nodejs
```
### Lancement 
Il suffit de taper les commandes suivantes et vous serez prêt à jouer.
```bash
npm install
node main
```
Nous vous recommandons de mettre le terminal en plein écran pour mieux visualiser le jeu.

## Fonctionalités

### Animation stylé
Afin de mettre en valeur le projet, nous avons travaillé une animation à l'ouverture. Cette animation nous a beaucoup appris sur le fonctionnement des promesses, des callbacks et de la gestion asynchrone en JavaScript. Cette fonction encapsule une animation ASCII dans une promesse, illustrant comment JavaScript permet de structurer des opérations asynchrones de manière claire et efficace. 

### Interface d'utilisateurs
![](https://i.imgur.com/OCj7ItE.png)

<br>
Les joueurs interagissent avec le jeu directement depuis le terminal. Grâce à une interface dynamique, les joueurs sont guidés tout au long de la partie par des indications claires sur les actions possibles et les choix effectués. Elle affiche aussi les plateaux de jeu pour chaque participant, leurs lettres disponibles, ainsi que le nombre d'indices restants. À chaque fin de tour, l'écran se rafraîchit pour offrir une expérience visuelle optimisée.

### Système de logs
Un système de logs est aussi utilisé afin de permettre à tout moment aux joueurs de consulter l'historique des modifications des tableaux et des mots.
Cet historique sera accessible depuis un fichier appelé `jarnac_coups.txt` au dossier de lancement du jeu.

<div style="text-align: center;">
    <img width="90%" src="https://image.noelshack.com/fichiers/2024/05/5/1706838481-capture-d-ecran-du-2024-02-02-02-46-11.png">
</div>

### Vérification des mots
Afin de vérifier les mots insérés par les joueurs, nous utilisons une API externe. A chaque fois qu'un joueur tente de rentrer un mot, une requête est envoyée pour vérifier si le mot existe ou non. Attention : l'API est anglaise donc le jeu doit se jouer en anglais.

Pour faire les requêtes à l'API : Nous utilisons le module axios. Il est préinstallé avec le projet mais s'il disparait, réinstaller avec:

```bash
npm install axios
```
### Indices
Tous les joueurs commencent avec trois indices. Si un joueur est bloqué et qu'il a peu de lettres, il peut demander un indice. On recherche alors parmi toutes les combinaisons de ses lettres s'il est possible de former un mot avec.

### Échange des lettres
Nous avons décidé d'ajouter une règle au jeu pour éviter les passages successifs. Les joueurs peuvent échanger 3 de leurs lettres contre des nouvelles. Ces lettres échangées ne retournent pas dans le jeu, et chaque joueur est limité à 3 échanges.

## Structure du code

Afin d'éviter un seul fichier de 500+ lignes de code, nous avons décidé de découper le jeu de la façon suivante:

#### main.js
C'est le fichier principal du jeu. Quand il est exécuté, lance l'animation d'introduction, initialise les plateaux des joueurs avec les noms correspondants, la distribution des lettres et l'affichage. C'est le seul fichier que doit manipuler l'utilisateur (lancer le programme).

#### utils.js
Ici on regroupe une collection de fonctions pratiques pour le jeu, allant de la récupération d'entrées utilisateur à la vérification de l'existence de mots via une API. Les ici on trouve les outils pour faire les interactions avec les joueurs, gérer les lettres et mots, et même pour agréablement afficher les plateaux. 

#### jarnac.js
Dans ce fichier sont les fonctions associées aux coups de jarnac
#### animation.js
C'est le module appelé pour lancer l'animation d'introduction.
#### next_turn.js
Ce fichier sert a actualiser le jeu avec chaque interaction et passage de tour. C'est ici que les fonctions du module `utils.js` sont majoritairement appliquées. Il gère aussi la fin et la fermeture du jeu.

#### check_word.js
Apporte des fonction pour vérifier si la modification d'un plateau est permise ou pas.

#### cartes.js
Ce petit module sert à créer le sac des lettres et a en piocher quand nécessaire.
