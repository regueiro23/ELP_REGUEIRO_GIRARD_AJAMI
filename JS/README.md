# Jarnac! 

A noter dans le Readme :




Afin de vérifier les mots insérés par les joueurs, nous utilisons une API externe ( la même que ELM ). A chaque fois qu'un joueur tente de rentrer un mot, une requête est envoyée pour vérifier si le mot existe ou non. Attention : l'API est anglaise donc les mots entrés doivent être anglais

Pour faire les requetes à l'API : Utilisation du module axios. Pré-installé avec le projet mais si il disparrait, réinstaller avec npm install axio ( attention ça fait disparaitre les modules qu'on a créés il faut les remettre à la main )



Tous les joueurs commencent avec trois indices. Si un joueur est bloqué et qu'il a peu de lettres, il peut demander un indice. On recherche alors parmis toutes les combinaisons de ses lettres s'il est possible de former un mot avec.
