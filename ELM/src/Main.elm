-----------------------
-------GUESS IT !------
-----------------------
 
module Main exposing (..)
 
import Browser
import Html exposing (Html, div, h1, h2, h3, input, text, button, ul, li, p)
import Html.Attributes exposing (style, type_, value, placeholder, disabled)
import Html.Events exposing (onInput, onClick)
import String
import Random
import Http
import Time
import Json.Decode exposing (Decoder, field, list, string, map, map2)

-----------------------
---------MODEL---------
-----------------------

--GameState : état du jeu ( menu de départ, jeu en cours ou jeu terminé)
type GameState  
    = MainMenu
    | InGame
    | GameOver
 
--WordStatus et WordStatusType : pour afficher les mots devinés et passés de manière différente
type alias WordStatus =
    { word : String
    , status : WordStatusType
    }
type WordStatusType
    = Guessed
    | Skipped

styleList : WordStatusType -> List (Html.Attribute Msg) --cette fonction permet de donner un affichage différent aux mots rentrés en fonction de leur statut, c'est-à dire si ils on été devinés ou pas
styleList status =
    case status of
        Guessed -> [] --si Guessed, alors on met le mot tel quel
        Skipped -> [ style "text-decoration" "line-through" ] --on barre les mots qui n'ont pas été devinés (Skipped)

--Model : contient la plupart des variables et paramètres de la partie
type alias Model =
    { randomWord : String
    , guessInput : String
    , guessSuccess : Bool
    , wordList : List String
    , isLoading : Bool
    , definitions : List WordDefinition
    , score : Int
    , gameState : GameState
    , timeElapsed : Int
    , selectedTimer : Int
    , guessedWords : List WordStatus
    }

--initialModel : modèle initial, tout vide et/ou par défaut
initialModel : Model
initialModel =
    { randomWord = ""
    , guessInput = ""
    , guessSuccess = False
    , wordList = []
    , isLoading = True
    , definitions = []
    , score = 0
    , gameState = MainMenu
    , timeElapsed = 60
    , selectedTimer = 60
    , guessedWords = []
    }

--Msg : on définit tous les messages échangés pendant l'execution du programme
type Msg
    = NewRandomWordIndex Int
    | UpdateGuess String
    | StartGame
    | RestartGame
    | SkipWord
    | NewWord
    | WordsLoaded (Result Http.Error String)
    | DefinitionsLoaded (Result Http.Error (List WordDefinition))
    | Tick Time.Posix
    | UpdateSelectedTimer Int


----------------------
---------INIT---------
----------------------

init : () -> (Model, Cmd Msg)
init _ =
    ( initialModel --initialisation du modèle
    , Http.get
        { 
        url = "http://localhost:8000/static/mots.txt"  --url du serveur local, à modifier en fonction de votre configuration
        , expect = Http.expectString WordsLoaded
        }
    )


------------------------
---------UPDATE---------
------------------------

