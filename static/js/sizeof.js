$(document).ready(function() {
    $('button#getSize').click(sizeof);
});

function sizeof() {
    var tempCode = editor.getSelectedText();
    if (tempCode === '') {
        $('#sizeofResults').html('<div class="bs-callout bs-callout-danger"><h4>Error</h4><p>No selected code.</p></div>');
        return
    }
    $.ajax({
        url: "/sizeof",
        method: 'POST',
        data: {
            code: tempCode
        },
        success: function(resp) {
            if (resp.error) {
                // var div = $('<div class="bs-callout bs-callout-danger">' + resp.output + '</div>')
                $('#sizeofResults').html('<div class="bs-callout bs-callout-danger"><h4>Error</h4><p>' + resp.output + '</p></div>');
                return
            }
            $('#sizeofResults').html(resp.output);
        },
        error: function() {
            console.log('sizeof ajax error');
            $('#sizeofResults').html('<div class="bs-callout bs-callout-danger"><h4>Server Error</h4><p>Please try again</p></div>');
        }
    });
}
