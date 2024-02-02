module.exports = { demarrage }
const { questionAsync } = require('./utils.js');

///////////////////////////////////////////////////
// animation.js : Contient l'animation de départ //
///////////////////////////////////////////////////

//animationDepart : Crée l'animation de départ
function animationDepart() {
    return new Promise((resolve, reject) => {
        (function() {
    const artAscii= `
     ██  █████  ██████  ██    █  █████   ██████
     ██ ██   ██ ██   ██ ███   █ ██   ██ ██     
     ██ ███████ ██████  █ ██  █ ███████ ██      
██   ██ ██   ██ ██   ██ █  ██ █ ██   ██ ██      
 █████  ██   ██ ██   ██ █   ███ ██   ██  ██████ `;

            const lignes = artAscii.split('\n').slice(1);
            const largeurLettre = 8;
            const mot = "jarnac";
            const largeurTotale = largeurLettre * mot.length;

            function afficherLettre(index, delai) {
                const avantLettres = ' '.repeat(index * largeurLettre);
                const partiesLettre = lignes.map(ligne => avantLettres + ligne.substring(index * largeurLettre, (index + 1) * largeurLettre));
                console.clear();
                console.log(partiesLettre.join('\n'));
            }

            function afficherMotComplet() {
                console.clear();
                console.log(artAscii);
                resolve(); // Résout la promesse une fois l'animation terminée
            }

            function demarrerAnimation(nbRepetitions, delaiInitial) {
                if (nbRepetitions <= 0) {
                    afficherMotComplet();
                    return;
                }

                mot.split('').forEach((lettre, index) => {
                    setTimeout(() => afficherLettre(index, delaiInitial), index * delaiInitial);
                });

                const prochainDelai = delaiInitial * 0.5;
                setTimeout(() => demarrerAnimation(nbRepetitions - 1, prochainDelai), mot.length * delaiInitial);
            }

            demarrerAnimation(20, 500);
        })();
    });
}

//demarrage : Attend la fin de l'animation de départ
async function demarrage() {
    await animationDepart();
    console.log("\n")
    await questionAsync("Appuyez sur entrée pour démarrer....")
}
