// 更新 url ，但仅替换 args ，如果不存在则添加
function updateLocation(key, value) {
    const args = [];
    if (window.location.search.length === 0) {
        args.push(key + '=' + value)
    } else {
        let exists = false
        let searches = window.location.search.substring(1).split("&");
        for (let search of searches) {
            if (search.startsWith(key + '=')) {
                exists = true
                if (value !==""){
                    args.push(key + '=' + value)
                }
            } else {
                    args.push(search);
            }
        }
        if (!exists && value !=="") {
            args.push(key + '=' + value)
        }
    }
    let url = window.location.pathname + '?' + args.join("&");
    if (window.location.hash.length > 0) {
        url = url + '#' + window.location.hash
    }
    window.history.replaceState({}, '', url);
}

function getOrSetUrlParam(key,def) {
    let result = getUrlParam(key,def);
    if (result !== ""){
        updateLocation(key,result)
    }
    return result
}


function getUrlParam(key,def) {
    for (let param of window.location.search.substring(1).split("&")) {
        if (param.startsWith(key + '=')) {
            return param.split('=', 2)[1]
        }
    }
    return def
}