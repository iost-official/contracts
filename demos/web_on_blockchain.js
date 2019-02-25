class Web {

    init() {
        storage.put("nonce", JSON.stringify(0));
    }

    deploy(webstring){
        const jn = storage.get("nonce");
        const id = JSON.parse(jn);
        storage.put("w"+id, webstring, tx.publisher);
        storage.put("nonce", JSON.stringify(id + 1));
        return "w"+id
    }
}

module.exports = Web;
