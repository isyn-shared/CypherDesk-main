if (DEBUG) {
    console.log("RSA_BITS:", RSA_BITS);
    console.log("PASSPHRASE", PASSPHRASE);
}

const privateKey = cryptico.generateRSAKey(PASSPHRASE, RSA_BITS);
if (DEBUG) console.log('privateKey', privateKey);
const publicKey = cryptico.publicKeyString(privateKey);
if (DEBUG) console.log('publicKey', publicKey);

let serverPublicKey = null;

// Moved to ws.js in initial events
// sendEvent("publicKey", {key: publicKey.cipher});