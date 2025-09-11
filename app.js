// 导入Koa框架
const Koa = require('koa');
// 导入Koa路由模块
const Router = require('@koa/router');
// 导入文件系统和路径模块
const fs = require('fs');
const path = require('path');

// 插件
// 获取环境变量插件
const dotenv = require('dotenv');

// 创建一个Koa对象表示web app本身
const app = new Koa();
const router = new Router();
const bodyParser = require('koa-bodyparser'); // 解析请求体
const cors = require('@koa/cors'); // 处理跨域
const compress = require('koa-compress'); // 响应压缩
// const helmet = require('koa-helmet'); // 安全相关
// 读取环境变量
dotenv.config();

// 使用中间件
// 安全头部设置
app.use(async (ctx, next) => {
    // 设置安全响应头
    ctx.set('X-Content-Type-Options', 'nosniff');
    ctx.set('X-Frame-Options', 'DENY');
    ctx.set('X-XSS-Protection', '1; mode=block');
    ctx.set('Strict-Transport-Security', 'max-age=31536000; includeSubDomains');
    ctx.set('Referrer-Policy', 'strict-origin-when-cross-origin');

    // 移除暴露服务器信息的头部
    ctx.remove('X-Powered-By');

    await next();
});

// 请求频率限制
const requestCounts = new Map();
app.use(async (ctx, next) => {
    const clientId = ctx.request.ip;
    const currentTime = Date.now();
    const windowTime = 60000; // 1分钟窗口
    const maxRequests = 2000; // 每分钟最大2000次请求

    if (!requestCounts.has(clientId)) {
        requestCounts.set(clientId, { count: 1, resetTime: currentTime + windowTime });
    } else {
        const clientData = requestCounts.get(clientId);
        if (currentTime > clientData.resetTime) {
            clientData.count = 1;
            clientData.resetTime = currentTime + windowTime;
        } else {
            clientData.count++;
            if (clientData.count > maxRequests) {
                ctx.status = 429;
                ctx.body = { code: 429, message: '太多次请求,请稍后再试' };
                return;
            }
        }
    }

    await next();
});

const getServerStatus = require('./utils/Modules/performance');

// 导入路由
const {
    TechnologyStack,
    WebPushRouter,
    ImgVerifyRouter,
    EmailVerifyRouter,
    UserRouter,
    SuperUserRouter,
    superServerStatus,
    SuperUserManageRouter,
    SuperGithubRouter,
    SuperFileRouter,
    SuperSystemConfigRouter,
    SuperAccessRouter,
    SuperAuthenticationRouter,
    BilibiliFollowInfoRouter
} = require('./router/index');

app.use(cors({
    allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
    // 生产环境中应该限制具体域名
    origin: function (ctx) {
        const origin = ctx.header.origin;
        // 在生产环境中，应该验证origin是否在允许列表中
        const allowedOrigins = ['https://localhost:5173', 'https://homepages.hongkong.atomglimpses.cn'];
        return allowedOrigins.includes(origin) ? origin : false;
        return origin;
    },
    credentials: true
}));

// 使用bodyparser
app.use(bodyParser());
// 使用 koa-compress 中间件
app.use(compress({
    threshold: 1024, // 超过 1KB 才进行压缩
    flush: require('zlib').Z_SYNC_FLUSH, // 立即刷新压缩数据
    gzip: {
        level: 9 // gzip 压缩级别，范围 0-9，越高压缩率越大
    }
}));
// app.use(helmet());

// 检测是否安装成功Koa
// app.use(async ctx => {
//     ctx.body = 'Hello World! Koa is working! Welcome to Hong Kong!';
// });

// 检查路由的正常(GET)
router.get('/test/get', async (ctx) => {
    console.log("触发了测试用的GET请求,网络已到达服务器");
    const date = new Date();
    const params = ctx.request.query; // 获取查询参数
    console.log("查询参数:", params);
    ctx.state = 200;
    ctx.body = {
        code: 200,
        message: 'Hello World! Koa get request is working! Welcome to Hong Kong!',
        type: 'test',
        source: ctx.req.socket.remoteAddress,
        time: date.toISOString()
    };
});

