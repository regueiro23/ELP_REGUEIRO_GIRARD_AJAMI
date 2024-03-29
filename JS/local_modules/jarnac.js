module.exports = { demandeJarnac, coupdeJarnac }

const { questionAsync, detectionFin, estVide, completerListe, trouverLigneVide, retirerLettres, MotExiste } = require('./utils.js')
const { finPartie, changePlayer, nextTurn, displayStatus } = require('./next_turn.js')
const { canFormWord, canModifyWord } = require('./check_word.js')
const { tirerCarte } = require('./cartes.js')
const fs = require('fs').promises;

////////////////////////////////////////////////////////////////
// jarnac.js : Contient les fonctions liées au coup de Jarnac //
////////////////////////////////////////////////////////////////



//demandeJarnac : Demande au joueur si il souhaite réaliser un coup de Jarnac une fois le tour de l'autre joueur est terminé
async function demandeJarnac(players, currentPlayerIndex,pioche) {

    const currentPlayer = players[currentPlayerIndex]; //Le joueur qui vient de terminer son tour
    const otherPlayerIndex = (currentPlayerIndex + 1) % players.length; //On récupère l'index de l'autre joueur
    const otherPlayer = players[otherPlayerIndex]; //Le joueur qui va réaliser (ou pas) le coup de Jarnac

    let reponse_jarnac = await questionAsync(`${otherPlayer.name}, voulez-vous tenter un coup de Jarnac, un double coup de Jarnac, ou rien? (tapez simple, double ou passer):  `);

    reponse_jarnac = reponse_jarnac.toLowerCase();
    if (reponse_jarnac === "simple"){
        console.log(`${otherPlayer.name}, vous allez tenter un coup de Jarnac.`);
        coupdeJarnac(players, currentPlayerIndex, 1 ,pioche);
        return;
    }

    else if (reponse_jarnac === "double"){
        console.log(`${otherPlayer.name}, vous allez tenter un double coup de Jarnac.`);
        coupdeJarnac(players, currentPlayerIndex, 2 ,pioche);
        return;
    }

    else if (reponse_jarnac === "passer"){
        //Si le joueur décide de ne pas faire de coup de Jarnac, on vérifie si son opposant a gagné
        if (detectionFin(currentPlayer.pool)){
            finPartie(players);
            return;
        }
        changePlayer(players,currentPlayerIndex,pioche);
        return;
    } 

    else{
        console.log("Erreur : tapez simple, double ou passer")
        demandeJarnac(players, currentPlayerIndex,pioche);
        return;
    }
    
}

