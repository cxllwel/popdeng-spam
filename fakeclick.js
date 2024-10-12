function base64url(source) {
    return CryptoJS.enc.Base64.stringify(source).replace(/=+$/, "").replace(/\+/g, "-").replace(/\//g, "_");;
}

function jwtSign(payload, secret) {
    const encodedHeader = base64url(CryptoJS.enc.Utf8.parse(JSON.stringify({"alg": "HS256"})));
    const encodedPayload = base64url(CryptoJS.enc.Utf8.parse(JSON.stringify(payload)));
    const encodedSignature = base64url(CryptoJS.HmacSHA256(encodedHeader + "." + encodedPayload, secret));
    return encodedHeader + "." + encodedPayload + "." + encodedSignature;
}

fetch("https://cdnjs.cloudflare.com/ajax/libs/crypto-js/4.2.0/crypto-js.min.js").then(r=>r.text()).then(eval).then(function(){
    localStorage.setItem("_d", jwtSign({"click":999_999_999_999,"iat":parseInt(Date.now()/10000),"By":"Larinax999"}, "Please, don\"t share me."));
    console.log("Done");
});

/* get ur countryCode
self.__next_f[7][1].split('{"countryCode":"')[1].split('","initLeaderboard"')[0]
*/