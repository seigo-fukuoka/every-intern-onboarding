#!/usr/bin/env node

const readline = require('readline');

// ゲームクラス
class Echo {
    constructor() {
        this.rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
        });
    }

  // ゲーム開始
    start() {
        console.log('=================================');
        console.log('ユーザーの入力を繰り返します');
        console.log('=================================');

        this.promptEcho();
    }

    promptEcho() {
        this.rl.question(`好きな文字を入力してください: `, (input) => {
            console.log("入力されたのは: " +input + " ですね \n");
            this.promptEcho();
        });
    }
}

// ゲーム開始
const game = new Echo();
game.start();