// 检查路由的正常(POST)
router.post('/test/post', async (ctx) => {
    console.log("触发了测试用的POST请求,网络已到达服务器");
    const params = ctx.request.body;
    console.log("POST参数:", params);
    const date = new Date();
    ctx.state = 200;
    ctx.body = {
        code: 200,
        message: 'Hello World! Koa post request is working! Welcome to Hong Kong!',
        type: 'test',
        source: ctx.req.socket.remoteAddress,
        time: date.toISOString()
    };
});

// 使用路由
app.use(router.routes());
app.use(router.allowedMethods());
app.use(TechnologyStack.routes()); // 技术栈图片路由
app.use(WebPushRouter.routes()); // WebPush路由
app.use(ImgVerifyRouter.routes()); // 图形验证码路由
app.use(EmailVerifyRouter.routes()); // 邮箱验证码路由
app.use(UserRouter.routes()); // 用户路由
app.use(SuperUserRouter.routes()); // 超级用户路由
app.use(superServerStatus.routes()); // 服务器状态路由
app.use(SuperUserManageRouter.routes()); // 超级用户管理路由
app.use(SuperGithubRouter.routes()); // Github路由
app.use(SuperFileRouter.routes()); // 文件路由
app.use(SuperSystemConfigRouter.routes()); // 系统配置路由
app.use(SuperAccessRouter.routes()); // 访问路由
app.use(SuperAuthenticationRouter.routes()); // 认证路由
app.use(BilibiliFollowInfoRouter.routes()); // B站关注的信息路由

// 安全配置
const ALLOWED_STATIC_DIR = path.resolve(__dirname, 'public', 'static');
const ALLOWED_EXTENSIONS = [
    '.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg', '.ico',
    '.js', '.mjs', '.jsx', '.ts', '.tsx',
    '.css', '.scss', '.sass', '.less',
    '.woff', '.woff2', '.ttf', '.otf', '.eot',
    '.txt', '.md', '.json'
];
const MAX_FILE_SIZE = 50 * 1024 * 1024; // 50MB

let logs = []; // 用于存储日志
// 备份原始 console.log
const originalLog = console.log;
const originalError = console.error;
const wss1 = new WebSocket.Server({ noServer: true });
const getServerStatusWss = new WebSocket.Server({ noServer: true });

// 拦截 console.log 并存储日志
console.log = (...args) => {
    const message = `[LOG] ${new Date().toISOString()} - ${args.join(" ")}`;
    logs.push(message);
    originalLog.apply(console, args); // 仍然打印到终端
    wss1.clients.forEach((client) => {
        client.send(JSON.stringify(logs));
    })
};

// 拦截 console.error 并存储日志
console.error = (...args) => {
    const message = `[ERROR] ${new Date().toISOString()} - ${args.join(" ")}`;
    logs.push(message);
    originalError.apply(console, args);
    wss1.clients.forEach((client) => {
        client.send(JSON.stringify(logs));
    })
};

router.get("/logs", async (ctx) => {
    const token = ctx.header.authorization;
    if (!token) {
        ctx.status = 401;
        ctx.body = { code: 401, message: '未登录' };
        return;
    }
    jwt.verify(token, SECRET_KEY, (err, decoded) => {
        if (err) {
            ctx.status = 401;
            ctx.body = { code: 401, message: '登录过期，请重新登录' };
            return;
        }
    });
    ctx.body = logs.slice(-200); // 只返回最近 50 条日志，避免数据过大
});

// 重启服务
router.post('/reset', async (ctx) => {
    const token = ctx.header.authorization;
    if (!token) {
        ctx.status = 401;
        ctx.body = { code: 401, message: '未登录' };
        return;
    }
    jwt.verify(token, SECRET_KEY, (err, decoded) => {
        if (err) {
            ctx.status = 401;
            ctx.body = { code: 401, message: '登录过期，请重新登录' };
            return;
        }
    });
    ctx.status = 200;
    ctx.body = { code: 200, message: '重置成功' };

    exec('pm2 restart app', (err, stdout, stderr) => {
        if (err) {
            console.error('重启失败:', err);
            return;
        }
        console.log('重启成功:', stdout);
    });
})

