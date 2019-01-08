// Storage.js
class Test {
    init() {
        storage.put("value1", "foobar")
    }
    read() {
        console.log(storage.get("value1"))
    }
    change(someone) {
        storage.put("value1", someone)
    }
}
module.exports = Test;