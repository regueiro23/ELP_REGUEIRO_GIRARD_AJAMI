-- Guess it!

module Main exposing (..)

import Browser
import Html exposing (Html, div, h1, h2, input, text, button, ul, li, p)
import Html.Attributes exposing (style, type_, value, placeholder, disabled)
import Html.Events exposing (onInput, onClick)
import String
import Random
import Http
import Time
import Json.Decode exposing (Decoder, field, list, string, map, map2)

-- MODEL

type GameState
    = MainMenu
    | InGame
    | GameOver

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
    }

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
    }

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

-- INIT

init : () -> (Model, Cmd Msg)
init _ =
    ( initialModel
    , Http.get
        { url = "http://localhost:5017/static/mots.txt"
        , expect = Http.expectString WordsLoaded
        }
    )

-- UPDATE

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        WordsLoaded (Ok data) ->
            let
                words = String.split " " data
            in
            ( { model | wordList = words, isLoading = False }, Cmd.none )

        WordsLoaded (Err _) ->
            ( { model | isLoading = False }, Cmd.none )

        NewRandomWordIndex index ->
            let
                newWord = List.drop index model.wordList |> List.head |> Maybe.withDefault ""
            in
            ( { model | randomWord = newWord, guessInput = "", guessSuccess = False, gameState = InGame }
            , Http.get { url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ newWord, expect = Http.expectJson DefinitionsLoaded definitionDecoder }
            )

        UpdateGuess newGuess ->
            let
                success = newGuess == model.randomWord
            in
            ( { model | guessInput = newGuess, guessSuccess = success, score = if success then model.score + 1 else model.score }, Cmd.none )

        NewWord ->
            ( { model | guessSuccess = False}
            , Random.generate NewRandomWordIndex (Random.int 0 (List.length model.wordList - 1))
            )

        SkipWord ->
            ( { model | score = model.score - 2, guessSuccess = False }
            , Random.generate NewRandomWordIndex (Random.int 0 (List.length model.wordList - 1))
            )

        StartGame ->
            ( { model | gameState = InGame, score = 0, guessSuccess = False, timeElapsed = model.selectedTimer }
            , Random.generate NewRandomWordIndex (Random.int 0 (List.length model.wordList - 1))
            )

        RestartGame ->
            ( { model | gameState = MainMenu }, Cmd.none )

        DefinitionsLoaded (Ok defs) ->
            ( { model | definitions = defs }, Cmd.none )

        DefinitionsLoaded (Err _) ->
            ( { model | definitions = [] }, Cmd.none )

        Tick _ ->
            if model.gameState == InGame && model.timeElapsed > 0 then
                let
                    newTime = model.timeElapsed - 1
                    newGameState = if newTime <= 0 then GameOver else InGame -- Passer à GameOver si le temps est écoulé
                in
                ( { model | timeElapsed = newTime, gameState = newGameState }, Cmd.none )
            else
                ( model, Cmd.none )

        UpdateSelectedTimer newTime ->
            ( { model | selectedTimer = newTime }, Cmd.none )

-- VIEW

view : Model -> Html Msg
view model =
    div
        [ style "text-align" "center", style "font-family" "'Roboto', sans-serif", style "padding" "20px" ]
        [ h1 [ style "font-size" "2.5em", style "margin-bottom" "20px" ] [ text "Guess it!" ]
        , case model.gameState of
            MainMenu ->
                div
                    [ style "margin" "20px auto", style "width" "80%", style "max-width" "500px" ]
                    [ div [ style "margin-bottom" "10px" ] [ text "Set Timer (seconds): " ]
                    , input
                        [ type_ "range"
                        , Html.Attributes.min "20"
                        , Html.Attributes.max "120"
                        , value (String.fromInt model.selectedTimer)
                        , onInput (String.toInt >> Maybe.withDefault 5 >> UpdateSelectedTimer)
                        , style "width" "100%"
                        ]
                        []
                    , div [ style "margin-bottom" "20px" ] [ text (String.fromInt model.selectedTimer ++ " seconds") ]
                    , button
                        [ onClick StartGame, disabled model.isLoading, style "margin" "10px 0", style "padding" "10px 20px", style "font-size" "1em", style "cursor" "pointer" ]
                        [ text "Start Game" ]
                    ]

            InGame ->
                div
                    [ style "display" "flex", style "justify-content" "center", style "align-items" "flex-start", style "margin" "20px auto", style "max-width" "800px" ]
                    [ div [ style "flex" "1", style "padding-right" "20px" ]
                        [ if List.isEmpty model.definitions then
                            div [] [ text "Loading definition..." ]
                          else
                            ul [ style "text-align" "left" ] (List.concatMap viewDefinition model.definitions)
                        ]
                    , div [ style "flex" "1" ]
                        [ input
                            [ placeholder "Guess the word..."
                            , onInput UpdateGuess
                            , value model.guessInput
                            , disabled model.guessSuccess
                            , style "width" "100%", style "padding" "10px", style "margin-bottom" "10px"
                            ]
                            []
                        , if model.guessSuccess then
                            div []
                            [ button [ onClick NewWord, style "padding" "10px 20px", style "font-size" "1em", style "cursor" "pointer" ] [ text "New Word" ]
                            , p [] [ text "Congratulations, you guessed the word!" ]
                            ]
                          else
                            button [ onClick SkipWord, style "padding" "10px 20px", style "font-size" "1em", style "cursor" "pointer" ] [ text "Skip Word" ]
                        , p [] [ text ("Time: " ++ String.fromInt model.timeElapsed ++ "s") ]
                        , p [] [ text ("Score: " ++ String.fromInt model.score) ]
                        ]
                    ]

            GameOver ->
                div
                    [ style "margin" "20px auto", style "width" "80%", style "max-width" "500px" ]
                    [ h2 [] [ text "Game Over" ]
                    , button
                        [ onClick RestartGame, style "margin" "10px 0", style "padding" "10px 20px", style "font-size" "1em", style "cursor" "pointer" ]
                        [ text "Restart" ]
                    , h2 [] [ text ("Score: " ++ String.fromInt model.score) ]
                    ]
        ]

viewDefinition : WordDefinition -> List (Html Msg)
viewDefinition def =
    [ -- Pour se donner un indice : ul [] [text (def.word) ]
      ul [] (List.concatMap viewMeaning def.meanings)
    ]

viewMeaning : Meaning -> List (Html Msg)
viewMeaning meaning =
    [ li [] [ text meaning.partOfSpeech ]
    , ul [] (List.map viewDefinitionDetail meaning.definitions)
    ]

viewDefinitionDetail : DefinitionDetail -> Html Msg
viewDefinitionDetail detail =
    li [] [ text detail.definition ]

-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions model =
    case model.gameState of
        InGame -> Time.every 1000 Tick
        _ -> Sub.none

main =
    Browser.element
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view
        }

-- JSON

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

definitionDecoder : Decoder (List WordDefinition)
definitionDecoder =
    list
        (map2 WordDefinition
            (field "word" string)
            (field "meanings" (list meaningDecoder))
        )

meaningDecoder : Decoder Meaning
meaningDecoder =
    map2 Meaning
        (field "partOfSpeech" string)
        (field "definitions" (list definitionDetailDecoder))

definitionDetailDecoder : Decoder DefinitionDetail
definitionDetailDecoder =
    map DefinitionDetail (field "definition" string)
