class Bet {
    constructor(player, amount, number, player_seed, hash, inviter) {
        this.amount = amount;
        this.number = number;
        this.hostHash = hash;
        this.player = player;
        this.playerSeed = player_seed;
        this.hostSeed = "0000000000000000";
        this.inviter = inviter;
        this.status = "pending"
    }

    reveal(seed) {
        if (IOSTCrypto.sha3(seed) !== this.hostHash) {
            throw "illegal seed"
        }
        this.hostSeed = seed;
    }

    dice() {
        const hash = IOSTCrypto.sha3(this.playerSeed + this.hostSeed);
        let sum = 0;
        for (let i = 0; i < hash.length; i ++) {
            sum += hash.charCodeAt(i)
        }
        return sum % 100 + 1
    }
}

class Dice {
    init(){
        storage.put("nonce", JSON.stringify(0))
    }
    bet(amount, number, player_seed, hash, inviter) {
        if (number < 3 || number > 96) {
            throw "illegal bet number"
        }
        blockchain.deposit(tx.publisher, amount, "");

        const b = new Bet(tx.publisher, amount, number, player_seed, hash, inviter);
        const nonce = this._new_bet(b);
        return "bet"+nonce
    }
    reveal(nonce, seed, isVIP) {
        const bet = this._load_bet(nonce);

        // pay et
        this._pay_et(bet.amount, bet.player);
        this._pay_et(bet.amount / 10, bet.inviter);

        bet.reveal(seed);
        if (bet.dice() < bet.number) {
            bet.status = "win";
            const base = isVIP?99:98;
            let amount = new Float64(bet.amount);
            const pay = amount.multi(base).div(bet.number - 1);
            blockchain.deposit(bet.player, pay, "");
        } else {
            bet.status = "lose"
            // TODO Transfer to private account
        }
        this._save_bet(bet);
        return bet
    }
    del(nonce) {
        const bet = this._load_bet("bet"+nonce);
        if (!blockchain.requireAuth(bet.player, "active")) {
            throw "only player can delete bet records"
        }
        if (bet.status === 'pending') {
            throw "bet unfinished"
        }
        storage.del("bet"+nonce)
    }
    _pay_et(iostAmount, player) {

    }
    _new_bet(bet) {
        const nonce = JSON.parse(storage.get("nonce"));
        this._save_bet(bet, nonce);
        storage.put("nonce", (nonce + 1).toString());
        return nonce
    }
    _save_bet(bet, nonce) {
        storage.put("bet"+nonce, JSON.stringify(bet), bet.player);
    }
    _load_bet(nonce) {
        const jsonStr = storage.get("bet"+nonce);
        if (!jsonStr) {
            throw "bet not exist"
        }
        const obj = JSON.parse(jsonStr);
        return Object.assign(new Bet, obj)
    }
}

module.exports = Dice;