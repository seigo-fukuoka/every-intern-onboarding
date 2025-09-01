//hit and blow game

//readlineモジュールをインポート
const readline = require("readline");

//ゲームクラスを定義
class HitAndBlow {
    constructor() {
        this.tryCount = 0; // 追加: 試行回数をカウントする変数
        this.secretArr = [];
        this.DIGITS = 0;
        this.history = []; // 追加: 回答履歴を保存する配列
        this.rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout
        });
    }
    // 難易度選択メソッド
    selectDifficulty() {
        this.rl.question(
            "難易度を選択してください:\n" +
            "簡単(3桁)→「1」を入力してください\n" +
            "普通(4桁)→「2」を入力してください\n" +
            "難しい(5桁)→「3」を入力してください\n",
            (difficulty) => {
                if (difficulty === "1") {
                    this.DIGITS = 3;
                } else if (difficulty === "2") {
                    this.DIGITS = 4;
                } else if (difficulty === "3") {
                    this.DIGITS = 5;
                } else {
                    console.log("無効な入力です。再度選択してください。");
                    this.selectDifficulty();
                    return;
                }
            console.log(`それでは、${this.DIGITS}桁のHit & Blowゲームを開始します`);
            this.generateSecretAnswer();
        });
    }


    generateSecretAnswer() {
        // 0から9までの配列（山札）を用意
        const source = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
        const answer = [];

        // 山札から4枚のカードをランダムに引く
        for (let i = 0; i < this.DIGITS; i++) {
            const randIndex = Math.floor(Math.random() * source.length);
            const pickedNumber = source.splice(randIndex, 1)[0];
            answer.push(pickedNumber);
        }

        this.secretArr = answer;
        this.gameStart();
    }

    // ゲーム開始メソッド
    gameStart() {
        console.log("=================================");
        console.log("Hit & Blowゲームを開始します");
        console.log("=================================");
        this.judgeAnswer();
    }

    // ユーザーの入力を受け付けるメソッド
    judgeAnswer() {
        // 追加: 現在の試行履歴を表示
        console.log("--- これまでの履歴 ---");
        this.history.forEach((record , index) => {
            console.log(`${index + 1} 回目: ユーザーの推測 - ${record.Guess} | hit - ${record.hit} | blow - ${record.blow}`);
        });
        console.log("--------------------");
        this.rl.question(`${this.DIGITS}桁の数字を入力してください: `, (input) => {
            const inputArray = input.split('').map(Number);
            
            // 入力値のチェック
            if (
                input.length !== this.DIGITS ||
                inputArray.some(isNaN) 
            ) {
                console.log(`無効な入力です。${this.DIGITS}桁の重複しない数字を入力してください`);
                this.judgeAnswer();
                return;
            }

            this.tryCount++; // 追加: 試行回数をカウント
            console.log("入力されたのは" + input + "ですね");
            let hit = 0;
            let blow = 0;

            // hitとblowを計算
            for (let i = 0; i < inputArray.length; i++) {
                for (let j = 0; j < this.secretArr.length; j++) {
                    // 同じ数値が無ければcontinueする
                    if (inputArray[i] !== this.secretArr[j]) {
                        continue;
                    }
                    if (i === j) {
                        hit++;
                    } else {
                        blow++;
                    }
                    break;
                }
            }

            console.log("Hit: " + hit + " Blow: " + blow);
            console.log("--------------------");
            console.log(`試行回数: ${this.tryCount}`); // 追加: 試行回数を表示
            this.history.push({
                Guess: input,
                hit: hit,
                blow: blow,
            })
            
            // hitが4ならゲーム終了
            if (hit === this.DIGITS) {
                console.log("おめでとうございます！正解です！");
                this.askToPlayAgain();
            } else {
                console.log("もう一度入力してください");
                this.judgeAnswer();
            }
        });
    }
    // ユーザーに再プレイを尋ねるメソッド
    askToPlayAgain() {
        this.rl.question("もう一度プレイしますか？ (y/n): ", (answer) => {
            if (answer.toLowerCase() === "y") {
                console.log("ゲームをリセットします...");
                this.resetGame();
            } else if (answer.toLowerCase() === "n") {
                console.log("ゲームを終了します。ありがとうございました。");
                this.rl.close();         
            } else {
                console.log("無効な入力です。再度入力してください。");
                this.askToPlayAgain();
            }
        });
    }
    // ゲームをリセットするメソッド
    resetGame() {
        this.tryCount = 0;
        this.secretArr = [];
        this.history = [];
        this.DIGITS = 0;
        this.start();
    }
    start() {
        this.selectDifficulty();
    }
}

// ゲーム開始
const hitAndBlow = new HitAndBlow();
hitAndBlow.start();