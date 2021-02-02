new Vue({
    el: "#app",
    delimiters: ['${', '}'],
    data() {
        return {
            loading: true,
            uid: uid,
            url: 'ws://' + location.host + "/socket/" + uid,
            status: 'connecting',
            userID: 0,
            clientID: null,
            ownerID: null,
            permission: 0,
            clients: [],

            markdownHTML: '',

            closed: false,

            conn: null,
            cm: null,
            client: null,
            editorAdapter: null,
        }
    },
    mounted() {
        this.cm = CodeMirror.fromTextArea(document.getElementById('editor'), {
            mode: {
                name: 'gfm',
                tokenTypeOverrides: {
                    emoji: "emoji"
                }
            },
            lineNumbers: true,
            readOnly: 'nocursor',
        });

        this.cm.on('changes', () => {
            let rendered = md.render(this.cm.getValue());
            this.markdownHTML = filterXSS(rendered, filterXSSOptions);
        })

        this.initConnection()
        this.loading = false
    },

    methods: {
        initConnection() {
            this.conn = new SocketConnection(this.url);
            this.conn.on('open', () => {
                this.setStatus('connecting')
                this.conn.send('join', {});
            });

            this.conn.on('close', (evt) => {
                this.setStatus('disconnected')
                // Retry if it's disconnect accidentally.
                if (this.closed) {
                    return
                }
                this.initConnection()
                console.log("closed")
            });

            this.conn.on('error', () => {
                this.closed = true
                window.location.href = '/';
            })

            this.conn.on('doc', (data) => {
                this.cm.setValue(data.document);
                let serverAdapter = new ot.SocketConnectionAdapter(this.conn);
                this.editorAdapter = new ot.CodeMirrorAdapter(this.cm);
                this.client = new ot.EditorClient(data.revision, data.clients, serverAdapter, this.editorAdapter);

                this.setStatus('online')
                this.ownerID = data.owner_id
                this.clients = data.clients
                this.permission = data.permission
            });

            this.conn.on('clients', (data) => {
                this.clients = data;
            });

            this.conn.on('permission', (data) => {
                // Check permission
                // Guest View: 0 1 3 Edit: 0
                // User View: 0 1 2 3 4 Edit: 0 1 2
                let canView = false
                let canEdit = false

                if (this.userID !== this.ownerID) {
                    if (this.userID === 0) { // Guest
                        switch (data) {
                            case 0:
                                canEdit = true
                            case 1:
                            case 3:
                                canView = true
                        }
                    } else { // User
                        switch (data) {
                            case 0:
                            case 1:
                            case 2:
                                canEdit = true
                            case 3:
                            case 4:
                                canView = true
                        }
                    }
                } else {
                    canView = true
                    canEdit = true
                }

                this.editorAdapter.canEdit = canEdit;
                this.cm.setOption('readOnly', !canEdit);

                this.permission = data;
            });

            this.conn.on('registered', (data) => {
                this.userID = data.user_id
                this.clientID = data.client_id

                this.editorAdapter.canEdit = !data.read_only;
                this.cm.setOption('readOnly', data.read_only);
            });

            this.conn.on('join', (data) => {
                console.log(data);
            });

            this.conn.on('quit', (data) => {
                console.log(data);
            });
        },

        // Remove the repeat clients.
        getUsers() {
            let u = {}
            let users = []

            this.clients.forEach((client) => {
                // Every guest are different.
                if (u[client.user_id] === undefined || client.user_id === 0) {
                    users.push(client)
                    u[client.user_id] = true
                }
            })
            return users
        },

        setPermission(permission) {
            this.conn.send('permission', permission)
        },

        setStatus(status) {
            this.status = status
        },
    }
})