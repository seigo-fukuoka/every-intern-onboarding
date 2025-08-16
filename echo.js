#!/usr/bin/env node

const readline = require('readline');

// ゲームクラス
class Echo {
    constructor() {
        this.rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
        
        });
const answer = new Set();
while (answer.size < 4){
    const randomNumber = Math.floor(Math.random() * 10);
    answer.add(randomNumber); 
}
        this.secret = Array.from(answer);
    }

  // ゲーム開始
    start() {
        console.log('=================================');
        console.log('Hit & Blowゲームを開始します');
        console.log('=================================');

        this.promptEcho();
    }

    promptEcho() {
        this.rl.question(`4桁の数字を入力してください: `, (input) => {
            console.log("入力されたのは: " +input + " ですね \n");
            const inputArray = input.split('').map(Number);
            let hit = 0;
            let blow = 0;
            for (let i = 0; i < inputArray.length ; i ++){
                for (let j = 0; j < this.secret.length; j++){
                    if ( inputArray[i] === this.secret[j] ){
                        if ( i === j ){
                            hit ++;
                        }else{
                            blow ++;
                        }
                        }
                    }

                }
        console.log ( "Hit:" + hit + "Blow:" + blow);
            if (hit === 4) {
                console.log("おめでとうございます！正解です！");
                this.rl.close();
            } else {
                console.log("もう一度入力してください。\n");
                this.promptEcho();
            }
                });
            }
        }

// ゲーム開始
const game = new Echo();
game.start();
