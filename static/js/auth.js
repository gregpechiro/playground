
function login(s) {
    if (s !== 'login please') {
        return;
    }
    $('<div class="modal fade" id="myModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">'+
        '<div class="modal-dialog modal-sm" role="document">' +
            '<div class="modal-content">' +
                '<div class="modal-header">' +
                    '<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>' +
                    '<h4 class="modal-title" id="myModalLabel">Login</h4>' +
                '</div>' +
                '<div class="modal-body">' +
                    '<form action="/login" method="post">' +
                        '<label>Username</label>' +
                        '<div class="form-group">' +
                            '<input name="username" class="form-control" type="text">' +
                        '</div>' +
                        '<label>Password</label>' +
                        '<div class="form-group">' +
                            '<input name="password" class="form-control" type="password">' +
                        '</div>' +
                        '<button class="btn btn-primary btn-block">Login</button>' +
                    '</form>' +
                '</div>'+
            '</div>' +
        '</div>' +
    '</div>').modal('show');
}

function logout() {
    window.location.href = '/logout';
}
