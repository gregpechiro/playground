// wraooer for run
function runWrap() {
    // reset any errors
    editor.getSession().clearAnnotations();
    $('#outpre').css('color', '');
    $('#outpre').text('Processing');
    // get code
    var doc = editor.getValue();
    if (doc !== '') {
        // save code
        saveCode(settings.language, doc);
        if (settings.formatOnRun) {
            // format then run
            format(doc, true);
            return;
        }
        // run without format
        run(doc);
        return;
    }
    $('#outpre').text('');
}

// run
function run(doc) {
    // run post
    $.ajax({
        url: '/run',
        method: 'POST',
        data: {
            dat: doc,
            language:$('select#language').val(),
            version: $('select#version').val()
        },
        success: function(resp) {
            // check for returned error
            $('#outpre').css('color', '');
            if (resp.error) {
                $('#outpre').css('color', 'orangered');
                markLines(resp.output);
            }
            // display output
            $('#outpre').text(resp.output);
        },
        // display server error
        error: function() {
            console.log('run ajax error');
            $('#outpre').css('color', 'orangered');
            $('#outpre').text('Server error. Please try again.');
        }
    });
}

// wrapper for format
function formatWrap() {
    // reset errors
    editor.getSession().clearAnnotations();
    $('#outpre').text('Processing');
    $('#outpre').css('color', '');
    // get code
    var doc = editor.getValue();
    if (doc !== '') {
        // save code
        saveCode(settings.language, doc);
        // format vpde
        format(doc, false);
        return;
    }
    $('#outpre').text('');
}

// format
function format(doc, doRun) {
    // format post
    $.ajax({
        url: '/format',
        method: 'POST',
        data: {
            dat: doc,
            imp: (document.getElementById("import").checked ? 'true' : 'false'),
            language: $('select#language').val(),
            version: $('select#version').val()
        },
        success: function(resp) {
            // check for returned error
            if (resp.error) {
                // display error
                $('#outpre').css('color', 'orangered');
                $('#outpre').text(resp.output);
                markLines(resp.output);
                return;
            }
            // replace editor with formated code
            editor.setValue(resp.output, 1);
            // save code
            saveCode(settings.language, resp.output);
            if (doRun) {
                // run if necessary
                $('#outpre').text('Formated\nProcessing');
                run(resp.output);
                return;
            }
            // print output
            $('#outpre').text('Formated');
            return;

        },
        // display server error
        error: function() {
            console.log('format ajax error');
            $('#outpre').css('color', 'orangered');
            $('#outpre').text('Server error. Please try again.');
            return;
        }
    });
}

$(document).ready(function() {

    // run click
    $('#run').click(runWrap);

    // format click
    $('#format').click(formatWrap);

    // import click
    $('#import').click(function() {
        if (this.checked) {
            settings.import = true;
        } else {
            settings.import = false;
        }
        saveSettings(settings);
    });

    // export click
    $('#export').click(function() {
        var blob = new Blob([editor.getValue()], {type : 'text/plain'});
        $('a#hiddenExport').attr('href', window.URL.createObjectURL(blob));
        $('a#hiddenExport').attr('download', 'main.go');
        $('a#hiddenExport')[0].click();
    });

    // share click
    $('#share').click(function() {
        // reset errors
        $('#outpre').css('color', '');
        $('#outpre').text('');
        //get code
        var doc = editor.getValue();
        // save code
        saveCode(settings.language, doc);
        if (doc !== '') {
            // post share
            $.ajax({
                url: '/share',
                method: 'POST',
                data: {
                    dat: doc
                },
                success: function(resp) {
                    // check for returned error
                    if (resp.error) {
                        // display error
                        $('#outpre').css('color', 'orangered');
                        $('#outpre').text(resp.output);
                        return
                    }
                    // set path in url
                    var hd = {"code": editor.getValue()}
                    var path = "/" + resp.output;
                    window.history.pushState(hd, "", path);
                },
                // display server error
                error: function() {
                    console.log('share ajax error');
                    $('#outpre').css('color', 'orangered');
                    $('#outpre').text('Server error. Please try again.');
                }
            })
        }
    });

    // reset click
    $('#resetCode').click(function() {
        const code = defaultCode[settings.language];
        console.log(code);
        editor.setValue(code, 1);
        saveCode(settings.language, code);
    });

    // favorite click
    $('button#addFavorite').click(function() {
        // get favorite name
        var name = $('input#favoriteName').val();
        if (name === '') {
            // display no name error
            $('#outpre').css('color', 'orangered');
            $('#outpre').text('No name entered for favorite');
            $('#favoriteModal').modal('hide');
            $('#favoriteModal').modal('hide');
            return
        }
        // check settings favorites array
        if (!Array.isArray(settings.favorites)) {
            settings.favorites = [];
        }
        if (settings.favorites.length >= 10) {
            // display max favorites error
            $('#outpre').css('color', 'orangered');
            $('#outpre').text('Maximum number of favorites reached. Please remove some first.');
            $('#favoriteModal').modal('hide');
            return
        }
        // reset errors
        $('#outpre').css('color', '');
        $('#outpre').text('');
        // get code
        var doc = editor.getValue();
        //save code
        saveCode(settings.language, doc);
        if (doc !== '') {
            // share post
            $.ajax({
                url: '/share',
                method: 'POST',
                data: {
                    dat: doc
                },
                success: function(resp) {
                    // check for returned error
                    if (resp.error) {
                        // display error
                        $('#outpre').css('color', 'orangered');
                        $('#outpre').text(resp.output);
                        return
                    }
                    // reset name input
                    $('input#favoriteName').val('');
                    // get share path
                    var path = "/" + resp.output;
                    // double check settings favorites array
                    if (!Array.isArray(settings.favorites)) {
                        settings.favorites = [];
                    }
                    // search existing favorites
                    for (var i = 0; i < settings.favorites.length; i++) {
                        // update favorite if name already exists
                        if (settings.favorites[i].name === name) {
                            settings.favorites[i].path = path;
                            // save favorites
                            saveSettings(settings);
                            // display favorites
                            displayFavorites();
                            return
                        }
                    }
                    // add new favorite
                    settings.favorites.push({'name': name, 'path': path});
                    // save favorites
                    saveSettings(settings);
                    // display favorites
                    displayFavorites();
                },
                // display server error
                error: function() {
                    $('#outpre').css('color', 'orangered');
                    $('#outpre').text('Server error. Please try again.');
                    $('#favoriteModal').modal('hide');
                }
            })
        }
    });

});
