function getCookie(cname) {
    var name = cname + '=';
    var ca = document.cookie.split(';');
    for(var i=0; i<ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1);
        if (c.indexOf(name) == 0) return c.substring(name.length,c.length);
    }
    return "";
}

function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+d.toUTCString();
    document.cookie = cname + '=' + cvalue + '; ' + expires +'; path=/';
}

function delCookie(cname) {
    document.cookie = cname + '=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/';
}

function saveSettings(settings) {
    var s = JSON.stringify(settings);
    s = btoa(s);
    //setCookie('settings', s, 365)


    localStorage.setItem('settings', s);
}

function getSettings() {
    //var settings = getCookie('settings');

    var settings = localStorage.getItem('settings');
    if (settings !== undefined) {
        settings = atob(settings);
        if (settings !== '' && settings !== undefined) {
            return JSON.parse(settings);
        }
        saveSettings({'editor':{}})
    }
    return {'editor':{}};


}
