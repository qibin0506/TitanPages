### TitanPages, A fast, simple static blog builder, powered by golang.

TitanPages是一个静态博客生成器, 可以配合githubpages轻松的搭建自己的博客.

### 依赖项目

[https://github.com/russross/blackfriday](https://github.com/russross/blackfriday)

[https://github.com/toqueteos/webbrowser](https://github.com/toqueteos/webbrowser)

demo中使用的主题来自:[http://kywk.github.io/](http://kywk.github.io/)

#### 展示网站请转到: [https://qibin0506.github.io/](https://qibin0506.github.io/)

### 更新日志

#### 2016/7/8更新
1. 修复生成目录标题错误问题
2. 修改生成目录列表按照时间排序

#### 2016/7/2更新
1. 修复生成摘要时保持markdown格式bug
2. 加入-help参数, 可以在命令行查看使用文档,使用方法: `tt -help keyword` e.g. `tt -help build`查看build的使用方法

### 如何使用

#### step 1. 

下载源码编译源码(linux用户, 可以直接下载`tt`文件;windows用户可直接下载tt.exe)

#### step 2. 

创建文件, 在你的工作空间用命令行运行以下命令:

`tt -type create -file 你的文件名称`

例如: `tt -type create -file 我的第一篇博客`

#### step 3. 

写作, 打开`/raw/你的文件名称`文件, 进行文章的书写(注意: 文章的格式必须是markdown的)

#### step 4. 

编译markdown文件,写作完成后, 运行命令:

`tt -type build -file 你的文件名称 [-author 作者] [-tmpl 要使用的模板文件]`

例如: `tt -type build -file 我的第一篇博客 -author 亓斌 -tmpl ./content.html`

(注意: []中的参数为可选参数, 具体content.html模板如何书写会在下面介绍)

现在在/html目录下会生成对应文件名的html文件.

#### step 5.

生成目录, 运行命令:

`tt -type cate`

运行该命令, 在/html目录中会生成一个`category.auto.js`的javascript文件.

#### step 6.

文章模板文件content.html的书写:

1. 使用占位符`{{.Title}}`表示文章的标题 
2. 使用占位符`{{.Date}}`表示文章的日期
3. 使用占位符`{{.Author}}`表示文章的作者
4. 使用占位符`{{.Desc}}`表示文章的描述
5. 使用占位符`{{.Content}}`表示文章内容

**注意: 关于占位符`{{.desc}}`的说明: 建议将这个描述放在`<meta name='description'></meta>`中,这样,在生成目录的时候才会产生摘要信息.**

#### step 7.

关于自动生成的`category.auto.js`文件的说明, 这个文件是关于文章索引信息的, 我们需要在目录页调用这个文件里的函数:

1. `pageCount()` 函数会返回分页页码总数(默认分页大小为5)
2. `getQueryString(query)` 函数可以获取指定的querystring参数, 通常我们用来获取当前页码
3. `get(currentPage)` 函数会根据当前页码返回数据数组, 该数组中包含了索引页需要的信息

索引信息数组中包含的信息如下:

1. `title` 文章的标题
2. `date` 文章生成的时间
3. `desc` 文章的简要描述

demo中的例子: 

``` javascript
window.onload = function() {
	var page = getQueryString("page")
	var count = pageCount()
	if (page == null) {
		page = 1
	}else {
		page = parseInt(page)
	}

	if(page > 1) {
		document.getElementById("nav").innerHTML += "<a class='newer-posts' href='?page="+(page - 1)+"'>← Newer Posts</a>"
	}

	document.getElementById("nav").innerHTML += "<span class='page-number'>Page "+page+" of "+count+"</span>"

	if(page < count) {
		document.getElementById("nav").innerHTML += "<a class='older-posts' href='?page="+(page + 1)+"'>← Older Posts</a>"
	}
	
	if (page <= count) {
		var result = get(page)
		for (var i=0;i<result.length;i++) {
			document.getElementById("content").innerHTML += "<article class='post'><header class='post-header'><span class='post-meta'><time datetime='"+result[i].date+"' itemprop='datePublished'>"+result[i].date+"</time><h2 class='post-title'><a href='./html/"+result[i].title+".html'>"+result[i].title+"</a></h2></header><section class='post-excerpt'><p>"+result[i].desc+"</p> <p><a href='./html/"+result[i].title+".html' class='excerpt-link'>Read More...</a></p></section></article>"
		}
	}
}
```

### 联系我

我的博客: [http://blog.csdn.net/qibin0506](http://blog.csdn.net/qibin0506)

我的邮箱: <a href="mailto:qibin0506@gmail.com">qibin0506@gmail.com</a>
