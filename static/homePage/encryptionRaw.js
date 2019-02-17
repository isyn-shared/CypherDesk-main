const NodeRSA = require('node-rsa');
window.encryptionKey = new NodeRSA({b: RSA_BITS}, 'pkcs1-private-pem', {
    'encryptionScheme': 'pkcs1'
});
if (DEBUG) console.log('privateKey', window.encryptionKey.exportKey());

window.NodeRSA = NodeRSA;

if (DEBUG) {
    console.log("RSA_BITS:", RSA_BITS);
    // console.log("PASSPHRASE", PASSPHRASE);
}

if (DEBUG) console.log('publicKey', window.encryptionKey.exportKey('pkcs1-public'));
if (DEBUG) {
    console.log("Testing key;");
    let plainText = "CYKA BLYAT";
    let encText = encryptionKey.encrypt(plainText);
    console.log('encText', encText);
    let decText = encryptionKey.decrypt(encText);
    console.log('decText', decText);

}


// Moved to ws.js in initial events
// sendEvent("publicKey", {key: publicKey});