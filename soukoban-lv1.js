//readlineモジュールをインポート
const readline = require("readline");

// 盤面全体を管理するStageクラスを定義
class Stage {
    constructor() {
        this.map = [
            '#####',
            '#.o #',
            '# @ #',
            '# o.#',
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

        // ゴールの場所を(X,Y)座標で把握する
        this.goalPositions = [];
        this.map.forEach((row, y) => {
            let index = -1;
            while ((index = row.indexOf('.', index + 1)) !== -1) {
                this.goalPositions.push({ x: index, y: y });
        }})
        console.log("ゴールの座標:", this.goalPositions);

        // 見つけた座標でPlayerのインスタンスの生成
        this.player = new Player(playerX, playerY);

        // プレイヤーの位置から"@"を削除"
        const playerRow = this.map[playerY]; // playerRowはプレイヤーのいる行
        this.map[playerY] = playerRow.substring(0 , playerX) + " " + playerRow.substring(playerX + 1);

    }
    // moveメソッドをmovePlayerメソッドに変更
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
            // 荷物のある行を文字列から配列に変換し、荷物があった場所を空白にしてからもう一度文字列に変換する
            const boxRow = this.map[nextY].split("");
            boxRow[nextX] = " ";
            this.map[nextY] = boxRow.join("");

            // 荷物の移動先の行を文字列から配列に変換し、荷物の移動先を荷物にしてからもう一度文字列に変換する
            const boxNextRow = this.map[boxNextY].split("");
            boxNextRow[boxNextX] = "o";
            this.map[boxNextY] = boxNextRow.join("");
                    
        }
        this.player.x = nextX;
        this.player.y = nextY;
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
        return this.goalPositions.every(pos => {
            return this.map[pos.y][pos.x] === 'o';
        });
    }
}

//プレイヤークラスを定義（プレイヤーの座標のみを管理）
class Player {
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }
}

// ゲームクラスを定義（入力を受付
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
        console.log('🎉 クリアおめでとう！ 🎉');
        process.exit(); // ゲームを終了する
        }
        });
    }

    start() {
        this.stage.display(); // Stageクラスのdisplayメソッドを呼び出して盤面を表示        
    }
}

const game = new Game();
game.start();

