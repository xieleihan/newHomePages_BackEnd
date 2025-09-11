const code = require("svg-captcha"); // 引入图片验证码

function svgCaptcha() {
    return code.create({
        size: 6,
        ignoreChars: "0o1iIl",
        noise: 3,
        color: true,
        background: "#fff",
        fontSize: 60
    });
}

module.exports = svgCaptcha;