$(document).ready(function() {

    function run(doc) {
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
        });
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
                console.log('ajax success');
                console.log(resp);
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
                return
            },
            error: function() {
                console.log('ajax error');
                $('#outpre').css('color', 'orangered');
                $('#outpre').text('Server error. Please try again.');
                return
            }
        });
    }

    $('#run').click(function() {
        // editor.getSession().clearAnnotations();
        // $('#outpre').css('color', '');
        // $('#outpre').text('Processing');
        // var doc = editor.getValue();
        // settings.code = doc;
        // saveSettings(settings);
        // if (doc !== '') {
        //     $.ajax({
        //         url: '/run',
        //         method: 'POST',
        //         data: {
        //             dat: doc
        //         },
        //         success: function(resp) {
        //             console.log('ajax success');
        //             if (!resp.error) {
        //                 $('#outpre').css('color', '');
        //             } else {
        //                 $('#outpre').css('color', 'orangered');
        //                 markLines(resp.output);
        //             }
        //             $('#outpre').text(resp.output);
        //         },
        //         error: function() {
        //             console.log('ajax error');
        //             $('#outpre').css('color', 'orangered');
        //             $('#outpre').text('Server error. Please try again.');
        //         }
        //     })
        // }
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
        }
        $('#outpre').text('');
    });

    $('#format').click(function() {
        // editor.getSession().clearAnnotations();
        // $('#outpre').text('Processing');
        // $('#outpre').css('color', '');
        // var doc = editor.getValue();
        // settings.code = doc;
        // saveSettings(settings);
        // if (doc !== '') {
        //     $.ajax({
        //         url: '/format',
        //         method: 'POST',
        //         data: {
        //             dat: doc,
        //             imp: (document.getElementById("import").checked ? 'true' : 'false')
        //         },
        //         success: function(resp) {
        //             console.log('ajax success');
        //             console.log(resp);
        //             test = resp;
        //             if (!resp.error) {
        //                 inputChanged();
        //                 editor.setValue(resp.output, 1);
        //                 settings.code = resp.output;
        //                 saveSettings(settings);
        //                 $('#outpre').text('Formated');
        //                 return
        //             }
        //             $('#outpre').css('color', 'orangered');
        //             $('#outpre').text(resp.output);
        //             return
        //         },
        //         error: function() {
        //             console.log('ajax error');
        //             $('#outpre').css('color', 'orangered');
        //             $('#outpre').text('Server error. Please try again.');
        //             return
        //         }
        //     })
        // }
        // $('#outpre').text('');
        editor.getSession().clearAnnotations();
        $('#outpre').text('Processing');
        $('#outpre').css('color', '');
        var doc = editor.getValue();
        if (doc !== '') {
            settings.code = doc;
            saveSettings(settings);
            format(doc, false);
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
