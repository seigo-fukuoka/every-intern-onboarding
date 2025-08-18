//hit and blow game

//1. readlineモジュールをインポート
const readline = require("readline");

//2. ゲームクラスを定義
class HitAndBlow {
    constructor(){
        this.rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout
        });
        //3. 4桁の数字をランダムに生成
        const answer = new Set();
        while ( answer.size < 4){
            const randomNumber = Math.floor(Math.random() *10);
            answer.add(randomNumber);
        }
        // Setを配列に変換
        this.secretArr = Array.from(answer);
    }
    
    //4. ゲーム開始メソッド
    gameStart(){
        console.log("=================================");
        console.log("Hit & Blowゲームを開始します");
        console.log("=================================");
        this.judgeAnswer();
    }

    //5. ユーザーの入力を受け付けるメソッド
    judgeAnswer(){
        this.rl.question(`4桁の数字を入力してください: `, (input) => {
            console.log("入力されたのは" + input + "ですね");
            const inputArray = input.split('').map(Number);
            let hit = 0;
            let blow = 0;
            //6. hitとblowを計算
            for ( let i = 0; i < inputArray.length; i ++){
                for ( let j = 0; j < this.secretArr.length; j ++){
                    //同じ数値があれば少なくともblowが1増える
                    if (inputArray[i] === this.secretArr[j]){
                        // インデックスまで同じならhitが1増える
                        if (i === j){
                            hit++;
                        }else{
                            blow++;
                        }
                    }
                }
            }
            console.log("Hit: " + hit +  " Blow " + blow);
            //7. hitが4ならゲーム終了
            if (hit === 4){
                console.log("おめでとうございます！正解です！");
                this.rl.close();
            }else{
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
const echoGame = new HitAndBlow();
echoGame.start();
