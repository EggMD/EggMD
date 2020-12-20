(function () {
    'use strict';

    function setStatus(status) {
        App.statusBadge.connecting.hide();
        App.statusBadge.online.hide();
        App.statusBadge.disconnected.hide();

        switch (status) {
            case 'online':
                App.statusBadge.online.show();
                return;
            case 'connecting':
                App.statusBadge.connecting.show();
                return;
            case 'disconnected':
                App.statusBadge.disconnected.show();
                return;
        }
    }

    function setClientCount(count) {
        App.statusBadge.online.text(count.toString() + ' ONLINE');
    }

    window.App = {
        conn: null,
        cm: null,
        statusBadge: {
            connecting: $('#connectingBadge'),
            online: $('#onlineBadge'),
            disconnected: $('#disconnectedBadge'),
        },
        markdown: $('#markdown')
    };

    setStatus('empty')

    App.cm = CodeMirror.fromTextArea(document.getElementById('editor'), {
        mode: {
            name: 'gfm',
            tokenTypeOverrides: {
                emoji: "emoji"
            }
        },
        lineNumbers: true,
        readOnly: 'nocursor',
    });


    var refresh = _.debounce(function () {
        updateView()
    }, 100)

    App.cm.on('change', function () {
        refresh()
    })

    let url = 'ws://' + location.host + "/socket/" + uid
    var conn = App.conn = new SocketConnection(url);

    conn.on('open', function () {
        setStatus('connecting')
        App.conn.send('join', {});
    });

    conn.on('close', function (evt) {
        setStatus('disconnected')
        console.log("closed")
    });

    conn.on('doc', function (data) {
        App.cm.setValue(data.document);
        var serverAdapter = new ot.SocketConnectionAdapter(conn);
        var editorAdapter = new ot.CodeMirrorAdapter(App.cm);
        App.client = new ot.EditorClient(data.revision, data.clients, serverAdapter, editorAdapter);

        setStatus('online')
        setClientCount(data.clients.length);
    });

    conn.on('clients', function (data) {
        setClientCount(data.length);
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