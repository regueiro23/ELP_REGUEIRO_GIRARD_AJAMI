const { nextTurn,displayStatus } = require('./local_modules/next_turn.js');
const { creerPioche, tirerCarte } = require('./local_modules/cartes.js');
const { questionAsync } = require('./local_modules/utils.js');
const { demarrage } = require('./local_modules/animation.js')

////////////////////////////////////////////////////
// main.js : Fichier principal. Démarre la partie //
////////////////////////////////////////////////////


//La classe Jarnac permet d'initier les mains/plateaux des joueurs et de lancer la partie
class Jarnac {
    constructor() {
        this.pioche=[]
        creerPioche(this.pioche)
        this.players = [
            { pool: Array(8).fill(null).map(() => Array(9).fill(' ')), letters: tirerCarte(6,this.pioche), indices:3 },
            { pool: Array(8).fill(null).map(() => Array(9).fill(' ')), letters: tirerCarte(6,this.pioche), indices:3 }
        ];
        this.currentPlayerIndex = 0;

    }

    async startGame() {

        let nom_joueur1 = await questionAsync("Nom du premier joueur : ");
        let nom_joueur2 = await questionAsync("Nom du deuxième joueur : ");
        this.players[0]['name']=nom_joueur1;
        this.players[1]['name']=nom_joueur2;

        displayStatus(this.players);
        nextTurn(this.players, this.currentPlayerIndex,this.pioche);

    }
}


async function debutJeu(){
    //On initialise l'animation de départ et on attend qu'elle se termine
    await demarrage()

    const game = new Jarnac();
    game.startGame();
}


debutJeu()
