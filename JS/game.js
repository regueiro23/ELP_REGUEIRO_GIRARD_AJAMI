const fs = require('fs');
const readline = require('readline');
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});
 
const pioche = [];
 
/*
// Initialiser la pioche avec les lettres et leurs quantités
const lettres = {
  'A': 14, 'B': 4, 'C': 7, 'D': 5, 'E': 19, 'F': 2, 'G': 4, 'H': 2, 'I': 11, 'J': 1,
  'K': 1, 'L': 6, 'M': 5, 'N': 9, 'O': 8, 'P': 4, 'Q': 1, 'R': 10, 'S': 7, 'T': 9,
  'U': 8, 'V': 2, 'W': 1, 'X': 1, 'Y': 1, 'Z': 2
};
*/

const lettres = {
    'A': 4, 'E': 4, 'L': 4, 'P': 4,
  };

const dictionnaire = ["apple"];

Object.keys(lettres).forEach(lettre => {
  for (let i = 0; i < lettres[lettre]; i++) {
    pioche.push(lettre);
  }
});
 

// Fonction pour tirer des cartes de la pioche
function tirerCarte(nombre_cartes=6) {
  const cartesTirees = [];
  for (let i = 0; i < nombre_cartes; i++) {
    if (pioche.length === 0) {
      break; // Arrête si la pioche est vide
    }
    const index = Math.floor(Math.random() * pioche.length);
    cartesTirees.push(...pioche.splice(index, 1));
  }
  return cartesTirees;
}
 
class JarnacGame {
    constructor() {
        this.players = [
            { name: "Joueur 1", pool: Array(8).fill(null).map(() => Array(9).fill(' ')), letters: tirerCarte() },
            { name: "Joueur 2", pool: Array(8).fill(null).map(() => Array(9).fill(' ')), letters: tirerCarte() }
        ];
        this.currentPlayerIndex = 0;
    }
 
    startGame() {
        this.displayStatus();
        this.nextTurn();
    }
 
    displayStatus() {
        console.log('\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n')
        this.players.forEach(player => {
            console.log(`Lettres de ${player.name}: ${player.letters}`);
            console.log(`Pool de ${player.name}:`);
            player.pool.forEach((line, index) => console.log(`${index + 1}: ${line}`));
        });
    }
 
    nextTurn() {
        const currentPlayer = this.players[this.currentPlayerIndex];
        rl.question(`${currentPlayer.name}, choisissez une ligne pour écrire ou modifier votre mot (ou entrez 'passer' pour terminer votre tour) : `, (lineNumber) => {
            if (lineNumber.toLowerCase() === 'passer') {
                this.changePlayer();
                return;
            }
 
            // Convertit lineNumber en nombre
            const parsedNumber = parseInt(lineNumber, 10);
 
            // Vérifie si parsedNumber est un entier et si la conversion est réussie (pas NaN)
            if (!Number.isInteger(parsedNumber) || isNaN(parsedNumber)) {
                console.log("Veuillez entrer un numéro de ligne valide.");
                this.nextTurn();
                return;
            }
 
            const lineIndex = parsedNumber - 1;
            if (lineIndex < 0 || lineIndex >= currentPlayer.pool.length) {
                console.log("Ligne invalide. Réessayez.");
                this.nextTurn();
                return;
            }
 
            rl.question("Entrez votre mot : ", (word) => {
                word = word.toUpperCase().split('');
                if (word.length>9){
                    console.log("Le mot est trop long, il doit faire max 9 lettres");
                    this.nextTurn();
                    return;
                }
                if (word.length<3){
                    console.log("Le mot est trop court, il doit faire min 3 lettres");
                    this.nextTurn();
                    return;
                }
                word=this.completerListe(word);
                if (!this.canFormWord(word, currentPlayer.letters)) {
                    console.log("Le mot contient des lettres qui ne sont pas disponibles dans votre liste de lettres.");
                    this.nextTurn();
                    return;
                }

                if (!this.MotExiste(word)) {
                    console.log("Ce mot n'existe pas!");
                    this.nextTurn();
                    return;
                }

                const existingWord = currentPlayer.pool[lineIndex];
 
 
                if (this.estVide(existingWord) && this.canFormWord(word, currentPlayer.letters)) {
 
                    currentPlayer.pool[lineIndex] = this.completerListe(word);
                    currentPlayer.letters=this.retirerLettres(word,currentPlayer.letters)
 
                    // Verifier si lettres bonus gagnées
 
                } else if (!this.estVide(existingWord) && this.canModifyWord(word, existingWord)) {
 
                    

                    let lettresAretirer=this.retirerLettres(existingWord,word);

                    for (let lettre of lettresAretirer){
                        const index = currentPlayer.letters.indexOf(lettre)
                        if (index !== -1) {
                            currentPlayer.letters.splice(index, 1);
                        }
                    }
                        
                            
                    

                    currentPlayer.pool[lineIndex] = word;
                    // Verifier si lettres bonus gagnées
 
                } else {
                    console.log("Modification invalide. Réessayez.");
                    this.nextTurn();
                    return;
                }
 
                this.displayStatus();
                this.nextTurn();
            });
        });
    }


    estVide(ligne) {
        return ligne.every(caractere => caractere === ' ');
    }
 
    retirerLettres(mot, lettres) {
        const motCopie = [...mot]; // Créer une copie de 'mot' pour éviter de modifier le tableau original
        return lettres.filter(lettre => {
            const index = motCopie.indexOf(lettre);
            if (index > -1) {
                motCopie.splice(index, 1); // Retirer l'élément trouvé de la copie de 'mot'
                return false;
            }
            
            return true;
        })
    }
    
 
    completerListe(liste) {
        const longueurCible = 9;
        while (liste.length < longueurCible) {
            liste.push(' ');
        }
        return liste;
    }
 
    canFormWord(listeLettres, lettresDisponibles) {
        listeLettres = listeLettres.filter(item => item !== ' ');
        return listeLettres.every(lettre => lettresDisponibles.includes(lettre));
    } 

    canModifyWord(newWord, existingWord) {
        const filteredNewWord = newWord.filter(item => item !== ' ');
        const newWordSet = new Set(filteredNewWord);
        const filteredExistingWord = existingWord.filter(item => item !== ' ');

        for (const letter of filteredExistingWord){
            if (!newWordSet.has(letter)){
                return false;
            }
        }

        return true;
    }

    MotExiste(word) {
        const trimmedWord = word.filter(char => char !== ' ').join('').toLowerCase();
        return dictionnaire.includes(trimmedWord);
    }
 
    getAddedLetters(newWord, existingWord) {
        return newWord.slice(existingWord.length);
    }
 
    updateLetters(player, usedLetters) {
        for (const char of usedLetters) {
            player.letters = player.letters.replace(char, '');
        }
        player.letters += tirerCarte(usedLetters.length);
    }
 
    changePlayer() {
        this.currentPlayerIndex = (this.currentPlayerIndex + 1) % this.players.length;
        this.displayStatus();
        this.nextTurn();
    }
}
 
 
const game = new JarnacGame();
game.startGame();