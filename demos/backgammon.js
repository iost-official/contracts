class Game {
    constructor(a, b, size) {
        this.a = a;
        this.b = b;
        this.size = size;
        this.board = [];
        for (let i = 0; i < this.size; i++) {
            this.board.push([]);
            for (let j = 0; j < this.size; j++) {
                this.board[i].push(0)
            }
        }
        this.winner = null
    }

    move(id, x, y) {
        if (this.winner !== null) {
            return "this game has come to a close"
        }
        if (this.board[x][y] !== 0) {
            return "this slot has marked"
        }
        if (id === this.a) {
            this.board[x][y] = 1;
        } else if (id === this.b) {
            this.board[x][y] = 2;
        } else {
            return "you are not a player"
        }
        if (this._result(x, y)) {
            this.winner = id
        }
    }

    _result(x, y) {
        return this._count(x, y, 1, 0) >= 5 ||
            this._count(x, y, 0, 1) >= 5 ||
            this._count(x, y, 1, 1) >= 5 ||
            this._count(x, y, 1, -1) >= 5;
    }

    _count(x, y, stepx, stepy) {
        let count = 1;
        const color = this.board[x][y];
        let cx = x;
        let cy = y;
        for (let i = 0; i < 4; i ++) {
            cx += stepx;
            cy += stepy;
            if (!this._checkBound(cx) || !this._checkBound(cy)) break;
            if (color !== this.board[cx][cy]) break;
            count++
        }
        cx = x;
        cy = y;
        for (let i = 0; i < 4; i ++) {
            cx -= stepx;
            cy -= stepy;
            if (color !== this.board[cx][cy]) break;
            count++
        }
        return count
    }

    _checkBound(i) {
        return !(i < 0 || i >= this.size);

    }

    static fromJSON(json) {
        const obj = JSON.parse(json);
        let g = new Game(obj.a, obj.b);
        g.board = obj.board;
        return g
    }
}

class Backgammon {
    constructor() {
    }

    init() {
        storage.put("nonce", JSON.stringify(0));
    }

    newGameWith(b, size) {
        const jn = storage.get("nonce");
        const id = JSON.parse(jn);
        const newGame = new Game(tx.publisher, b, size);
        this._saveGame(b, newGame);
        storage.put("nonce", JSON.stringify(id + 1));
        return "no"+id
    }

    move(id, x, y) {
        let g = this._readGame(id);

        if (!g._checkBound(x) || !g._checkBound(y))
            throw "input out of bounds";

        const rtn = g.move(tx.publisher, x, y);
        if (rtn !== 0) {
            throw rtn
        }
    }

    accomplish(id) {
        let game = this._readGame(id);
        if (!BlockChain.requireAuth(game.a, "active") && !BlockChain.requireAuth(game.b, "active")) {
            throw "require auth error"
        }
        this._releaseGame(id)
    }

    _readGame(id) {
        const gj = storage.mapGet("games", "no" + id);
        return Game.fromJSON(gj)
    }

    _saveGame(id, game) {
        storage.mapPut("games", "no" + id, JSON.stringify(game), game.a);
    }

    _releaseGame(id) {
        storage.mapDel("games", "no" + id)
    }
}

module.exports = Backgammon;
