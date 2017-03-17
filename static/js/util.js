// save settings to local storage
function saveSettings(settings) {
    var s = JSON.stringify(settings);
    s = btoa(unescape(encodeURIComponent(s)));
    localStorage.setItem('settings', s);
}

// return saved settings from local storage
function getSettings() {
    var s = localStorage.getItem('settings');
    if (s !== undefined && s !== null && s !== '') {
        s = decodeURIComponent(escape(atob(s)));
        if (s !== '' && s !== undefined && s[0] === '{') {
            return JSON.parse(s);
        }
        saveSettings({'editor':{}})
    }
    return {'editor':{}};
}

// mark errors on lines
function markLines(output) {
    // divied error into array by \n
    var outputList = output.split('\n');
    var annotations = [];
    // loop output errors
    for (var i = 0; i < outputList.length; i++) {
        if (outputList[i] !== '') {
            // get line number
            var lineList = outputList[i].split(':');
            var lineNum =+ lineList[0].replace('Line ', '');
            // create pretty error message
            var msg = lineList.slice(1, lineList.length).join(':');
            // add annotation to annotation array
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
    // add annotations to ace
    editor.getSession().setAnnotations(annotations);
}

// display favorites modal links
function displayFavorites() {
    // reset favorites links
    $('#favoriteLinks').html('');
    // add space above links
    var links = '<div class="form-group"></div>';
    // loop all settings favorites
    for (var i = 0; i < settings.favorites.length; i++) {
        // create link button and delete button for every link
        links += '<div class="form-group clearfix">'+
            '<div class="col-lg-10">' +
                '<a href="' + settings.favorites[i].path + '" class="btn btn-default btn-block">' + settings.favorites[i].name + '</a>'+
            '</div>'+
            '<div class="col-lg-2">'+
                '<button class="btn btn-danger btn-block remove-favorite" data-idx="' + i + '">'+
                    '<i class="fa fa-trash"></i>'+
                '</button>'+
            '</div>'+
        '</div>';
    }
    // set favorite links
    $('#favoriteLinks').html(links);
}