update : Msg -> Model -> (Model, Cmd Msg) --fonction qui définit comment le modèle est mis à jour en fonction des messages reçus
update msg model =
    case msg of
        WordsLoaded (Ok data) ->  --cas où les mots sont bien chargés depuis le serveur
            let
                words = String.split " " data  --isole chaque mot du fichier et les stocke dans une liste
            in
            ( { model | wordList = words, isLoading = False }, Cmd.none ) --mise à jour de wordList avec words, isLoading est mis en False pour indiquer que le chargement est terminé
 
        WordsLoaded (Err _) ->  --cas où les mots n'ont pas pu être chargés ( erreur )
            ( { model | isLoading = False }, Cmd.none )
 
        NewRandomWordIndex index ->  --à partir d'un index aléatoire, cette fonction extrait le mot correspondant à l'index de la liste de mots du modèle
            let
                newWord = List.drop index model.wordList |> List.head |> Maybe.withDefault ""  --on récupère le mot aléatoire
            in
            ( { model | randomWord = newWord, guessInput = "", guessSuccess = False, gameState = InGame }  --mise à jour du modèle
            , Http.get { url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ newWord, expect = Http.expectJson DefinitionsLoaded definitionDecoder } --requete API pour récupérer le JSON contenant les définitions
            )
 
        UpdateGuess newGuess -> --dès que l'utilisateur écrit dans la zone de réponse, on vérifie si son entrée correspond au mot recherché
            let
                success = newGuess == model.randomWord --success est définie comme True si le mot est deviné
                newWordStatus =
                    { word = newGuess
                    , status = if success then Guessed else Skipped --un nouvel état de mot est créé avec la devinette et son statut (Guessed si elle correspond, Skipped sinon)
                    }
                newGuessedWords = if success then newWordStatus :: model.guessedWords else model.guessedWords --la liste de mots devinés est mise à jour en ajoutant le nouveau mot deviné. Cette liste sera affichée pendant et à la fin de la partie
            in
            ( { model | guessInput = newGuess, guessSuccess = success, score = if success then model.score + 1 else model.score, guessedWords = newGuessedWords }
            , Cmd.none )
 
        NewWord -> --branche appelée après le clic sur le bouton "New Word"
            let
                newWordStatus =
                    { word = model.randomWord
                    , status = Guessed --Le mot recherché est marqué comme deviné
                    }
 
                isWordAlreadyGuessed =
                    List.any (\w -> w.word == newWordStatus.word) model.guessedWords --on vérifie si le mot actuel a déjà été trouvé en le comparant avec les mots déjà trouvés dans le modèle
 
                newGuessedWords =
                    if isWordAlreadyGuessed then
                        model.guessedWords 
                    else
                        newWordStatus :: model.guessedWords --si il n'a pas enore été deviné, alors on l'ajoute dans la liste de mots devinés
            in
            ( { model | guessSuccess = False, guessedWords = newGuessedWords } --le succès est remis à False et la liste de mots devinés est mise à jour 
            , Random.generate NewRandomWordIndex (Random.int 0 (List.length model.wordList - 1)) --on appelle NewRandomIndex pour obtenir un nouveau mot
            )
 
        SkipWord -> --branche appelée après le clic sur le bouton "Skip Word"
            let
                skippedWord = { word = model.randomWord, status = Skipped } 
                newGuessedWords = skippedWord :: model.guessedWords --le mot est ajouté à la liste de mots trouvés avec le status Skipped
            in
            ( { model | score = model.score - 2, guessSuccess = False, guessedWords = newGuessedWords } --le score est décrementé de 2, le succès est mis à False et la liste de mots devinés est mise à jour
            , Random.generate NewRandomWordIndex (Random.int 0 (List.length model.wordList - 1))
            )
 
        StartGame -> --branche appelée lors du démarrage du jeu
            ( { model | gameState = InGame, score = 0, guessSuccess = False, timeElapsed = model.selectedTimer } --l'état du jeu est mis à jour à InGame, le score est réinitialisé à 0, le succès de la devinette est réinitialisé à False, et le temps écoulé est défini sur la valeur sélectionnée avec le minuteur
            , Random.generate NewRandomWordIndex (Random.int 0 (List.length model.wordList - 1)) --on appelle NewRandomWordIndex pour obtenir un nouveau mot
            )
 
        RestartGame -> --branche exécutée lorsque le jeu est redémarré
            ( { model | gameState = MainMenu, guessedWords = [] }, Cmd.none ) --On retourne au main menu
 
        DefinitionsLoaded (Ok defs) -> --lorsque les définitions d'un mot sont chargées avec succès, cette branche du motif de cas est exécutée
            ( { model | definitions = defs }, Cmd.none ) --ces définitions sont mises à jour dans le modèle
 
        DefinitionsLoaded (Err _) -> --branche exécutée lorsque une erreur se produit lors du chargement des définitions
            ( { model | definitions = [] }, Cmd.none ) --la liste des définitions est réinitialisée à vide
 
        Tick _ -> --cette branche est exécutée à chaque intervalle de l'horloge
            if model.gameState == InGame && model.timeElapsed > 0 then --si le jeu est en cours et qu'il reste du temps, le temps écoulé est mis à jour en décrémentant d'une unité
                let
                    newTime = model.timeElapsed - 1
                    newGameState = if newTime <= 0 then GameOver else InGame --si le temps est écoulé, l'état du jeu est mis à jour à GameOver
                in
                ( { model | timeElapsed = newTime, gameState = newGameState }, Cmd.none )
            else
                ( model, Cmd.none )
 
        UpdateSelectedTimer newTime -> --Mise à jour du temps quand l'utilisateur interagit avec la scroll bar
            ( { model | selectedTimer = newTime }, Cmd.none ) 
 


------------------------
----------VIEW----------
------------------------
 
