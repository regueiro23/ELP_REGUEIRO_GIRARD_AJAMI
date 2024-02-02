# Jarnac! 

<div style="text-align: center;">
    <img width="100%" src="https://i.ibb.co/Qp4dq11/A-convertir-online-video-cutter-com.gif">
</div>


A noter dans le Readme :

Parler du systeme de logs:
<div style="text-align: center;">
    <img width="100%" src="https://image.noelshack.com/fichiers/2024/05/5/1706838481-capture-d-ecran-du-2024-02-02-02-46-11.png">
</div>

Afin de vérifier les mots insérés par les joueurs, nous utilisons une API externe ( la même que ELM ). A chaque fois qu'un joueur tente de rentrer un mot, une requête est envoyée pour vérifier si le mot existe ou non. Attention : l'API est anglaise donc les mots entrés doivent être anglais

Pour faire les requetes à l'API : Utilisation du module axios. Pré-installé avec le projet mais si il disparrait, réinstaller avec npm install axio ( attention ça fait disparaitre les modules qu'on a créés il faut les remettre à la main )

Tous les joueurs commencent avec trois indices. Si un joueur est bloqué et qu'il a peu de lettres, il peut demander un indice. On recherche alors parmis toutes les combinaisons de ses lettres s'il est possible de former un mot avec.

Afin de mettre en valeur le projet, nous avons travaillé une animation à l'ouverture. Cette animation nous a beaucoup appris sur le fonctionnement des promesses, des callbacks et de la gestion asynchrone en JavaScript. Cette fonction encapsule une animation ASCII dans une promesse, illustrant comment JavaScript permet de structurer des opérations asynchrones de manière claire et efficace. 

Expliquer le découpage en modules :
main.js fichier principal puis :
    - utils.js
    - jarnac.js
    - animation.js
    - next_turn.js
    - check_word.js
    - cartes.js

