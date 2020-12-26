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
            ownerID: null,
            permission: 0,
            clients: [],

            conn: null,
            cm: null,
            client: null,
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
                console.log("closed")
            });

            this.conn.on('doc', (data) => {
                this.cm.setValue(data.document);
                let serverAdapter = new ot.SocketConnectionAdapter(this.conn);
                let editorAdapter = new ot.CodeMirrorAdapter(this.cm);
                this.client = new ot.EditorClient(data.revision, data.clients, serverAdapter, editorAdapter);

                this.setStatus('online')
                this.clients = data.clients
                this.permission = data.permission
            });

            this.conn.on('clients', (data) => {
                this.clients = data;
            });

            this.conn.on('permission', (data) => {
                this.permission = data;
            });

            this.conn.on('registered', (data) => {
                this.userID = data.user_id
                this.ownerID = data.owner_id
                this.cm.setOption('readOnly', false);
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