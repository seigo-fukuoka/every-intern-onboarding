//hit and blow game

//1. readlineモジュールをインポート
const readline = require("readline");

//2. ゲームクラスを定義
class HitAndBlow {
    constructor() {
        this.rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout
        });
        this.DIGITS = 4;

        // 1. 0から9までの配列（山札）を用意
        const source = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
        const answer = [];

        // 2. 山札から4枚のカードをランダムに引く
        for (let i = 0; i < this.DIGITS; i++) {
            // 山札からランダムな位置を選ぶ
            const randIndex = Math.floor(Math.random() * source.length);
            // その位置の数字を1枚抜き取り、山札から削除する
            const pickedNumber = source.splice(randIndex, 1)[0];
            // 抜き取った数字を答えの配列に追加する
            answer.push(pickedNumber);
        }

        this.secretArr = answer;
    }

    //4. ゲーム開始メソッド
    gameStart() {
        console.log("=================================");
        console.log("Hit & Blowゲームを開始します");
        console.log("=================================");
        this.judgeAnswer();
    }

    //5. ユーザーの入力を受け付けるメソッド
    judgeAnswer() {
        this.rl.question(`${this.DIGITS}桁の数字を入力してください: `, (input) => {
            const inputArray = input.split('').map(Number);
            
            // 入力値のチェック
            if (
                input.length !== this.DIGITS ||
                inputArray.some(isNaN) ||
                new Set(inputArray).size !== this.DIGITS
            ) {
                console.log(`無効な入力です。${this.DIGITS}桁の重複しない数字を入力してください`);
                this.judgeAnswer();
                return;
            }

            console.log("入力されたのは" + input + "ですね");
            let hit = 0;
            let blow = 0;

            //6. hitとblowを計算
            for (let i = 0; i < inputArray.length; i++) {
                for (let j = 0; j < this.secretArr.length; j++) {
                    //同じ数値が無ければcontinueする
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
            
            //7. hitが4ならゲーム終了
            if (hit === this.DIGITS) {
                console.log("おめでとうございます！正解です！");
                this.rl.close();
            } else {
                console.log("もう一度入力してください");
                this.judgeAnswer();
            }
        });
    }

    //8. ゲーム開始メソッドを呼び出す
    start() {
        this.gameStart();
    }
}

//9. ゲーム開始
const hitAndBlow = new HitAndBlow();
hitAndBlow.start();