function saveSettings(settings) {
    var s = JSON.stringify(settings);
    s = btoa(s);
    localStorage.setItem('settings', s);
}

function getSettings() {
    var s = localStorage.getItem('settings');
    if (s !== undefined && s !== null && s !== '') {
        s = atob(s);
        if (s !== '' && s !== undefined && s[0] === '{') {
            return JSON.parse(s);
        }
        saveSettings({'editor':{}})
    }
    return {'editor':{}};

}
