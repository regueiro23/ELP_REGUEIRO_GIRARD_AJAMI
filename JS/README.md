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
sudo apt-get install npm
```
### Lancement 
Il suffit de taper la commande suivante et vous serez prêt à jouer.
```bash
node main.js
```
Nous vous recommandons de mettre le terminal en plein écran pour mieux visualiser le jeu.

## Fonctionement
### Interface d'utilisateurs:

Les joueurs peuvent interagir avec le jeu via le terminal. Celui-ci montrera les deux tableaux du jeu,
Les lettres disponibles pour chaque joueur ainsi que le nombre d'indices restants.

Le système dynamique de logs permet de guider les joueurs en affichant les choix possibles du jeu et les options choisies. Après le tour du joueur, l'écran s'actualise pour optimiser l'affichage.

<div style="text-align: center;">
    <img width="100%" src="https://image.noelshack.com/fichiers/2024/05/5/1706838481-capture-d-ecran-du-2024-02-02-02-46-11.png">
</div>

### Verification des mots
Afin de vérifier les mots insérés par les joueurs, nous utilisons une API externe. A chaque fois qu'un joueur tente de rentrer un mot, une requête est envoyée pour vérifier si le mot existe ou non. Attention : l'API est anglaise donc le jeu est en anglais.

Pour faire les requetes à l'API : Nous utilisons le module axios. Il est pré-installé avec le projet mais s'il disparrait, réinstaller avec npm install axio.

Tous les joueurs commencent avec trois indices. Si un joueur est bloqué et qu'il a peu de lettres, il peut demander un indice. On recherche alors parmis toutes les combinaisons de ses lettres s'il est possible de former un mot avec.

Afin de mettre en valeur le projet, nous avons travaillé une animation à l'ouverture. Cette animation nous a beaucoup appris sur le fonctionnement des promesses, des callbacks et de la gestion asynchrone en JavaScript. Cette fonction encapsule une animation ASCII dans une promesse, illustrant comment JavaScript permet de structurer des opérations asynchrones de manière claire et efficace. 

### Structure du code

Afin d'eviter un seul fichier de 500+ lignes de code, nous avons decidé de decouper le jeu de la facon suivante:

#### main.js

C'est le fichier principal du jeu. Quand il est executé, lance l'animation d'introduction, initialise les plateaux des joueurs avec les noms correspondants, la distribution des lettres et l'affichage. C'est le seul fichier que doit manipuler l'utilisateur (lancer le programme).

#### utils.js

Ici on regroupe une collection de fonctions pratiques pour le jeu, allant de la récupération d'entrées utilisateur à la vérification de l'existence de mots via une API. Les ici on trouve les outilis pour faire les interactions avec les joueurs, gérer les lettres et mots, et même pour agreablement afficher les plateaux. 

#### jarnac.js
Dans ce fichier sont les fonctions associés aux coups de jarnac
#### animation.js
C'est le module appelé pour lancer l'animation d'introduciton.
#### next_turn.js
Ce fichier sert a actualiser le jeu avec chaque interaction et passage de tour. C'est ici que le fonctions du module `utils.js` sont majoritairement apliquées. Il gere aussi la fin et la fermeture du jeu.

#### check_word.js
Apporte des fonction pour verifier si la modification d'un plateau est permise ou pas.

#### cartes.js
Ce petit module sert a creer le sac des lettres et a en piocher quand neccessaire.