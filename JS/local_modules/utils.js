module.exports = {questionAsync, estVide, retirerLettres, completerListe, detectionFin, trouverLigneVide, MotExiste,afficheMatrice, chercheIndice}
const axios = require('axios'); // Module axios utilisé pour les requêtes HTTP à l'API
const readline = require('readline');
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// utils.js : Contient les petites fonctions utilisées à tout moment pendant le déroulement du programme//
//////////////////////////////////////////////////////////////////////////////////////////////////////////



//questionAsync : Récupère l'input utilisateur pour une question donnée
function questionAsync(question) {
    return new Promise(resolve => {
      rl.question(question, resolve);
    });
}

//estVide : Vérifie si une liste est remplie de caractères ' '
function estVide(ligne) {
    return ligne.every(caractere => caractere === ' ');
}

//retirerLettres : Renvoie la liste lettres en ayant supprimé toutes les lettres présentes dans la liste mot 
function retirerLettres(mot, lettres) {
    const motCopie = [...mot];
    return lettres.filter(lettre => {
        const index = motCopie.indexOf(lettre);
        if (index > -1) {
            motCopie.splice(index, 1);
            return false;
        }
            
        return true;
    })
}
    
//completerListe : Ajoute à une liste de caractères des caractères vides pour que sa taille finale soit de 9
//                 Permet de remplir les lignes du plateau avec des caractères vides
function completerListe(liste) {
    const longueurCible = 9;
    while (liste.length < longueurCible) {
        liste.push(' ');
    }
    return liste;
}

//detectionFin : Verifie si la matrice donnée en entrée a encore des lignes vides ou non
function detectionFin(pool){
    for (const line of pool) {
        if (estVide(line)){
            return false;
        }
    }
    return true;
}

//trouverLigneVide : Renvoie l'id de la première ligne vide de la liste
function trouverLigneVide(pool){
    for (let i = 0; i < 8; i++){
        const ligne = pool[i];
        if (estVide(ligne)) {
            return i;
        }
    }
}

//MotExiste(mot) : Vérifie si un mot existe en faisant une requête à une API de dictionnaire
async function MotExiste(mot) {
    try {
        mot = mot.join('').toLowerCase()
        const response = await axios.get(`https://api.dictionaryapi.dev/api/v2/entries/en/${mot}`);
        // Si la réponse est réussie (code 200), le mot existe
        return response.status === 200;
    } catch (error) {
        // En cas d'erreur, renvoyer false
        return false;
    }
}

//permute : Génère toutes les permutations d'un tableau
function permute(permutation) {
    let length = permutation.length,
        result = [permutation.slice()],
        c = new Array(length).fill(0),
        i = 1, k, p;

    while (i < length) {
        if (c[i] < i) {
            k = i % 2 && c[i];
            p = permutation[i];
            permutation[i] = permutation[k];
            permutation[k] = p;
            ++c[i];
            i = 1;
            result.push(permutation.slice());
        } else {
            c[i] = 0;
            ++i;
        }
    }
    return result;
}

//chercheIndice : utilise permute() et MotExiste() pour tenter de trouver un mot qui peut etre formé avec une liste de lettres
async function chercheIndice(lettres) {
    let permutations = permute(lettres);
    for (let perm of permutations) {
        if (await MotExiste(perm)) {
            console.log(`Mot trouvé : ${perm.join('')}`);
            return;
        }
    }
    console.log("Aucun mot trouvé.");
}

//afficheMatrice : stylise l'affichage des plateaux des joueurs dans la console
function afficheMatrice(matrice) {
    const ANSI_RESET = "\x1b[0m";
    const ANSI_BLACK_BACKGROUND = "\x1b[40m";
    const ANSI_WHITE = "\x1b[37m";
    const ANSI_GRID_COLOR = "\x1b[90m";

    matrice.forEach((ligne, indexLigne) => {
        let ligneAffichage = ANSI_WHITE + (indexLigne + 1) + " " + ANSI_RESET + "| ";
        ligne.forEach((cellule) => {
            ligneAffichage += ANSI_BLACK_BACKGROUND + ANSI_WHITE + ` ${cellule} ` + ANSI_RESET;
        });
        console.log(ANSI_GRID_COLOR + ligneAffichage + ANSI_RESET);
    });
}