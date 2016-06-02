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

function markLines(output) {
    var outputList = output.split('\n');
    var annotations = [];
    for (var i = 0; i < outputList.length; i++) {
        if (outputList[i] !== '') {
            var lineList = outputList[i].split(':');
            var lineNum =+ lineList[0].replace('Line ', '');
            var msg = lineList.slice(1, lineList.length).join(':');
            annotations.push(
                {
                    row: (lineNum - 1),
                    column: 0,
                    text: msg,
                    type: "error"
                }
            );
        }
    }
    editor.getSession().setAnnotations(annotations);
}