view : Model -> Html Msg --cette fonction prend le modèle en entrée et retourne la page web en HTML. Elle permet l'affichage de la page web pour les différents états du jeu
view model =
    div
        [ style "text-align" "center", style "font-family" "'Roboto', sans-serif", style "padding" "20px" ]
        [ h1 [ style "font-size" "2.5em", style "margin-bottom" "20px" ] [ text "Guess it!" ]
        , case model.gameState of --on utilise case of pour gérer l'affichage en fonction des états du jeu
            MainMenu ->
                div --à l'intérieur de ce div on retrouve tous les éléments du menu principal
                    [ style "margin" "20px auto", style "width" "80%", style "max-width" "500px" ]
                    [ div [ style "margin-bottom" "10px" ] [ text "Set Timer (seconds): " ]
                    , input --étiquette pour définir le timer
                        [ type_ "range"
                        , Html.Attributes.min "20" --valeur min
                        , Html.Attributes.max "120" --valeur max
                        , value (String.fromInt model.selectedTimer)
                        , onInput (String.toInt >> Maybe.withDefault 5 >> UpdateSelectedTimer) --onInput permet de détecter les changements de la valeur de l'entrée. Lorsque la valeur change, la fonction UpdateSelectedTimer est appelée pour mettre à jour le modèle avec la nouvelle valeur de la minuterie
                        , style "width" "100%"
                        ]
                        []
                    , div [ style "margin-bottom" "20px" ] [ text (String.fromInt model.selectedTimer ++ " seconds") ]
                    , button --bouton start game
                        [ onClick StartGame, disabled model.isLoading, style "margin" "10px 0", style "padding" "10px 20px", style "font-size" "1em", style "cursor" "pointer" ] --onClick appelle la fonction StartGame. Si le jeu est en cours, le bouton est désactivé, ce qui est déterminé par la valeur de model.isLoading
                        [ text "Start Game" ]
                    ]
 
            InGame ->
            --à l'intérieur de ce div on retrouve tous les éléments affichés pendant que le jeu est en cours
            --cette mise en page se fait en 3 colonnes:
                div --cette première colonne affiche les définitions du mot actuel
                    [ style "display" "flex", style "justify-content" "center", style "align-items" "flex-start", style "margin" "20px auto", style "max-width" "800px" ]
                    [ div [ style "flex" "1", style "padding-right" "20px" ]
                        [ if List.isEmpty model.definitions then --si aucune définition est chargée, un texte "Loading deifnition..." est affiché
                            div [] [ text "Loading definition..." ]
                          else
                            ul [ style "text-align" "left" ] (List.concatMap viewDefinition model.definitions) --sinon, on crée une liste à puces dans laquelle sont affichées les définitions. On utilise List.concatMap pour appliquer la fonction viewDefinition à chaque élément de la liste model.definitions et concaténer les résultats
                        ]
                    , div --cette deuxième colonne affiche les mots devinés jusqu'à présent avec leur statut (guessed ou skipped)
                    [ style "flex" "1", style "border-left" "1px solid #ccc", style "padding-left" "20px" ] --style "flex" "1" nous permet de la positionner au centre, style "border-left" "1px solid #ccc" positionne la bordure gauche
                        [ h2 [] [ text "Guessed Words" ]
                        , ul [] --on crée une liste de puces pour les mots devinés (ou pas)
                            (List.map
                                (\wordStatus ->
                                    li (styleList wordStatus.status ) [ text wordStatus.word ] --les mots devinés sont dans une liste li avec une mise en forme différente en fonction de son statut (Guessed ou Skipped), définie par la fonction styleList
                                )
                                model.guessedWords
                            )
                        ]
                    , div [ style "flex" "1" ] --cette troisième colonne comprend un champ de saisie où l'utilisateur peut deviner le mot
                        [ input --ce input permet au joueur d'insérer le mot qui correspond à la définition
                            [ placeholder "Guess the word..."
                            , onInput UpdateGuess --onInput appelle la fonction UpdateGuess
                            , value model.guessInput --value affiche la valeur de l'entrée du joueur
                            , disabled model.guessSuccess 
                            , style "width" "100%", style "padding" "10px", style "margin-bottom" "10px"
                            ]
                            []
                        --si le mot a été trouvé, le bouton "New Word" s'affiche
                        , if model.guessSuccess then
                            div []
                            [ button [ onClick NewWord, style "padding" "10px 20px", style "font-size" "1em", style "cursor" "pointer" ] [ text "New Word" ]
                            , p [] [ text "Congratulations, you guessed the word!" ]
                            ]
                          else --sinon c'est le bouton "Skip Word"
                            button [ onClick SkipWord, style "padding" "10px 20px", style "font-size" "1em", style "cursor" "pointer" ] [ text "Skip Word" ]
                        , p [] [ text ("Time: " ++ String.fromInt model.timeElapsed ++ "s") ] --le temps restant est affiché
                        , p [] [ text ("Score: " ++ String.fromInt model.score) ] --le score est affiché
                        ]
                    ]
 
            GameOver ->
                div --à l'intérieur de ce div on retrouve tous les éléments lorsque le jeu est fini
                    [ style "margin" "20px auto", style "width" "80%", style "max-width" "500px" ]
                    [ h2 [] [ text "Game Over" ] --texte "Game Over"
                    , button --bouton pour restart
                        [ onClick RestartGame, style "margin" "10px 0", style "padding" "10px 20px", style "font-size" "1em", style "cursor" "pointer" ] --onClick est utilisé pour appeler la fonctionRestartGame lorsque le bouton est appuyé
                        [ text "Restart" ]
                    , h2 [] [ text ("Score: " ++ String.fromInt model.score) ] --texte affichant le score
                    , div [ style "margin-bottom" "20px" ]
                        [ h3 [] [ text "Guessed Words:" ]
                        --on crée une liste de puces et utilise la fonction List.map pour parcourir chaque élément de la liste model.guessedWords, qui contient toutes les structures WordStatus avec les mots devinés et leur état
                        , ul [] (List.map (\wordStatus -> li (styleList wordStatus.status) [ text wordStatus.word ]) model.guessedWords) --affiche les mots devinés et passés pendant la dernière partie
 
                        ]
                    ]
        ]

