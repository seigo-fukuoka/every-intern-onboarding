//readlineモジュールをインポート
const readline = require("readline");

// Stageクラスを定義（盤面全体の管理者）
// 役割：盤面の状態（壁、荷物、プレイヤーの位置）をすべて把握し、ゲームのルールを実行する責任者
class Stage {
    constructor() {
        // 盤面を定義し、配列として保持しておく
        this.mapData = [
            "#########",
            "#.   o  #",
            "#   @   #",
            "# o   . #",
            "#       #",
            "#       #",
            "#########",     
        ];
        this.map = this.mapData.map(row => row.split("")); // 文字列を配列に変換
        
        // ゴールの場所を(X,Y)座標で把握する、
        this.goalPositions = [];
        this.map.forEach((row, y) => {
            let index = -1;
            while ((index = row.indexOf('.', index + 1)) !== -1) {
                this.goalPositions.push({ x: index, y: y });
        }})
        // プレイヤーの初期位置を設定
        let playerX;
        let playerY;

        this.map.forEach((row, y) => {
            const x = row.indexOf("@");
            if (x !== -1) { //indexOfは見つからない場合-1を返す
                playerX = x; // プレイヤーのX座標
                playerY = y; // プレイヤーのY座標
            }
        });

        // 見つけた座標でPlayerのインスタンスの生成
        this.player = new Player(playerX, playerY);

        // プレイヤーの位置から"@"を削除"
        this.map[playerY][playerX] = " ";

    }
    // Playerを移動させるメソッド
    movePlayer(dx, dy) {
        // 移動先の座標を計算
        const nextX = this.player.x + dx;
        const nextY = this.player.y + dy;
        // 移動先が壁なら何もしない
        if (this.map[nextY][nextX] === "#") {
            return;
        }
        // 移動先が荷物なら、荷物の一個先をチェック
        if (this.map[nextY][nextX] === "o") {
            const boxNextX = nextX + dx;
            const boxNextY = nextY + dy;
            //荷物の一個先が壁か荷物なら何もしない
            //早期リターンってやつ
            if (this.map[boxNextY][boxNextX] === "#" || this.map[boxNextY][boxNextX] === "o") {
                return;
            } 
            // returnしなかったら荷物を移動する
            this.map[nextY][nextX] = " "; // 荷物を移動した場所を空白にする

            // 荷物の移動先の行を文字列から配列に変換し、荷物の移動先を荷物にしてからもう一度文字列に変換する
            this.map[boxNextY][boxNextX] = "o";
                    
        }

        this.player.x = nextX;
        this.player.y = nextY;
    }

    // 盤面全体を表示するメソッド
    display () {
        console.clear();
        console.log("=================================");
        console.log("Sokoban Level 1");
        console.log("=================================");
        console.log("操作方法: w(上), a(左), s(下), d(右), r(リセット), q(終了)");
        console.log("=================================");
        const player = this.player;
        // 元のマップをコピーする、普通にコピーすると浅いコピーになってしまい、元のマップに影響が出るらしい
        const viewMap = JSON.parse(JSON.stringify(this.map));
        // ゴール位置を表示する（荷物がない場合）
        this.goalPositions.forEach(goal => {
            if (this.map[goal.y][goal.x] === " ") {
                viewMap[goal.y][goal.x] = ".";
            }
        });
        // プレイヤーがいる行を文字列から配列に変換
        viewMap[this.player.y][this.player.x] = "@";      // プレイヤーの位置に"@"を置く
        // 文字列に戻してマップに反映
        viewMap.forEach(rowArray => {
            console.log(rowArray.join(""));
        })
    }
    // クリア判定を行うメソッド
    // 盤面上のゴールの座標を把握しておき、すべての座標に荷物が置かれているかチェックする
    // GameクラスのisClearメソッドから呼び出される
    isClear() {
        for (let i = 0; i < this.goalPositions.length; i++) {
            if (this.map[this.goalPositions[i].y][this.goalPositions[i].x] !== "o") {
                return false;
            }
        }
        return true;
    }
}

//プレイヤークラスを定義（プレイヤーの座標のみを管理する駒）
// 役割：
class Player {
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }
}

// ゲームクラスを定義（ゲーム全体の司令塔、支配人）
// 役割：ユーザーからのキー入力を受付、それをStageクラスへの命令に変換する司令塔
class Game {
    constructor() {
        this.stage = new Stage();
        this.setupInput();
    }
    // ユーザーからの入力を受け付けるメソッド
    setupInput() {
        readline.emitKeypressEvents(process.stdin);
        process.stdin.setRawMode(true);

        process.stdin.on('keypress', (str, key) => {
        // Qが押されたらプログラムを終了する
        if (key.name === "q") {
            process.exit();
        }

        if (key.name === "r") {
            this.reset();
            return;
        }

        //入力の分岐によって座標を変更
        if (key.name === "w") {
            this.stage.movePlayer(0, -1);
        } else if (key.name === "a") {
            this.stage.movePlayer(-1, 0);
        } else if (key.name === "s") {
            this.stage.movePlayer(0, 1);
        } else if (key.name === "d") {
            this.stage.movePlayer(1, 0);
        }
        // 毎回の入力後に、必ず盤面を再描画する
        this.stage.display();

        if (this.stage.isClear()) {
        console.log('クリアおめでとう！');
        process.exit(); // ゲームを終了する
        }
        });
    }
    // ゲームを開始するメソッド
    start() {
        this.stage.display(); // Stageクラスのdisplayメソッドを呼び出して盤面を表示        
    }
    // ゲームをリセットするメソッド
    reset() {
        this.stage = new Stage(); // 新しいStageインスタンスを生成
        this.stage.display(); // 盤面を再表示
    }
}

const game = new Game();
game.start();

