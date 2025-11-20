// Web Worker 用于执行模幂运算

function modPow(base, exponent, modulus) {
    if (modulus === 1n) return 0n;
    let result = 1n;
    base = base % modulus;
    while (exponent > 0n) {
        if (exponent % 2n === 1n) {
            result = (result * base) % modulus;
        }
        exponent = exponent >> 1n;
        base = (base * base) % modulus;
    }
    return result;
}

self.onmessage = (event) => {
    const { base, exponent, modulus } = event.data;
    try {
        console.log(`[Worker] 开始计算 modPow，base=${base.substring(0, 20)}..., exp=${exponent.substring(0, 20)}..., mod=${modulus.substring(0, 20)}...`);

        const result = modPow(BigInt(base), BigInt(exponent), BigInt(modulus));
        const resultStr = result.toString();

        console.log(`[Worker] 计算完成，结果长度: ${resultStr.length} chars`);

        self.postMessage({ success: true, result: resultStr });
    } catch (error) {
        console.error(`[Worker] 错误: ${error.message}`);
        self.postMessage({ success: false, error: error.message });
    }
};
