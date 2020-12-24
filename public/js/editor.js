new Vue({
    el: "#app",
    delimiters: ['${', '}'],
    data() {
        return {
            loading: true,
            uid: '33c2122f-9cf4-43af-8567-070a86b52d1c',
            url: 'ws://' + location.host + "/socket/" + uid,
            status: 'connecting',
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

            this.conn.on('registered', (clientId) => {
                this.cm.setOption('readOnly', false);
            });

            this.conn.on('join', (data) => {
                console.log(data);
            });

            this.conn.on('quit', (data) => {
                console.log(data);
            });
        },

        setPermission(permission) {
            this.conn.send('permission', permission)
        },

        setStatus(status) {
            this.status = status
        },
    }
})