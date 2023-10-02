const fs = require('fs')
const net = require('net')
const { pipeline } = require('stream')

const fileName = "client_1"

const fileStream = fs.createReadStream('./nums.txt')
const socket = new net.Socket()

class ConnectionState {
  static INITIALIZED = 0
  static SENT_REQUEST = 1

  constructor() {
    this.state = ConnectionState.INITIALIZED
  }

  setSentRequest = () => {
    this.state = ConnectionState.SENT_REQUEST
  }
}

socket.connect(8000, 'localhost', () => {
  console.log('connected')

  const connection = new ConnectionState()

  socket.on('data', (data) => processData(connection, data))
  socket.on('error', console.error)

  socket.write('x')
  const buffer = Buffer.alloc(16)
  buffer.write(fileName)
  socket.write(buffer)

  // assume we're going to write 2KB (not allowed)
  const buffer2 = Buffer.alloc(4)
  buffer2.writeUInt32BE(255)
  socket.write(buffer2)

  pipeline(fileStream, socket, () => socket.destroy())
})

const processData = (conn, data) => {
  if (conn.state = ConnectionState.SENT_REQUEST) {
    const res = data.readInt8(0)
    if (res == 1) {
      console.log('liked the request')
    } else {
      console.log("didn't like request")
    }
  }
}
