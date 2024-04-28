let quiz = {};
let questionIdx = -1;
let correctGuesses = -1;

function startGame(){
    let game = document.getElementById("game");
    let gameEndScreen = document.getElementById("gameEndScreen");
    game.style.display = "inherit";
    gameEndScreen.style.display = "none";

    questionIdx = 0;
    correctGuesses = 0;
    let gameTitle = document.getElementById("gameTitle");
    gameTitle.innerHTML = quiz.title;
    setupQuestion()
}

function setupQuestion(){
    let answers = document.getElementById("answers");
    let qTitle = document.getElementById("gameQuestionTitle");
    
    if (questionIdx < quiz.questions.length) {
        qTitle.innerHTML = quiz.questions[questionIdx].title;
        answers.innerHTML = '';
        let i = 0;
        for (const answer of quiz.questions[questionIdx].choices) {
            let answerButton = document.createElement('button');
            answers.appendChild(answerButton);
            answerButton.innerHTML = answer.title;
            answerButton.onclick = ( (x) => () => {answerQuestion(x);} )(i);
            i++;
        }
    } else {
        endGame();
    }
}

function answerQuestion(idx){
    let i = parseInt(idx);
    if (quiz.questions[questionIdx].choices[i].correct) {
        correctGuesses++;
    }
    setupQuestion(++questionIdx);
}

function endGame(){
    let game = document.getElementById("game");
    let gameEndScreen = document.getElementById("gameEndScreen");
    let gameEndScore = document.getElementById("gameEndScore");
    gameEndScore.innerHTML = `You got ${correctGuesses} of ${quiz.questions.length} correct.`;
    game.style.display = "none";
    gameEndScreen.style.display = "block";

    try{
        fetch("gameEnd");
    }catch(ex){
        
    }
}

async function loadGame(title){
    let response = await fetch("game?" + new URLSearchParams({
        title:title
    }));
    let gameQuiz = await response.json();
    quiz = gameQuiz;
    questionIdx = 0;
}

function prepAddGame(){
    let quiz = quizEntryToJSON();
    let str = JSON.stringify(quiz);
    document.getElementById("addQuizJSON").value = str;
}

function quizEntryToJSON(){
   let quizTitle = document.getElementById("addQuizTitle").value;
    let questionDIVs = document.getElementsByClassName("questionEntryCard");

    let quiz = {
        title: quizTitle,
        questions: []
    };

    for (const q of questionDIVs) {
        let questionTitle = q.querySelector('input[type="text"]').value;
        let question = {
            title: questionTitle,
            choices: []
        };
        let answerDIVs = q.getElementsByClassName("addQuestionAnswer");
        for (const a of answerDIVs) {

            let text = a.querySelector('input[type="text"]').value;
            let isCorrect = a.querySelector('input[type="checkbox"]').checked;
            let answer = {
                title: text,
                correct: isCorrect
            }
            question.choices.push(answer);
        }
        quiz.questions.push(question);
    }
    return quiz;
}