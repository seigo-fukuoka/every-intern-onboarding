// 追加要件
// 複数ステージの実装（最低3ステージ）
// ステージセレクト機能
    // 複数マップ用意しておき（Stageクラス）、ゲーム開始時にユーザーが選択できるようにする(Gameクラス)
// 移動回数のカウント
    // displayメソッドに移動回数を表示する機能を追加
// アンドゥ機能（1手戻る）
    // 現在の位置を保存しておき、"u"キーで一手戻る機能を実装
// txtファイルからステージデータを読み込む
    // むずそう

const fs = require('fs');

function loadLevelsFromFile(filePath) {
    try {
        // ファイルを巨大な1つの文字列として同期的に読み込む
        const fileContent = fs.readFileSync(filePath, 'utf8');

        // 文字列を処理して allMapData の形に変換する
        return fileContent
            .trim() // 文字列の最初と最後の余白や改行を削除
            .split(';') // ステージ区切り(;)で分割し、ステージごとの文字列の配列にする
            .map(levelString => levelString.trim()) // 各ステージの余白を削除
            .filter(levelString => levelString.length > 0) // 空のステージを削除
            .map(levelString => levelString.split('\n')); // 各ステージを改行(\n)で分割し、行の配列にする

    } catch (error) {
        // ファイルが読めなかった場合にエラーメッセージを表示して終了する
        console.error(`エラー: レベルファイルが読み込めませんでした。パス: ${filePath}`);
        process.exit(1); // プログラムを強制終了
    }
}

// 作成した関数を呼び出して、ファイルから allMapData を生成する
const allMapData = loadLevelsFromFile('levels.txt');

const MAP_SYMBOLS = {
PLAYER: '@',
BOX: 'o',
GOAL: '.',
WALL: '#',
FLOOR: " ",
BOX_ON_GOAL: "*"
};

const CONTROL_KEYS = {
    UP: "w",
    LEFT: "a",
    DOWN: "s",
    RIGHT: "d",
    RESET: "r",
    QUIT: "q",
    UNDO: "u"
};

//readlineモジュールをインポート
const readline = require("readline");

// Stageクラスを定義（盤面全体の管理者）
// 役割：盤面の状態（壁、荷物、プレイヤーの位置）をすべて把握し、ゲームのルールを実行する責任者
class Stage {
    constructor(selectedMapData, inputStageNumber) {
        this.inputStageNumber = inputStageNumber + 1;
        // 選択されたステージをallMapDataから抜き取り、配列として保持しておく
        this.map = selectedMapData.map(row => row.split("")); // 文字列を配列に変換        
        // ゴールの場所を(X,Y)座標で把握する
        this.goalPositions = [];
        // 荷物の場所を(X, Y)座標で把握する
        this.boxes = [];
        // プレイヤーの場所を(X, Y)座標で把握する
        let playerX;
        let playerY;
        

        this.map.forEach((row, y) => {
            row.forEach((char, x) => {
                if(char === MAP_SYMBOLS.GOAL) {
                    this.goalPositions.push({x: x, y: y})
                } else if (char === MAP_SYMBOLS.PLAYER) {
                    playerX = x;
                    playerY = y;
                    // 見つけた座標でPlayerのインスタンスの生成
                    this.player = new Player(playerX, playerY);
                    // プレイヤーの位置から"@"を削除"、displayメソッドを実行する前に盤面をまっさらにする必要がある
                    this.map[playerY][playerX] = MAP_SYMBOLS.FLOOR;
                } else if (char === MAP_SYMBOLS.BOX) {
                    this.boxes.push(new Box(x, y));
                    // 荷物の位置から"o"を削除
                    this.map[y][x] = MAP_SYMBOLS.FLOOR;
                }
            });
        });

        this.historyStack = []; // アンドゥ機能のための履歴スタック
        this.moveCount = 0; // 移動回数をカウントする変数
    } // ここまでコンストラクタ

