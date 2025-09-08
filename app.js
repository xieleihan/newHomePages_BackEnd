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
app.use(cors({
    allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
    // 允许所有域名访问
    origin: function (ctx) {
        return ctx.header.origin;
    }
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

// 获取静态文件列表的API
router.get('/api/static-files', async (ctx) => {
    try {
        const staticDir = path.join(__dirname, 'public', 'static');

        function walkDir(dir, baseDir) {
            let results = [];
            const files = fs.readdirSync(dir);

            files.forEach(file => {
                const filePath = path.join(dir, file);
                const stat = fs.statSync(filePath);

                if (stat && stat.isDirectory()) {
                    results = results.concat(walkDir(filePath, baseDir));
                } else {
                    const relativePath = path.relative(baseDir, filePath).replace(/\\/g, '/');
                    const ext = path.extname(file).toLowerCase().replace('.', '');
                    const size = stat.size;

                    results.push({
                        name: file,
                        path: '/' + relativePath,
                        ext: ext || 'unknown',
                        size: size,
                        sizeFormatted: formatBytes(size)
                    });
                }
            });

            return results;
        }

        function formatBytes(bytes) {
            if (bytes === 0) return '0 Bytes';
            const k = 1024;
            const sizes = ['Bytes', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }

        const files = walkDir(staticDir, staticDir);

        ctx.body = {
            code: 200,
            message: 'success',
            data: files
        };
    } catch (error) {
        console.error('Error reading static files:', error);
        ctx.status = 500;
        ctx.body = {
            code: 500,
            message: 'Error reading static files',
            error: error.message
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