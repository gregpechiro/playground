function saveSettings(settings) {
    var s = JSON.stringify(settings);
    s = btoa(s);
    localStorage.setItem('settings', s);
}

function getSettings() {
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
