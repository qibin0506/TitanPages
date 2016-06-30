
var defaultPageSize = 5
var arr = eval('[{"title":"我的第一篇博客","date":"2016-07-01 00:39:19","desc":"这是我的第一篇博客, 主要是用来做测试的, `TitanPages`是不是可以完美的生成静态博客呢?"}]')
function get(currentPage) {
	return getResult(currentPage, defaultPageSize);
}
function getResult(currentPage, pageSize) {
	currentPage = parseInt(currentPage);
	pageSize = parseInt(pageSize);
	var startIndex = (currentPage - 1) * pageSize;
	var endIndex = startIndex + pageSize;
	if (arr.length <= startIndex) { return null;}
	if (endIndex > arr.length) { endIndex = arr.length;}
	return arr.slice(startIndex, endIndex);
}
function pageCount() {
	return getPageCount(defaultPageSize);
}
function getPageCount(pageSize) {
	return Math.ceil(arr.length / pageSize);
}
function getQueryString(query) {
	var uri = window.location.search;
    var re = new RegExp("" +query+ "=([^&?]*)", "ig");
    return ((uri.match(re))?(uri.match(re)[0].substr(query.length+1)):null);
}
	