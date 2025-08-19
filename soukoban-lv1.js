//readlineモジュールをインポート
const readline = require("readline");

// 盤面を作るStageクラスを定義
class Stage {
    constructor() {
        this.map = [
            '#####',
            '#.o.#',
            '# @ #',
            '#.o.#',
            '#####',
        ];

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
        const playerRow = this.map[playerY]; // playerRowはプレイヤーのいる行
        this.map[playerY] = playerRow.substring(0 , playerX) + " " + playerRow.substring(playerX + 1);

    }

    // 盤面全体を表示するメソッド
    display () {
        console.clear();
        const player = this.player;
        // 元のマップをコピーする
        const viewMap = this.map.slice();
        // プレイヤーがいる行を文字列から配列に変換
        const playerRowArray = viewMap[player.y].split("");
        // プレイヤーの位置に"@"を置く
        playerRowArray[player.x] = "@";
        // 文字列に戻してマップに反映
        viewMap[player.y] = playerRowArray.join("");

        viewMap.forEach(row => {
            console.log(row);
        })
    }
    // クリア判定を行うメソッド
    // 盤面上の荷物がすべてゴールに置かれているかチェック
    isClear() {
    }
}

//プレイヤークラスを定義
class Player {
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }
    //プレイヤーの移動に関するメソッドを定義
    move(dx, dy, stage) {
        // 移動先の座標を計算
        const nextX = this.x + dx;
        const nextY = this.y + dy;
        // 移動先が壁なら何もしない
        if (stage.map[nextY][nextX] === "#") {
            return;
        }
        // 移動先が荷物なら、荷物の一個先をチェック
        if (stage.map[nextY][nextX] === "o") {
            const boxNextX = nextX + dx;
            const boxNextY = nextY + dy;
            //荷物の一個先が壁か荷物なら何もしない
            //早期リターンってやつ
            if (stage.map[boxNextY][boxNextX] === "#" || stage.map[boxNextY][boxNextX] === "o") {
                return;
            }
            // returnしなかったら荷物を移動する
            // 荷物のある行を文字列から配列に変換し、荷物があった場所を空白にしてからもう一度文字列に変換する
            const boxRow = stage.map[nextY].split("");
            boxRow[nextX] = " ";
            stage.map[nextY] = boxRow.join("");

            // 荷物の移動先の行を文字列から配列に変換し、荷物の移動先を荷物にしてからもう一度文字列に変換する
            const boxNextRow = stage.map[boxNextY].split("");
            boxNextRow[boxNextX] = "o";
            stage.map[boxNextY] = boxNextRow.join("");
                    
        }
        this.x = nextX;
        this.y = nextY;
    }
}

// ゲームクラスを定義
class Game {
    constructor() {
        this.stage = new Stage();
        this.setupInput();
    }
    setupInput() {
        readline.emitKeypressEvents(process.stdin);
        process.stdin.setRawMode(true);

        process.stdin.on('keypress', (str, key) => {
        // Ctrl+Cが押されたらプログラムを終了する
        if (key.ctrl && key.name === 'c') {
            process.exit();
        }

        // TODO: ここでキーに応じた移動処理を行う
        const player = this.stage.player;

        //入力の分岐によって座標を変更
        if (key.name === "w") {
            player.move(0, -1 ,this.stage);
        } else if (key.name === "a") {
            player.move(-1, 0, this.stage);
        } else if (key.name === "s") {
            player.move(0, 1 , this.stage);
        } else if (key.name === "d") {
            player.move(1, 0, this.stage);
        }
        // 毎回の入力後に、必ず盤面を再描画する
        this.stage.display();
        });
    }

    start() {
        this.stage.display(); // Stageクラスのdisplayメソッドを呼び出して盤面を表示        
    }
}

const game = new Game();
game.start();

