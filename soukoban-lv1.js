// ステージ上のモノを定数として定義する
const MAP_SYMBOLS = {
    PLAYER: '@',
    BOX: 'o',
    GOAL: '.',
    WALL: '#',
    FLOOR: " ",
    BOXONGOAL: "*"
};

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
        
        // ゴールの場所を(X,Y)座標で把握する
        this.goalPositions = [];
        // 荷物の場所を(X, Y)座標で把握する
        this.boxes = [];
        // プレイヤーの場所を(X, Y)座標で把握する
        let playerX;
        let playerY;

        this.map.forEach((row, y) => {
            row.forEach((char, x) => {
                if(char === MAP_SYMBOLS.GOAL){
                    this.goalPositions.push({x: x, y: y})
                } else if (char === MAP_SYMBOLS.PLAYER) {
                    playerX = x;
                    playerY = y;
                } else if (char === MAP_SYMBOLS.BOX) {
                    this.boxes.push(new Box(x, y));
                    this.map[y][x] = MAP_SYMBOLS.FLOOR;
                }
            });
        });
        // 見つけた座標でPlayerのインスタンスの生成
        this.player = new Player(playerX, playerY);
        // プレイヤーの位置から"@"を削除"
        this.map[playerY][playerX] = MAP_SYMBOLS.FLOOR;
    }
    // Playerを移動させるメソッド
    movePlayer(dx, dy) {
        // プレイヤーの移動先の座標を計算
        const nextX = this.player.x + dx;
        const nextY = this.player.y + dy;
        // 移動先が壁なら何もしない
        if (this.map[nextY][nextX] === MAP_SYMBOLS.WALL) {
            return;
        }
        // 移動先に荷物があるか、this.boxes 配列から検索する
        const targetBox = this.boxes.find(box => box.x === nextX && box.y === nextY);
        // 荷物があった場合
        // 荷物の一個先が壁か他の荷物だった場合、何もしない
        if (targetBox) {
            const boxNextX = targetBox.x + dx;
            const boxNextY = targetBox.y + dy;
            const isBlocked = this.map[boxNextY][boxNextX] === MAP_SYMBOLS.WALL ||
                            this.boxes.some(box => box.x === boxNextX && box.y === boxNextY);
            
            if(isBlocked){
                return;
            }

            targetBox.x += dx;
            targetBox.y += dy;

            this.player.move(dx, dy);
        } else { // 荷物がない場合
            this.player.move(dx, dy);
        }         
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
            if (this.map[goal.y][goal.x] === MAP_SYMBOLS.FLOOR) {
                viewMap[goal.y][goal.x] = MAP_SYMBOLS.GOAL;
            }
        });
        // boxes配列の情報を元に、荷物をviewMapに描画する
        // 荷物の位置とゴールの位置が被ってたら表記を変える
        this.boxes.forEach(box => {
            const isOnGoal = this.goalPositions.some(goal => goal.x === box.x && goal.y === box.y);
            if (isOnGoal){
                viewMap[box.y][box.x] = MAP_SYMBOLS.BOXONGOAL;
            } else {
                viewMap[box.y][box.x] = box.symbol;
            }     
        });

        // プレイヤーがいる行を文字列から配列に変換
        viewMap[this.player.y][this.player.x] = this.player.symbol;      // プレイヤーの位置に"@"を置く
        // 文字列に戻してマップに反映
        viewMap.forEach(rowArray => {
            console.log(rowArray.join(""));
        })
    }
    // クリア判定を行うメソッド
    // 盤面上のゴールの座標を把握しておき、すべての座標に荷物が置かれているかチェックする
    // GameクラスのisClearメソッドから呼び出される
    isClear() {
        return this.goalPositions.every(goal => {
            return this.boxes.some(box => box.x === goal.x && box.y === goal.y)
        })
    }
}

class MovableObject {
    constructor(x, y, symbol) {
        this.x = x;
        this.y = y;
        this.symbol = symbol; // 表示用の記号
    }
}

//プレイヤークラスを定義（プレイヤーの座標のみを管理する駒）
// 役割：
class Player extends MovableObject {
    constructor(x, y) {
        // super()で親のconstructorを呼び出す
        super(x, y, "@");
    }

    // プレイヤー専用のメソッドはここに追加できる
    move(dx, dy) {
        this.x += dx;
        this.y += dy;
    }
}

class Box extends MovableObject {
    constructor(x, y) {
        super(x, y, "o");
    }
}

// ゲームクラスを定義（ゲーム全体の司令塔、支配人）
// 役割：ユーザーからのキー入力を受付、それをStageクラスへの命令に変換する司令塔
class Game {
    constructor() {
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
            switch (key.name) {
                case "w":
                    this.stage.movePlayer(0, -1);
                    break;
                case "a":
                    this.stage.movePlayer(-1, 0);
                    break;
                case "s":
                    this.stage.movePlayer(0, 1);
                    break;
                case "d":
                    this.stage.movePlayer(1, 0);
                    break;
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
        this.stage = new Stage();
        this.setupInput();   
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

