// 导入Koa框架
const Koa = require('koa');
// 导入Koa路由模块
const Router = require('@koa/router');

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

// 检测是否安装成功Koa
// app.use(async ctx => {
//     ctx.body = 'Hello World! Koa is working! Welcome to Hong Kong!';
// });

// 检查路由的正常(GET)
router.get('/test/get', async (ctx) => {
    ctx.state = 200;
    ctx.body = {
        code: 200,
        message: 'Hello World! Koa get request is working! Welcome to Hong Kong!',
    };
});

// 检查路由的正常(POST)
router.post('/test/post', async (ctx) => {
    ctx.state = 200;
    ctx.body = {
        code: 200,
        message: 'Hello World! Koa post request is working! Welcome to Hong Kong!',
    };
});

// 运行
app.use(router.routes()).use(router.allowedMethods());
const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});