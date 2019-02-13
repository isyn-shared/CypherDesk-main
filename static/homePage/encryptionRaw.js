const NodeRSA = require('node-rsa');
window.encryptionKey = new NodeRSA({b: RSA_BITS});
if (DEBUG) console.log('privateKey', window.encryptionKey.exportKey());

window.NodeRSA = NodeRSA;

if (DEBUG) {
    console.log("RSA_BITS:", RSA_BITS);
    // console.log("PASSPHRASE", PASSPHRASE);
}

if (DEBUG) console.log('publicKey', window.encryptionKey.exportKey('pkcs1-public'));

// Moved to ws.js in initial events
// sendEvent("publicKey", {key: publicKey});