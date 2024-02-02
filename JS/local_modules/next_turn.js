module.exports = {nextTurn, displayStatus, finPartie, changePlayer}

const { demandeJarnac } = require('./jarnac.js');
const { questionAsync, estVide, completerListe, retirerLettres, MotExiste, afficheMatrice,chercheIndice } = require('./utils.js');
const { canFormWord, canModifyWord } = require('./check_word.js');
const { tirerCarte } = require('./cartes.js');
const fs = require('fs').promises;

///////////////////////////////////////////////////////////////////////
// nextTurn.js : Contient les fonctions liées à la gestion des tours //
///////////////////////////////////////////////////////////////////////

//nextTurn : Fonction principale du jeu, déroule un tour normal
async function nextTurn(players, currentPlayerIndex,pioche) {

    const currentPlayer = players[currentPlayerIndex];
    const lineNumber = await questionAsync(`${currentPlayer.name}, choisissez une ligne pour écrire ou modifier (entrez 'indice' pour obtenir un indice ou 'passer' pour passer votre tour) : `);
    
    //Si le joueur demande un indice, on vérifie qu'il lui reste des indices et qu'il ait 4 lettres ou moins en main
    if (lineNumber.toLowerCase()==="indice"){
        if (currentPlayer.letters.length<5){
            if (currentPlayer.indices!=0){
                currentPlayer.indices-=1;
                console.log("Recherche de mot à partir de vos lettres:")
                await chercheIndice(currentPlayer.letters)
                nextTurn(players, currentPlayerIndex,pioche);
                return;
            }
            else{
                console.log("Vous n'avez plus d'indices !")
                nextTurn(players, currentPlayerIndex,pioche);
                return;
            }
        }
        else{
            console.log("Vous avez trop de lettres, cherchez encore un peu ;)")
            nextTurn(players, currentPlayerIndex,pioche);
            return;
        }

    }

    //Si le joueur passe son tour, on propose à son adversaire de faire un coup de Jarnac
    if (lineNumber.toLowerCase() === 'passer') {
        demandeJarnac(players, currentPlayerIndex,pioche);
        return;
    }

    const parsedNumber = parseInt(lineNumber, 10);
    if (!Number.isInteger(parsedNumber) || isNaN(parsedNumber)) {
        console.log("Veuillez entrer un numéro de ligne valide.");
        nextTurn(players, currentPlayerIndex,pioche);
        return;
    }
    const lineIndex = parsedNumber - 1;
    if (lineIndex < 0 || lineIndex >= currentPlayer.pool.length) {
        console.log("Ligne invalide. Réessayez.");
        nextTurn(players, currentPlayerIndex,pioche);
        return;
    }


    let word = await questionAsync(`Entrez votre mot : `);
    word = word.toUpperCase().split('');
    if (word.length>9){
        console.log("Le mot est trop long, il doit faire max 9 lettres");
        nextTurn(players, currentPlayerIndex,pioche);
        return;
    }
    if (word.length<3){
        console.log("Le mot est trop court, il doit faire min 3 lettres");
        nextTurn(players, currentPlayerIndex,pioche);
        return;
    }

    const existingWord = currentPlayer.pool[lineIndex];
    let actionText = "";

    if (estVide(existingWord)) {
        
        if (!canFormWord(word, currentPlayer.letters)) {
            console.log("Le mot contient des lettres qui ne sont pas disponibles dans votre liste de lettres.");
            nextTurn(players, currentPlayerIndex,pioche);
            return;
        }
        if (!(await MotExiste(word))) {
            console.log("Ce mot n'existe pas! PS : Le mot entré doit être en anglais");
            nextTurn(players, currentPlayerIndex,pioche);
            return;
        }

        actionText = `${currentPlayer.name} a écrit dans la ligne ${lineIndex + 1}: "${word.join('')}"\n`;
        
        word=completerListe(word);
        currentPlayer.pool[lineIndex] = word;
        
        //On retire les lettres ajoutées pour former le mot de la main du joueur
        currentPlayer.letters=retirerLettres(word,currentPlayer.letters);

        //Le joueur gagne une nouvelle lettre
        currentPlayer.letters.push(tirerCarte(1,pioche)[0] );

    } else {
        
        if (!canModifyWord(word, existingWord.filter(item => item !== ' ') , currentPlayer.letters)){
            console.log("Cette modification est impossible. Réessayez!");
            nextTurn(players, currentPlayerIndex,pioche);
            return;
        }
        if (!(await MotExiste(word))) {
            console.log("Ce mot n'existe pas! PS : Le mot entré doit être en anglais");
            nextTurn(players, currentPlayerIndex,pioche);
            return;
        }

        //On retire les lettres ajoutées pour former le mot de la main du joueur
        let lettresAretirer=retirerLettres(existingWord,word);
        for (let lettre of lettresAretirer){
            const index = currentPlayer.letters.indexOf(lettre);
            currentPlayer.letters.splice(index, 1);
        }            
        
        actionText = `${currentPlayer.name} a modifié la ligne ${lineIndex + 1}: "${existingWord.filter(item => item !== ' ').join('')}"=>"${word.join('')}"\n`;

        word=completerListe(word);
        currentPlayer.pool[lineIndex] = word;

        //Le joueur gagne une nouvelle lettre
        currentPlayer.letters.push(tirerCarte(1,pioche)[0]);

    }

    //On logs le tour qui vient de passer
    await fs.appendFile('jarnac_coups.txt', actionText);

    displayStatus(players);
    //C'est encore au même joueur de jouer
    nextTurn(players, currentPlayerIndex,pioche);

}

//displayStatus : Affiche les mains et les indices des joueurs ainsi que les plateaux de jeu
function displayStatus(players) {
    console.clear();
    console.log("%c ", "font-size: 1px; padding: 166.5px 250px; background-size: 500px 333px; background: no-repeat url(https://cdn.cultura.com/cdn-cgi/image/width=1280/media/pim/45_249049_2_10_FR.jpg);");
    players.forEach(player => {

        console.log(`Lettres de ${player.name}: ${player.letters}`);
        console.log(`Nombre d'indices restants : ${player.indices}`);
        console.log(`Pool de ${player.name}:`);
        afficheMatrice(player.pool)
        console.log("\n");

    });
}

//changePlayer : Appelée lorsque le tour passe d'un joueur à un autre
function changePlayer(players, currentPlayerIndex,pioche) {
    currentPlayerIndex = (currentPlayerIndex + 1) % players.length;
    displayStatus(players);
    nextTurn(players, currentPlayerIndex,pioche);
}


//finPartie : Appelée quand une pool est remplie. Calcule et affiche les points et annonce le gagnant
function finPartie(players){ 

    let points_joueur0 = 0;
    let points_joueur1 = 0;

    console.log(`Fin de partie, un pool a été complété ! Qui sera notre grand vainqueur.........`);

    for (const line of players[0].pool) {
        points_joueur0 += (line.filter(item => item !== ' ').length)**2;
    }
    console.log(`${players[0].name} a obtenu ${points_joueur0}.`);

    
    for (const line of players[1].pool) {
        points_joueur1 += (line.filter(item => item !== ' ').length)**2;
    }
    console.log(`${players[1].name} a obtenu ${points_joueur1}.`);

    if (points_joueur0 > points_joueur1){
        console.log(`Félicitations! ${players[0].name} remporte le jeu.`);
    }
    else{
        console.log(`Félicitations! ${players[1].name} remporte le jeu.`);
    }
    process.exit(0);
    
}