    // Playerを移動させるメソッド
    movePlayer(dx, dy) {
        // 移動を行う前の盤面（床、壁、プレイヤー、荷物）の状態を保存
        this.historyStack.push(JSON.parse(JSON.stringify({
            player: this.player,
            boxes: this.boxes
        })));
        // 移動先の座標を計算
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
            
            if(isBlocked) {
                return;
            }

            targetBox.x += dx;
            targetBox.y += dy;

            this.player.move(dx, dy);
        } else { // 荷物がない場合
            this.player.move(dx, dy);
        } 
        if (this.moveCount > 0) {
            this.moveCount ++;    
        }
    }

    // 盤面全体を表示するメソッド
    display () {
        console.clear();
        console.log("=================================");
        console.log("Sokoban Level " + this.inputStageNumber);
        console.log("=================================");
        console.log("操作方法: w(上), a(左), s(下), d(右), r(リセット), q(終了)");
        console.log("=================================");
        console.log("移動回数: " + this.moveCount); // 移動回数を表示
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
                viewMap[box.y][box.x] = MAP_SYMBOLS.BOX_ON_GOAL;
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

    undo() {
        if (this.historyStack.length === 0) {
            return;
    }
        const lastState = this.historyStack.pop();
        this.player = new Player(lastState.player.x, lastState.player.y);
        this.boxes = lastState.boxes.map(undoData => new Box(undoData.x, undoData.y));
        this.moveCount --;  
        this.display(); // 盤面を際表示
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
        super(x, y, MAP_SYMBOLS.PLAYER);
    }

    // プレイヤー専用のメソッドはここに追加できる
    move(dx, dy) {
        this.x += dx;
        this.y += dy;
    }
}

class Box extends MovableObject {
    constructor(x, y) {
        super(x, y, MAP_SYMBOLS.BOX);
    }
}

// ゲームクラスを定義（ゲーム全体の司令塔、支配人）
// 役割：ユーザーからのキー入力を受付、それをStageクラスへの命令に変換する司令塔
class Game {
    constructor() {
        this.rl = readline.createInterface({
                    input: process.stdin,
                    output: process.stdout
                });
        this.selectedStageIndex = null;
    }
    // ユーザーからの入力を受け付けるメソッド
    setupInput() {
        readline.emitKeypressEvents(process.stdin);
        process.stdin.setRawMode(true);

        process.stdin.on('keypress', (str, key) => {
            switch (key.name) {
                case CONTROL_KEYS.QUIT:
                    process.exit();
                case CONTROL_KEYS.RESET:
                    this.reset();
                    return;
                case CONTROL_KEYS.UNDO:
                    this.stage.undo();
                case CONTROL_KEYS.UP:
                    this.stage.movePlayer(0, -1);
                    break;
                case CONTROL_KEYS.LEFT:
                    this.stage.movePlayer(-1, 0);
                    break;
                case CONTROL_KEYS.DOWN:
                    this.stage.movePlayer(0, 1);
                    break;
                case CONTROL_KEYS.RIGHT:
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

    selectStage() {
        console.log("ステージを選択してください: 「1」, 「2」, 「3」のどれかを入力してください");
        this.rl.question(
            `ステージを${allMapData.length}つの中から選択してください: 「1」, 「2」, 「3」のどれかを入力してください`, 
            (input) => {
                const selectedStageIndex = Number(input) - 1;  
                if (selectedStageIndex >= 0 && selectedStageIndex < allMapData.length) {
                    // 選択されたステージでクラスステージのインタスタンスを生成
                    this.selectedStageIndex = selectedStageIndex;
                    this.stage = new Stage(allMapData[this.selectedStageIndex], this.selectedStageIndex); 
                    this.setupInput();
                    this.stage.display();
                } else {
                    console.log("無効なステージ選択です。1, 2, 3 のいずれかを入力してください。");
                    this.selectStage(); // 再度選択を促す
                }
            }
        );
    }    
    // ゲームを開始するメソッド
    start() {
        this.selectStage();    
    }
    // ゲームをリセットするメソッド
    reset() {
        this.stage = new Stage(allMapData[this.selectedStageIndex], this.selectedStageIndex); // 新しいStageインスタンスを生成
        this.stage.display(); // 盤面を再表示
    }
}

const game = new Game();
game.start();

