$(document).ready(function() {

    document.getElementById('import').checked = settings.import;
    document.getElementById('formatOnRun').checked = settings.formatOnRun;
    $('#theme').val(editor.getTheme());
    $('#fontSize').val(editor.getFontSize());
    if (settings.editor.keys === 'vim') {
        $('#keybindings').val('vim');
    } else if (settings.editor.keys === 'emacs'){
        $('#keybindings').val('emacs');
    }
    if (settings.editor.load === 'load') {
        $('#load')[0].checked = true;
    }

    $('#theme').change(function() {
        var theme = $('#theme').val()
        editor.setTheme(theme);
        settings['editor']['theme'] = theme;
        saveSettings(settings);
    });

    $('#fontSize').change(function() {
        var fontSize =+ $('#fontSize').val();
        editor.setFontSize(fontSize);
        settings['editor']['fontSize'] = fontSize;
        saveSettings(settings);
    });

    $('#keybindings').change(function() {
        var bind = $('#keybindings').val();
        if (bind === 'ace') {
            editor.setKeyboardHandler("");
        } else if (bind === 'vim') {
            editor.setKeyboardHandler("ace/keyboard/vim");
        } else {
            editor.setKeyboardHandler("ace/keyboard/emacs");
        }
        settings['editor']['keys'] = bind;
        saveSettings(settings);
    });

    $('#formatOnRun').change(function() {
        settings['formatOnRun'] = this.checked;
        saveSettings(settings);
    });

    $('input[name="load"]').change(function() {
        settings['editor']['load'] = $('#load:checked').val();
        saveSettings(settings);
    });

    $('#reset').click(function() {
        settings = {'editor':{}};
        saveSettings(settings);
        location.reload();
    });

});
