(function () {
    'use strict';

    window.App = {
        conn: null,
        cm: null
    };

    App.cm = CodeMirror.fromTextArea(document.getElementById('editor'), {
        lineNumbers: true,
        readOnly: 'nocursor',
    });


    let url = 'ws://' + location.host + "/socket/" + shortID
    var conn = App.conn = new SocketConnection(url);

    conn.on('open', function () {
        App.conn.send('join', {});
    });

    conn.on('close', function (evt) {
        console.log("closed")
    });

    conn.on('doc', function (data) {
        App.cm.setValue(data.document);
        var serverAdapter = new ot.SocketConnectionAdapter(conn);
        var editorAdapter = new ot.CodeMirrorAdapter(App.cm);
        App.client = new ot.EditorClient(data.revision, data.clients, serverAdapter, editorAdapter);
    });

    conn.on('registered', function (clientId) {
        App.cm.setOption('readOnly', false);
    });

    conn.on('join', function (data) {
        console.log(data);
    });

    conn.on('quit', function (data) {
        console.log(data);
    });
}());