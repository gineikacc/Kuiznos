function startGame(){
    


let quiz = {
    "id": 101,
    "title": "Capitals",
    "author": "One Guy",
    "questions": [
        {"id": 101,
        "title": "Italy",
        "choices": [
            {
                "id": 101,
                "title": "Rome",
                "correct": true
            },
            {
                "id": 101,
                "title": "London",
                "correct": false
            },
            {
                "id": 101,
                "title": "Athens",
                "correct": false
            },
            {
                "id": 101,
                "title": "Oslo",
                "correct": false
            }
        ]
    } ,
    {"id": 101,
        "title": "Sweden",
        "choices": [
            {
                "id": 101,
                "title": "Rome",
                "correct": false
            },
            {
                "id": 101,
                "title": "London",
                "correct": false
            },
            {
                "id": 31,
                "title": "Athens",
                "correct": false
            },
            {
                "id": 101,
                "title": "Oslo",
                "correct": true
            }
        ]
    }
    ]
};


}

// type Quiz struct {
// 	Id        int               `form:"id" json:"id"`
// 	Title     string            `form:"title" json:"title"`
// 	Author    string            `form:"author" json:"author"`
// 	Questions []*Question_entry `form:"questions" json:"questions"`
// }

// type Question_entry struct {
// 	Id      int                `form:"id" json:"id"`
// 	Title   string             `form:"title" json:"title"`
// 	Choices []*Question_choice `form:"choices" json:"choices"`
// }

// type Question_choice struct {
// 	Id         int    `form:"id" json:"id"`
// 	Title      string `form:"title" json:"title"`
// 	Is_correct bool   `form:"correct" json:"correct"`

