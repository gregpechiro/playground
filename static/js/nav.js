$(document).ready(function() {
    $('#run').click(function() {
        editor.getSession().clearAnnotations();
        $('#outpre').css('color', '');
        $('#outpre').text('Processing');
        var doc = editor.getValue();
        settings.code = doc;
        saveSettings(settings);
        if (doc !== '') {
            $.ajax({
                url: '/run',
                method: 'POST',
                data: {
                    dat: doc
                },
                success: function(resp) {
                    console.log('ajax success');
                    if (!resp.error) {
                        $('#outpre').css('color', '');
                    } else {
                        $('#outpre').css('color', 'orangered');
                        markLines(resp.output);
                    }
                    $('#outpre').text(resp.output);
                },
                error: function() {
                    console.log('ajax error');
                    $('#outpre').css('color', 'orangered');
                    $('#outpre').text('Server error. Please try again.');
                }
            })
        }
    });

    $('#format').click(function() {
        editor.getSession().clearAnnotations();
        $('#outpre').text('Processing');
        $('#outpre').css('color', '');
        var doc = editor.getValue();
        settings.code = doc;
        saveSettings(settings);
        if (doc !== '') {
            $.ajax({
                url: '/format',
                method: 'POST',
                data: {
                    dat: doc,
                    imp: (document.getElementById("import").checked ? 'true' : 'false')
                },
                success: function(resp) {
                    console.log('ajax success');
                    console.log(resp);
                    if (!resp.error) {
                        inputChanged();
                        editor.setValue(resp.output, 1);
                        settings.code = resp.output;
                        saveSettings(settings);
                        $('#outpre').text('Formated');
                        return
                    }
                    $('#outpre').css('color', 'orangered');
                    $('#outpre').text(resp.output);
                },
                error: function() {
                    console.log('ajax error');
                    $('#outpre').css('color', 'orangered');
                    $('#outpre').text('Server error. Please try again.');
                }
            })
        }
        $('#outpre').text('');
    });

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
                    console.log('ajax success');
                    if (resp.error) {
                        $('#outpre').css('color', 'orangered');
                        $('#outpre').text(resp.output);
                        return
                    }
                    var hd = {"code": editor.getValue()}
                    path = "/" + resp.output;
                    //var url = origin(window.location) + path;
                    //shareURL.show().val(url).focus().select();
                    window.history.pushState(hd, "", path);
                },
                error: function() {
                    console.log('ajax error');
                    $('#outpre').css('color', 'orangered');
                    $('#outpre').text('Server error. Please try again.');
                }
            })
        }
    });

});
