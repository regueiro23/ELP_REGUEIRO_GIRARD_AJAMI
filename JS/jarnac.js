import * as sac.js
function creerJeu(){
    let sac = [];
    function addLetter(letter, n){
        if (n>0) {
            sac.push(letter);
            addLetter(letter,n-1);
        }
    }
    //remplisage du sac
    addLetter('A', 14);
    addLetter('B', 4);
    addLetter('C', 7);
    addLetter('D', 5);
    addLetter('E', 19);
    addLetter('F', 2);
    addLetter('G', 4);
    addLetter('H', 2);
    addLetter('I', 11);
    addLetter('J', 1);
    addLetter('K', 1);
    addLetter('L', 6);
    addLetter('M', 5);
    addLetter('N', 9);
    addLetter('O', 8);
    addLetter('P', 4);
    addLetter('Q', 1);
    addLetter('R', 10);
    addLetter('S', 7);
    addLetter('T', 9);
    addLetter('U', 8);
    addLetter('V', 2);
    addLetter('W', 1);
    addLetter('X', 1);
    addLetter('Y', 1);
    addLetter('Z', 2);

    //melanger le sac
    sac = melanger(sac);
    tapis = [];

    const etat = {
        sac: sac,
        tapis: tapis,
    }

    console.log(sac)
}

function melanger(sac){
    for(let i = sac.length -1; i>0;i--){
        const j = Math.floor(Math.random()*(i+1));
        [sac[i],sac[j]] = [sac[j], sac[i]]
    }
    return sac; 
}
let jeu = creerJeu()


/*
creerJeu(commencerjeu)

// ou ca 
creerJeu()
.then(commencerjeu)
.catch(erreurCreationJeu)

//utiliser ca  ^
async creerJeu(){}

await creerJeu()
commencerJeu()*/