// 关闭服务
router.post('/stop', async (ctx) => {
    const token = ctx.header.authorization;
    if (!token) {
        ctx.status = 401;
        ctx.body = { code: 401, message: '未登录' };
        return;
    }
    jwt.verify(token, SECRET_KEY, (err, decoded) => {
        if (err) {
            ctx.status = 401;
            ctx.body = { code: 401, message: '登录过期，请重新登录' };
            return;
        }
    });
    ctx.status = 200;
    ctx.body = { code: 200, message: '停止成功' };

    exec('pm2 stop app', (err, stdout, stderr) => {
        if (err) {
            console.error('停止失败:', err);
            return;
        }
        console.log('停止成功:', stdout);
    });
});

router.get('/processes', async (ctx) => {
    const psList = (await import('ps-list')).default;
    const processes = await psList();
    ctx.body = {
        total: processes.length,
        list: processes.slice(0, 10) // 只返回前 10 个进程
    };
});

// 获取静态文件列表的API (安全版本)
router.get('/api/static-files', async (ctx) => {
    try {
        // 安全检查：验证客户端IP (可选，根据需求启用)
        // const clientIP = ctx.request.ip;
        // if (!isAllowedIP(clientIP)) {
        //     ctx.status = 403;
        //     ctx.body = { code: 403, message: 'Access denied' };
        //     return;
        // }

        // 添加简单的访问频率限制 (可选)
        const currentTime = Date.now();
        const clientId = ctx.request.ip + ctx.request.header['user-agent'];
        // 这里可以实现更复杂的频率限制逻辑

        function isValidPath(targetPath) {
            const resolvedPath = path.resolve(targetPath);
            return resolvedPath.startsWith(ALLOWED_STATIC_DIR);
        }

        function isAllowedFile(filePath) {
            const ext = path.extname(filePath).toLowerCase();
            return ALLOWED_EXTENSIONS.includes(ext);
        }

        function sanitizeFileName(name) {
            // 移除路径遍历字符和特殊字符
            return name.replace(/[\.\/\\:*?"<>|]/g, '');
        }

        function walkDir(dir, baseDir) {
            let results = [];

            // 安全检查：确保目录在允许范围内
            if (!isValidPath(dir)) {
                console.warn(`Attempted to access unauthorized directory: ${dir}`);
                return results;
            }

            try {
                const files = fs.readdirSync(dir);

                files.forEach(file => {
                    try {
                        const filePath = path.join(dir, file);

                        // 安全检查：验证路径
                        if (!isValidPath(filePath)) {
                            console.warn(`Skipping invalid path: ${filePath}`);
                            return;
                        }

                        const stat = fs.statSync(filePath);

                        // 检查文件大小限制
                        if (stat.size > MAX_FILE_SIZE) {
                            console.warn(`File too large, skipping: ${file} (${stat.size} bytes)`);
                            return;
                        }

                        if (stat && stat.isDirectory()) {
                            // 递归处理子目录
                            results = results.concat(walkDir(filePath, baseDir));
                        } else if (isAllowedFile(filePath)) {
                            const relativePath = path.relative(baseDir, filePath).replace(/\\/g, '/');
                            const ext = path.extname(file).toLowerCase().replace('.', '');

                            results.push({
                                name: sanitizeFileName(file),
                                path: '/' + relativePath,
                                ext: ext || 'unknown',
                                size: stat.size,
                                sizeFormatted: formatBytes(stat.size),
                                lastModified: stat.mtime.toISOString()
                            });
                        }
                    } catch (fileError) {
                        console.error(`Error processing file ${file}:`, fileError.message);
                        // 继续处理其他文件，不中断整个过程
                    }
                });
            } catch (dirError) {
                console.error(`Error reading directory ${dir}:`, dirError.message);
            }

            return results;
        }

        function formatBytes(bytes) {
            if (bytes === 0) return '0 Bytes';
            const k = 1024;
            const sizes = ['Bytes', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }

        // 检查静态目录是否存在
        if (!fs.existsSync(ALLOWED_STATIC_DIR)) {
            ctx.status = 404;
            ctx.body = {
                code: 404,
                message: 'Static directory not found'
            };
            return;
        }

        const files = walkDir(ALLOWED_STATIC_DIR, ALLOWED_STATIC_DIR);

        // 记录访问日志
        console.log(`Static files API accessed by ${ctx.request.ip} at ${new Date().toISOString()}`);

        ctx.body = {
            code: 200,
            message: 'success',
            data: files,
            timestamp: new Date().toISOString()
        };
    } catch (error) {
        // 安全的错误处理：不暴露详细的系统信息
        console.error('Error in static-files API:', error);
        ctx.status = 500;
        ctx.body = {
            code: 500,
            message: 'Internal server error',
            timestamp: new Date().toISOString()
        };
    }
});

router.post('/proxy', async (ctx) => {
    const { url,fileType } = ctx.request.body;
    if (!url) {
        ctx.status = 400;
        ctx.body = { code: 400, message: '缺少url参数' };
        return;
    }
    if (!fileType) {
        ctx.status = 400;
        ctx.body = { code: 400, message: '缺少存放路径' };
        return;
    }

    try {
        // 确保目录存在
        const saveDir = path.join(__dirname, `public/static/${fileType}`);
        if (!fs.existsSync(saveDir)) {
            fs.mkdirSync(saveDir, { recursive: true });
        }

        // 从url里取文件名
        const fileName = path.basename(new URL(url).pathname);
        const filePath = path.join(saveDir, fileName);

        // 请求远程资源，保存为文件
        // 检查本地是否已经有相同文件名的文件
        if (!fs.existsSync(filePath)) {
            const response = await axios.get(url, {
                responseType: 'arraybuffer', // 二进制流
                headers: {
                    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3'
                }
            });

            fs.writeFileSync(filePath, response.data);
        }


        // 拼接静态访问路径（前提：Koa挂载了static目录）
        const fileUrl = `/static/${fileType}/${fileName}`;

        ctx.body = {
            code: 200,
            message: '下载成功',
            url: fileUrl
        };
    } catch (error) {
        ctx.status = 500;
        ctx.body = { code: 500, message: '请求失败', error: error.message };
    }
});

router.post('/verifyFriend', async (ctx) => {
    const { password } = ctx.request.body; // 前端传递过来的密码
    console.log('收到的密码:', password);
    const FRIEND_PASSWORD = process.env.FRIEND_PASSWORD; // 预设的密码，存储在环境变量中
    if (password === null) {
        ctx.status = 400;
        ctx.body = { code: 400, message: '缺少密码参数' };
        return;
    } else {
        if (password == FRIEND_PASSWORD) {
            ctx.body = { code: 200, message: '验证成功', imgUrl: '/static/wechat.jpg' };
        } else {
            ctx.status = 401;
            ctx.body = { code: 401, message: '密码错误' };
        }
    }
})

// 运行
app.use(router.routes()).use(router.allowedMethods());

// 静态资源分发
app.use(require('koa-static')(__dirname + '/public'));

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});

app.on('upgrade', (request, socket, head) => {
    const { url } = request;
    if (url === '/logs') {
        wss1.handleUpgrade(request, socket, head, (ws) => {
            wss1.emit('connection', ws, request);
        });
    } else if (url === '/private/superServerStatus') {
        getServerStatusWss.handleUpgrade(request, socket, head, (ws) => {
            getServerStatusWss.emit('connection', ws, request);
        });
    } else {
        socket.destroy();
    }
});

wss1.on('connection', (ws) => {
    // 连接成功后，发送最近的日志
    ws.send(JSON.stringify(logs.slice(-200)));
});

// WebSocket 连接
getServerStatusWss.on('connection', async (ws) => {
    console.log('WebSocket 客户端已连接');

    // 发送初始数据
    ws.send(JSON.stringify(await getServerStatus()));

    // 每 1 秒发送一次服务器状态更新
    const interval = setInterval(async () => {
        if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify(await getServerStatus()));
        } else {
            clearInterval(interval);
        }
    }, 1000);

    // 关闭连接时清理定时器
    ws.on('close', () => {
        console.log('WebSocket 连接关闭');
        clearInterval(interval);
    });
});