// global variables
var settings = {};
var editor = ace.edit("editor");

// change path if input os changed
function inputChanged() {
    if ((window.location.pathname == "/")) {
        return;
    }
    window.history.pushState(null, "", "/");
}

$(document).ready(function() {
    // setup ace editor
    editor.session.setMode("ace/mode/golang");
    editor.renderer.setShowGutter(true);
    editor.setHighlightActiveLine(true);
    editor.setReadOnly(false);
    editor.setDisplayIndentGuides(true);
    editor.setOption("scrollPastEnd", 1);
    editor.renderer.setScrollMargin(15)
    editor.$blockScrolling = Infinity
    editor.setOptions({
        enableBasicAutocompletion: true,
        enableSnippets: true,
        enableLiveAutocompletion: true
    });
    $('textarea.ace_text-input').focus();

    // init settings from local storage
    settings = getSettings();
    if (settings.editor !== undefined && !$.isEmptyObject(settings.editor)) {
        editor.setTheme(((settings.editor.theme === '' || settings.editor.theme === undefined) ? 'ace/theme/monokai' : settings.editor.theme));
        editor.setFontSize(((settings.editor.fontSize === '' || settings.editor.fontSize === undefined) ? 12 : settings.editor.fontSize));
        if (settings.editor.keys === 'vim') {
            editor.setKeyboardHandler("ace/keyboard/vim");
        } else if (settings.editor.keys === 'emacs'){
            editor.setKeyboardHandler("ace/keyboard/emacs");
        }
    } else {
        editor.setTheme('ace/theme/monokai');
        editor.setFontSize(12);
    }

    if (!Array.isArray(settings.favorites)) {
        settings.favorites = [];
    }
    if (settings.favorites.length > 0) {
        displayFavorites();
    }

    // add key bindings to ace editor
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

    editor.commands.addCommand({
        name: "replace",
        bindKey: {win: "Ctrl-R", mac: "Command-Option-R"},
        exec: function(editor) {
            ace.config.loadModule("ace/ext/searchbox", function(e) {e.Search(editor, true)});
        }
    });

    // trigger path reset on update
    editor.on('change', function() {
        inputChanged();
    });

    // add window key bindings

    function onKeyDown(e) {
        // ctrl+r remove default
        if (e.ctrlKey) { // ctrl
            if (e.keyCode == 82) { // +r
                e.preventDefault();
                return
            }
        }

        // ctrl+shift+s run sizeof
        if (e.ctrlKey) { // ctrl
            if (e.shiftKey) { // + shift
                if (e.keyCode == 83) { // +s
                    e.preventDefault();
                    $('#flipper').removeClass('flip');
                    $('#flipTrigger').text('Settings');
                    sizeof();
                    return
                }
            }
        }

        if (e.keyCode == 13) { // enter
            // enter + shift run
            if (e.shiftKey) { // +shift
                e.preventDefault();
                runWrap();
                return
            }
            // enter + ctrl format
            if (e.ctrlKey) { // +ctrl
                e.preventDefault();
                formatWrap();
                return
            }
        }
    }
    // register the handler
    document.addEventListener('keydown', onKeyDown, false);

    // add tips button
    $('#editor').prepend('<a class="btn control" data-toggle="modal" data-target="#tipsModal"><i class="fa fa-lg fa-question"></i></a>');

    // flipper
    $('#flipTrigger').click(function() {
        $('#flipper').toggleClass('flip');
        if (this.innerText === 'Settings') {
            this.innerText = 'Sizeof';
            return
        }
        this.innerText = 'Settings';
        return
    });
});
