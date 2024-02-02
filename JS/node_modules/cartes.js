module.exports = {creerPioche, tirerCarte}

//////////////////////////////////////////////////////////////////////
// cartes.js : Contient les fonctions liées aux cartes ( =lettres ) //
//////////////////////////////////////////////////////////////////////

//creerPioche : Génère la pioche de jeu à partir des quantités présentes dans les règles du jeu
function creerPioche(pioche){
    const lettres = {
        'A': 14, 'B': 4, 'C': 7, 'D': 5, 'E': 19, 'F': 2, 'G': 4, 'H': 2, 'I': 11, 'J': 1,
        'K': 1, 'L': 6, 'M': 5, 'N': 9, 'O': 8, 'P': 4, 'Q': 1, 'R': 10, 'S': 7, 'T': 9,
        'U': 8, 'V': 2, 'W': 1, 'X': 1, 'Y': 1, 'Z': 2
      };
      Object.keys(lettres).forEach(lettre => {
        for (let i = 0; i < lettres[lettre]; i++) {
          pioche.push(lettre);
        }
      });

}

//tirerCarte : Tire un nombre de cartes donné en paramètres dans la pioche et les renvoie sous forme de liste
function tirerCarte(nombre_cartes,pioche) {
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