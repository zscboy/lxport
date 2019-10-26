"use strict";
const net = require('net');
const WebSocket = require('ws');
const defaultForward = 22;

const websocketURL = "wss://host/path";
let portForward;
let localPort = 8005;

function serveConn(socket) {
    let url = websocketURL + "?port=" + portForward;
    const ws = new WebSocket(url);

    ws.on('open', function open() {
        socket.on('data', (data) => {
            ws.send(data);
        });
    });

    ws.on('message', (data) => {
        socket.write(data);
    });

    ws.on('close', () => {
        console.log('websocket close');
        socket.destroy();
    });

    ws.on('error', (err) => {
        console.log('websocket error:', err);
    });

    socket.on('close', function () {
        console.log('socket close');
        ws.close();
    });

    socket.on('error', function (err) {
        console.log('socket error:', err);
        ws.close();
    });
}

function main() {
    let args = require('yargs')
        .option('loc', {
            alias: 'l',
            describe: 'provide local port to bind'
        })
        .option('port', {
            alias: 'p',
            describe: 'provide a port to forward'
        })
        .help()
        .argv;

    portForward = args['p'];
	if (portForward === undefined) {
		portForward = defaultForward;
	}

    if (args['l'] !== undefined) {
        localPort = args['l'];
    }

    let server = net.createServer(socket => {
        serveConn(socket);
    });

    console.log("forward port:", portForward);
    console.log("server listen at:", localPort);
    server.listen(localPort).on("error", (err) => {
        console.log("listen failed:", err)
    });

}

main();

