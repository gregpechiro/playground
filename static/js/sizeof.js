$(document).ready(function() {
    $('button#getSize').click(sizeof);
});

// size function
function sizeof() {
    // get selected code
    var tempCode = editor.getSelectedText();
    if (tempCode === '') {
        // display error
        $('#sizeofResults').html('<div class="bs-callout bs-callout-danger"><h4>Error</h4><p>No selected code.</p></div>');
        return
    }
    // sizeof post
    $.ajax({
        url: "/sizeof",
        method: 'POST',
        data: {
            code: tempCode
        },
        success: function(resp) {
            // check for returned error
            if (resp.error) {
                // display error
                $('#sizeofResults').html('<div class="bs-callout bs-callout-danger"><h4>Error</h4><p>' + resp.output + '</p></div>');
                return
            }
            // display sizeof results
            $('#sizeofResults').html(resp.output);
        },
        // display server error
        error: function() {
            console.log('sizeof ajax error');
            $('#sizeofResults').html('<div class="bs-callout bs-callout-danger"><h4>Server Error</h4><p>Please try again</p></div>');
        }
    });
}
