module.exports = {canFormWord, canModifyWord}

//////////////////////////////////////////////////////////////////////////////
// check_words.js : Contient les fonctions liées à la vérification des mots //
//////////////////////////////////////////////////////////////////////////////


//canFormWord : Vérifie si un mot peut être formé à partir des lettres
//              Utilisée dans le cadre où le joueur tente d'écrire sur une ligne vide
function canFormWord(listeLettresMot, lettresDisponibles) {
    const lettresTemp = [...lettresDisponibles];
    for (const lettre of listeLettresMot) {
        const index = lettresTemp.indexOf(lettre);
        if (index !== -1) {
            lettresTemp.splice(index, 1);
        } else {
            return false;
        }
    }
    return true;
}

//canModifyWord : Vérifie si un mot peut être formé à partir des lettres de l'ancien mot et des lettres de la main du joueur
//                Utilisée dans le cadre où le joueur tente d'écrire sur une ligne où il y a déjà un mot
function canModifyWord(newWord, existingWord, letters) {
    //Si le nouveau mot est plus court que l'ancien
    if (newWord.length<=existingWord.length){
        return false;
    }

    // Vérifier si chaque lettre de l'ancien mot est présente dans le nouveau
    for (const letter of existingWord) {
        if (!newWord.includes(letter)) {
            return false;
        }
    }

    // Créer une banque de lettres composée des lettres de 'letters' et 'existingWord'
    const letterBank = [...letters, ...existingWord];

    // Parcourir le nouveau mot et vérifier la disponibilité de chaque lettre dans la banque
    for (const letter of newWord) {
        const index = letterBank.indexOf(letter);
        if (index === -1) {
            return false; // Lettre non disponible dans la banque
        }
        letterBank.splice(index, 1); // Supprimer la lettre utilisée de la banque
    }

    return true; // Toutes les lettres du nouveau mot sont disponibles dans la banque
}