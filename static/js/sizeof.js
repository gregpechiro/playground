$(document).ready(function() {
    $('button#getSize').click(sizeof);
});

function sizeof() {
    var tempCode = editor.getSelectedText();
    if (tempCode === '') {
        $('#sizeofResults').css('color', 'orangered');
        $('#sizeofResults').html('No selected code.');
        return
    }
    $.ajax({
        url: "/sizeof",
        method: 'POST',
        data: {
            code: tempCode
        },
        success: function(resp) {
            if (!resp.error) {
                $('#sizeofResults').css('color', '');
            } else {
                $('#sizeofResults').css('color', 'orangered');
            }
            $('#sizeofResults').html(resp.output);
        },
        error: function() {
            console.log('ajax error');
            $('#sizeofResults').css('color', 'orangered');
            $('#sizeofResults').html('Server error. Please try again.');
        }
    });
}
