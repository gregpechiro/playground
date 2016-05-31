
var settings = {};
var editor = ace.edit("editor");
$(document).ready(function() {
    editor.session.setMode("ace/mode/golang");
    editor.renderer.setShowGutter(true);
    editor.setHighlightActiveLine(true);
    editor.setReadOnly(false);
    editor.setDisplayIndentGuides(true);
    editor.$blockScrolling = Infinity
    $('textarea.ace_text-input').focus();

    settings = getSettings();
    if (settings.editor !== undefined && !$.isEmptyObject(settings.editor)) {
        editor.setTheme((settings.editor.theme == '' ? 'ace/theme/monokai' : settings.editor.theme));
        editor.setFontSize((settings.editor.fontSize == '' ? 12 : settings.editor.fontSize));
    } else {
        editor.setTheme('ace/theme/monokai');
        editor.setFontSize(12);
    }

    editor.commands.addCommand({
        name: 'moveLineDown',
        bindKey: {win: 'Ctrl-down',  mac: 'Command-down'},
        exec: function(editor) {
            editor.moveLinesDown();
        },
        readOnly: false
    });

    editor.commands.addCommand({
        name: 'moveLineUp',
        bindKey: {win: 'Ctrl-up',  mac: 'Command-up'},
        exec: function(editor) {
            editor.moveLinesUp();
        },
        readOnly: false
    });

    editor.commands.addCommand({
        name: 'moveCursorUp',
        bindKey: {win: 'Alt-Shift-up',  mac: 'Alt-Shift-up'},
        exec: function(editor) {
            editor.selectMoreLines(-1);
        },
        readOnly: false
    });

    editor.commands.addCommand({
        name: 'moveCursorDown',
        bindKey: {win: 'Alt-Shift-down',  mac: 'Alt-Shift-down'},
        exec: function(editor) {
            editor.selectMoreLines(1);
        },
        readOnly: false
    });
});
