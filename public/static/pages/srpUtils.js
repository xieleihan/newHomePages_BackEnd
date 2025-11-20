/**
 * SRP 注册辅助函数库（使用 PBKDF2 和 Web Worker）
 * 完全按照已验证的 TypeScript 版本实现
 */

// RFC 5054 3072-bit 参数 (N，768 个十六进制字符 = 3072 bits)
const RFC5054_N_HEX = 'FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE65381FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE65381FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE65381FFFFFFFFFFFFFFFF';
const RFC5054_g = 2n;

// 验证常量
console.log(`[SRPUtils] RFC5054_N_HEX 长度: ${RFC5054_N_HEX.length} (应该是 768)`);
if (RFC5054_N_HEX.length !== 768) {
    console.error(`[ERROR] N 常量长度错误！期望 768，实际 ${RFC5054_N_HEX.length}`);
}

let workerPool = null;

/**
 * 获取 Worker 池
 */
function getWorkerPool(poolSize = 2) {
    if (!workerPool) {
        workerPool = new CryptoWorkerPool(poolSize);
    }
    return workerPool;
}

/**
 * 销毁 Worker 池
 */
function destroyWorkerPool() {
    if (workerPool) {
        workerPool.destroy();
        workerPool = null;
    }
}

/**
 * Worker 池类 - 管理多个 Worker 实例
 */
class CryptoWorkerPool {
    constructor(poolSize = 2) {
        this.workers = new Map();
        this.queue = [];
        this.poolSize = poolSize;
        this.useWorkers = typeof Worker !== 'undefined';
        this.initializeWorkers();
    }

    initializeWorkers() {
        if (!this.useWorkers) {
            console.log('[WorkerPool] Web Worker 不可用，将使用主线程计算');
            return;
        }

        try {
            for (let i = 0; i < this.poolSize; i++) {
                const worker = new Worker('/static/pages/cryptoWorker.js');
                this.workers.set(worker, false); // false = 不忙
            }
            console.log(`[WorkerPool] 初始化 ${this.poolSize} 个 Worker`);
        } catch (error) {
            console.error('[WorkerPool] 初始化 Worker 失败:', error);
            this.useWorkers = false;
        }
    }

    async modPow(base, exponent, modulus) {
        if (!this.useWorkers) {
            // 降级到主线程计算
            console.log('[WorkerPool] 使用主线程进行 modPow 计算');
            return modPowDirect(base, exponent, modulus);
        }

        return new Promise((resolve, reject) => {
            const task = {
                resolve,
                reject,
                data: {
                    base: base.toString(),
                    exponent: exponent.toString(),
                    modulus: modulus.toString()
                }
            };
            this.queue.push(task);
            this.processQueue();
        });
    }

    processQueue() {
        if (this.queue.length === 0) return;

        let availableWorker = null;
        for (const [worker, busy] of this.workers) {
            if (!busy) {
                availableWorker = worker;
                break;
            }
        }

        if (!availableWorker) return; // 没有空闲 Worker

        const task = this.queue.shift();
        if (!task) return;

        this.workers.set(availableWorker, true); // 标记为忙

        const messageHandler = (event) => {
            availableWorker.removeEventListener('message', messageHandler);
            availableWorker.removeEventListener('error', errorHandler);
            this.workers.set(availableWorker, false); // 标记为空闲

            if (event.data.success) {
                task.resolve(BigInt(event.data.result));
            } else {
                task.reject(new Error(event.data.error));
            }

            this.processQueue();
        };

        const errorHandler = (error) => {
            availableWorker.removeEventListener('message', messageHandler);
            availableWorker.removeEventListener('error', errorHandler);
            this.workers.set(availableWorker, false);
            task.reject(new Error(error.message));
            this.processQueue();
        };

        availableWorker.addEventListener('message', messageHandler);
        availableWorker.addEventListener('error', errorHandler);

        console.log('[WorkerPool] 发送任务到 Worker，base=' + task.data.base.substring(0, 20) + '..., exp=' + task.data.exponent.substring(0, 20) + '..., mod=' + task.data.modulus.substring(0, 20) + '...');
        availableWorker.postMessage(task.data);
    }

    destroy() {
        this.workers.forEach((_, worker) => worker.terminate());
        this.workers.clear();
        this.queue = [];
    }
}

/**
 * 主线程 modPow 实现（降级方案）
 */
function modPowDirect(base, exponent, modulus) {
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

/**
 * 使用 PBKDF2 派生密钥
 */
async function computeHashedPassword(
    salt,
    password,
    iterations = 600000,
    keyLength = 64,
    hash = 'SHA-512'
) {
    console.log(`[PBKDF2] 开始计算，salt=${salt.substring(0, 20)}..., password=${password}, iterations=${iterations}`);

    // 盐十六进制转换成 Uint8Array
    const saltBytes = new Uint8Array(salt.match(/.{2}/g).map(byte => parseInt(byte, 16)));

    // 密码转换成 utf-8
    const encoder = new TextEncoder();
    const passwordBytes = encoder.encode(password);

    const keyMaterial = await window.crypto.subtle.importKey(
        'raw',
        passwordBytes,
        'PBKDF2',
        false,
        ['deriveBits']
    );

    const derivedBits = await window.crypto.subtle.deriveBits(
        {
            name: 'PBKDF2',
            salt: saltBytes,
            iterations,
            hash: { name: hash }
        },
        keyMaterial,
        keyLength * 8 // bits
    );

    const hashArray = Array.from(new Uint8Array(derivedBits));
    const result = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');

    console.log(`[PBKDF2] 完成，结果长度: ${result.length} chars (${result.substring(0, 30)}...)`);
    return result;
}

/**
 * 生成验证器 v = g^x mod N
 */
async function computerVerifier(xHex, nHex = RFC5054_N_HEX, gBigInt = RFC5054_g) {
    console.log(`[computerVerifier] 开始计算，xHex=${xHex.substring(0, 30)}..., nHex长度=${nHex.length}`);

    const xBigInt = BigInt('0x' + xHex);
    const nBigInt = BigInt('0x' + nHex);

    console.log(`[computerVerifier] 转换完成，计算 g^x mod N...`);

    const workerPool = getWorkerPool();
    const v = await workerPool.modPow(gBigInt, xBigInt, nBigInt);

    let vHex = v.toString(16);
    console.log(`[computerVerifier] 计算完成，vHex长度: ${vHex.length}`);

    return vHex;
}

/**
 * 生成随机盐
 */
function generateSalt(length = 16) {
    const array = new Uint8Array(length);
    window.crypto.getRandomValues(array);
    return Array.from(array, byte => byte.toString(16).padStart(2, '0')).join('');
}

/**
 * 返回盐和验证器
 */
async function returnVerifierAndSalt(password) {
    console.log(`[returnVerifierAndSalt] 开始生成盐和验证器`);

    const salt = generateSalt(16); // 16字节盐
    console.log(`[returnVerifierAndSalt] Salt: ${salt}`);

    const xHex = await computeHashedPassword(salt, password);
    console.log(`[returnVerifierAndSalt] 已计算 x: ${xHex.substring(0, 30)}...`);

    const verifier = await computerVerifier(xHex);
    console.log(`[returnVerifierAndSalt] 已计算 Verifier，长度: ${verifier.length}`);

    return { salt, verifier };
}
