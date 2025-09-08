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

console.log("测试环境变量是否正常注入:", process.env.RUNTIME_TEXT || '环境变量未设置');

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
    const maxRequests = 100; // 每分钟最大100次请求

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
                ctx.body = { code: 429, message: 'Too many requests' };
                return;
            }
        }
    }

    await next();
});

app.use(cors({
    allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
    // 生产环境中应该限制具体域名
    origin: function (ctx) {
        const origin = ctx.header.origin;
        // 在生产环境中，应该验证origin是否在允许列表中
        const allowedOrigins = ['http://localhost:3000', 'https://homepages.hongkong.atomglimpses.cn'];
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

// 运行
app.use(router.routes()).use(router.allowedMethods());

// 静态资源分发
app.use(require('koa-static')(__dirname + '/public'));

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});