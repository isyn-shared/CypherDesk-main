const NodeRSA = require('node-rsa');
const privateKey = new NodeRSA({b: RSA_BITS});
if (DEBUG) console.log('privateKey', privateKey.exportKey());

window.NodeRSA = NodeRSA;

if (DEBUG) {
    console.log("RSA_BITS:", RSA_BITS);
    console.log("PASSPHRASE", PASSPHRASE);
}

window.publicKey = privateKey.exportKey('pkcs8-public');
if (DEBUG) console.log('publicKey', window.publicKey);

// Moved to ws.js in initial events
// sendEvent("publicKey", {key: publicKey});