-- viewDefinition, viewMeaning et viewDefinitionDetail : gèrent le bon affichage des définitions
viewDefinition : WordDefinition -> List (Html Msg)
viewDefinition def =
    [ -- Pour se donner un indice : ul [] [text (def.word) ] Affiche le mot recherché ;)
      ul [] (List.concatMap viewMeaning def.meanings) --List.concatMap applique la fonction viewMeaning à chaque élément de la liste des significations (def.meanings) et concatène les résultats
    ]
 
viewMeaning : Meaning -> List (Html Msg) --cette fonction prend une Meaning en entrée et crée une liste contenant la partie du texte (meaning.partOfSpeech) de la signification, puis une liste non ordonnée contenant les détails de la définition (definitionDetail)
viewMeaning meaning =
    [ li [] [ text meaning.partOfSpeech ]
    , ul [] (List.map viewDefinitionDetail meaning.definitions) --on applique avec List.map la fonction viewDefinitionDetail aux éléments de la liste meaning.definitions
    ]
 
viewDefinitionDetail : DefinitionDetail -> Html Msg --cette fonction, appelée dans viewMeaning, prend un DefinitionDetail en entrée et retourne un message HTML. Elle crée une liste contenant la définition (detail.definition). Avec viewMeaning on va donc créer une liste de ces messages HTML pour chaque meaning du mot
viewDefinitionDetail detail =
    li [] [ text detail.definition ]
 


------------------------
-----SUBSCRIPTIONS------
------------------------
 
subscriptions : Model -> Sub Msg --cette fonction détermine quelles souscriptions doivent être actives en fonction de l'état actuel du jeu 
subscriptions model =
    case model.gameState of
        InGame -> Time.every 1000 Tick --si le jeu est en cours, on crée une souscription pour déclencher le messafe Tick toutes les secondes (every 1000)
        _ -> Sub.none --sinon, aucune souscription n'est crée
 

 
------------------------
----------MAIN----------
------------------------
 
main = --fonction qui initialise l'application
    Browser.element
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view
        }



------------------------
--DECODAGE DU JSON API--
------------------------
 
type alias WordDefinition =
    { word : String
    , meanings : List Meaning
    }
 
type alias Meaning =
    { partOfSpeech : String
    , definitions : List DefinitionDetail
    }
 
type alias DefinitionDetail =
    { definition : String }
 
--décodeurs de JSON, utilisés pour transformer les données JSON en structures de données Elm utilisables dans l'application
 
definitionDecoder : Decoder (List WordDefinition) --ce décodeur est utilisé pour décoder une liste de WordDefinition à partir d'un JSON. Un WordDefinition est une structure de données Elm qui contient un mot (word) et une liste de définitions (meanings), où chaque définition contient le type de mot (partOfSpeech) et une liste de détails de définitions (definitions). Se réferrer au API pour comprendre la mise en forme du dictionnaire fourni
 
definitionDecoder =
    list
        (map2 WordDefinition --map2 combine deux valeurs décodées avec une fonction, ici WordDefinition
            (field "word" string)
            (field "meanings" (list meaningDecoder))
        )
 
meaningDecoder : Decoder Meaning --ce décodeur est utilisé pour décoder une Meaning à partir d'un JSON
meaningDecoder =
    map2 Meaning --on combine les 2 valeurs décodées avec la fonction Meaning
        (field "partOfSpeech" string)
        (field "definitions" (list definitionDetailDecoder))
 
definitionDetailDecoder : Decoder DefinitionDetail --ce décodeur est utilisé pour décoder un DefinitionDetail à partir d'un JSON. Un DefinitionDetail contient les définitions des mots
definitionDetailDecoder =
    map DefinitionDetail (field "definition" string) --on applique la fonction DefinitionDetail à la chaîne décodée représentant la définition


