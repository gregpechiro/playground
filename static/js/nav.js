
function runWrap() {
    editor.getSession().clearAnnotations();
    $('#outpre').css('color', '');
    $('#outpre').text('Processing');
    var doc = editor.getValue();
    if (doc !== '') {
        settings.code = doc;
        saveSettings(settings);
        if (settings.formatOnRun) {
            format(doc, true);
            return
        }
        run(doc);
        return
    }
    $('#outpre').text('');
}

function run(doc) {
    $.ajax({
        url: '/run',
        method: 'POST',
        data: {
            dat: doc
        },
        success: function(resp) {
            if (!resp.error) {
                $('#outpre').css('color', '');
            } else {
                $('#outpre').css('color', 'orangered');
                markLines(resp.output);
            }
            $('#outpre').text(resp.output);
        },
        error: function() {
            console.log('run ajax error');
            $('#outpre').css('color', 'orangered');
            $('#outpre').text('Server error. Please try again.');
        }
    });
}

function formatWrap() {
    editor.getSession().clearAnnotations();
    $('#outpre').text('Processing');
    $('#outpre').css('color', '');
    var doc = editor.getValue();
    if (doc !== '') {
        settings.code = doc;
        saveSettings(settings);
        format(doc, false);
        return
    }
    $('#outpre').text('');
}

function format(doc, doRun) {
    $.ajax({
        url: '/format',
        method: 'POST',
        data: {
            dat: doc,
            imp: (document.getElementById("import").checked ? 'true' : 'false')
        },
        success: function(resp) {
            if (!resp.error) {
                inputChanged();
                editor.setValue(resp.output, 1);
                settings.code = resp.output;
                saveSettings(settings);
                if (doRun) {
                    $('#outpre').text('Formated\nProcessing');
                    run(resp.output);
                    return
                }
                $('#outpre').text('Formated');
                return
            }
            $('#outpre').css('color', 'orangered');
            $('#outpre').text(resp.output);
            markLines(resp.output);
            return
        },
        error: function() {
            console.log('format ajax error');
            $('#outpre').css('color', 'orangered');
            $('#outpre').text('Server error. Please try again.');
            return
        }
    });
}

$(document).ready(function() {

    $('#run').click(runWrap);

    $('#format').click(formatWrap);

    $('#import').click(function() {
        if (this.checked) {
            settings.import = true;
        } else {
            settings.import = false;
        }
        saveSettings(settings);
    });

    $('#export').click(function() {
        var blob = new Blob([editor.getValue()], {type : 'text/plain'});
        $('a#hiddenExport').attr('href', window.URL.createObjectURL(blob));
        $('a#hiddenExport').attr('download', 'main.go');
        $('a#hiddenExport')[0].click();
    });

    $('#share').click(function() {
        $('#outpre').css('color', '');
        $('#outpre').text('');
        var doc = editor.getValue();
        settings.code = doc;
        saveSettings(settings);
        if (doc !== '') {
            $.ajax({
                url: '/share',
                method: 'POST',
                data: {
                    dat: doc
                },
                success: function(resp) {
                    if (resp.error) {
                        $('#outpre').css('color', 'orangered');
                        $('#outpre').text(resp.output);
                        return
                    }
                    var hd = {"code": editor.getValue()}
                    path = "/" + resp.output;
                    window.history.pushState(hd, "", path);
                },
                error: function() {
                    console.log('share ajax error');
                    $('#outpre').css('color', 'orangered');
                    $('#outpre').text('Server error. Please try again.');
                }
            })
        }
    });

    $('#resetCode').click(function() {
        editor.setValue('package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello, custom playground")\n}', 1);
        settings.code = '';
        saveSettings(settings);
    });

});
