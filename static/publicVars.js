const RSA_BITS = 1024;
const PASSPHRASE = Math.round(Math.random() * 10e17).toString();

let DEBUG = true;
var serverPublicKey = null;
var encryptionKey = null;