//coupdeJarnac : Fonction qui gère le coup de Jarnac
//               Prend en paramètre un nombre égal à 1 ou 2 ( c'est un coup simple ou double )
async function coupdeJarnac(players, currentPlayerIndex, number,pioche){

    const currentPlayer = players[currentPlayerIndex];
    const otherPlayerIndex = (currentPlayerIndex + 1) % players.length;
    const otherPlayer = players[otherPlayerIndex];

    let lineNumber = await questionAsync('Choisissez une ligne pour écrire ou modifier le mot (ou entrez passer pour terminer le coup de Jarnac) : ');
    
    if (lineNumber.toLowerCase() === 'passer') {
        changePlayer(players, currentPlayerIndex,pioche);
        return;
    }
    const parsedNumber = parseInt(lineNumber, 10);
    if (!Number.isInteger(parsedNumber) || isNaN(parsedNumber)) {
        console.log("Veuillez entrer un numéro de ligne valide.");
        coupdeJarnac(players, currentPlayerIndex,number,pioche);
        return;
    }
    const lineIndex = parsedNumber - 1;
    if (lineIndex < 0 || lineIndex >= otherPlayer.pool.length) {
        console.log("Veuillez entrer un numéro de ligne valide.");
        coupdeJarnac(players, currentPlayerIndex,number,pioche);
        return;
    }


    let word = await questionAsync("Entrez votre mot : ");
    word = word.toUpperCase().split('');
    if (word.length>9){
        console.log("Le mot est trop long, il doit faire max 9 lettres");
        coupdeJarnac(players, currentPlayerIndex,number,pioche);
        return;
    }
    if (word.length<3){
        console.log("Le mot est trop court, il doit faire min 3 lettres");
        console.log(word);
        coupdeJarnac(players, currentPlayerIndex,number,pioche);
        return;
    }

    const existingWord = currentPlayer.pool[lineIndex];
    let actionText = "";

    if (estVide(existingWord)) {

        if (!canFormWord(word, currentPlayer.letters)) {
            console.log("Le mot contient des lettres qui ne sont pas disponibles dans la liste de lettres.");
            coupdeJarnac(players, currentPlayerIndex,number,pioche);
            return;
        }
        if (!(await MotExiste(word))) {
            console.log("Ce mot n'existe pas! PS : Le mot entré doit être en anglais");
            coupdeJarnac(players, currentPlayerIndex,number,pioche);
            return;
        }

        actionText = `Coup de Jarnac !! ${otherPlayer.name} vole les lettres de ${currentPlayer.name} pour écrire "${word.join('')}" !!\n`;
        
        word=completerListe(word);

        //Le mot volé est placé dans la première ligne vide du joueur qui réalise le coup
        lineIndexOther=trouverLigneVide(otherPlayer.pool);
        otherPlayer.pool[lineIndexOther] = word;
        
        //On retire les lettres ajoutées pour former le mot de la main de l'adversaire
        currentPlayer.letters=retirerLettres(word,currentPlayer.letters);

        //Le joueur gagne une nouvelle lettre
        otherPlayer.letters.push(tirerCarte(1,pioche)[0]);

    } else {
        
        let result = canModifyWord(word , existingWord.filter(item => item !== ' ') , currentPlayer.letters);
        if (!result){
            console.log("Cette modification est impossible. Réessayez!");
            coupdeJarnac(players, currentPlayerIndex,number,pioche);
            return;
        }
        if (!(await MotExiste(word))) {
            console.log("Ce mot n'existe pas! PS : Le mot entré doit être en anglais");
            coupdeJarnac(players, currentPlayerIndex,number,pioche);
            return;
        }

        //On retire les lettres ajoutées pour former le mot de la main de l'adversaire
        let lettresAretirer=retirerLettres(existingWord,word);
        for (let lettre of lettresAretirer){
            const index = currentPlayer.letters.indexOf(lettre);
            currentPlayer.letters.splice(index, 1);
        }

        actionText = `Coup de Jarnac !! ${otherPlayer.name} vole la ligne de ${currentPlayer.name} : "${existingWord.filter(item => item !== ' ').join('')}"=>"${word.join('')}" !!\n`;

        word=completerListe(word);

        //Le mot volé est placé dans la première ligne vide du joueur qui réalise le coup
        lineIndexOther = trouverLigneVide(otherPlayer.pool);
        otherPlayer.pool[lineIndexOther] = word; 
        currentPlayer.pool[lineIndex] = Array(9).fill(' ');

        //Le joueur gagne une nouvelle lettre
        otherPlayer.letters.push(tirerCarte(1,pioche)[0]);

    }

    //On logs le tour qui vient de passer
    await fs.appendFile('jarnac_coups.txt', actionText);

    //On vérifie si le joueur qui vient de réaliser le coup a gagné
    if (detectionFin(otherPlayer.pool)){
        finPartie(players);
        return;
    }

    //Si le coup était simple, on passe au tour du joueur
    if (number === 1){
        changePlayer(players, currentPlayerIndex,pioche);
        return;
    }

    //Si le coup était double, on rappelle la fonction coup de Jarnac avec un coup simple
    else if (number === 2){
        displayStatus(players);
        coupdeJarnac(players, currentPlayerIndex,1,pioche);
        return;
    }
}
