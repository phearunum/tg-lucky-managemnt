<template>
  <div class="app-container">
    <el-card class="terminal-container">
      <div class="terminal-output" ref="terminalOutput" @click="focusInput">
        <pre>{{ output }}</pre>
        <span class="prompt">{{ prompt }}</span>
        <span contenteditable="true" class="input" ref="input" @input="updateCommand" @keydown="handleKeydown"></span>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'

const command = ref('')
const output = ref('')
const prompt = ref('phearunum@Picha-Mac-mini ~ % ')
const commandHistory = ref([]) // Store command history
let historyIndex = ref(-1) // Track current index in history

// Initial prompt and last login message
const lastLogin = `Last login: Fri Nov  1 15:53:18 on ttys008
You have new mail.`
output.value = `${lastLogin}\n${prompt.value}` // Add a newline before the prompt

let socket = null

// Establish a WebSocket connection to the server
const connectWebSocket = () => {
  socket = new WebSocket('ws://localhost:1212/terminal')

  socket.onopen = () => {
    output.value += '\n[Connected to server]\n' // Connection confirmation with newline
  }

  socket.onmessage = (event) => {
    output.value += `\n${event.data}` // Append command output with leading new line
    output.value += `\n${prompt.value}` // Re-add the prompt after output
    scrollToBottom() // Scroll to bottom
  }

  socket.onerror = (error) => {
    console.error('WebSocket error:', error)
    output.value += '\n[Error connecting to server]' // Display error with newline
  }

  socket.onclose = () => {
    output.value += '\n[Disconnected from server]' // Notify disconnection with newline
  }
}

// Update the command based on user input
const updateCommand = (event) => {
  command.value = event.target.innerText // Update command with current input
}

// Handle keyboard input
const handleKeydown = (event) => {
  // Navigate command history with up and down arrows
  if (event.key === 'ArrowUp') {
    if (historyIndex.value < commandHistory.value.length - 1) {
      historyIndex.value++
      command.value = commandHistory.value[commandHistory.value.length - 1 - historyIndex.value] || ''
      clearInput()
      return
    }
  } else if (event.key === 'ArrowDown') {
    if (historyIndex.value > 0) {
      historyIndex.value--
      command.value = commandHistory.value[commandHistory.value.length - 1 - historyIndex.value] || ''
      clearInput()
      return
    }
  }

  if (event.key === 'Enter') {
    sendCommand()
  }
}

// Send the command to the server and update output
const sendCommand = () => {
  if (command.value.trim() === '') return // Avoid sending empty commands

  // Append the command to output with the prompt
  output.value += `\n${prompt.value}${command.value}` // Show the command in output

  // Send command to server
  socket.send(command.value)

  // Add command to history
  commandHistory.value.push(command.value) // Add command to history
  historyIndex.value = 0 // Reset history index
  command.value = '' // Clear command
  clearInput() // Clear input area

  scrollToBottom() // Scroll to bottom after adding command
}

// Clear the input area after sending the command
const clearInput = () => {
  const inputElement = $refs.input
  inputElement.innerText = '' // Clear the content of the editable div
  inputElement.focus() // Keep focus on input for immediate typing
}

// Scroll to the bottom of the terminal output
const scrollToBottom = () => {
  const terminalOutput = $refs.terminalOutput
  terminalOutput.scrollTop = terminalOutput.scrollHeight
}

// Focus on the input area when the terminal is clicked
const focusInput = () => {
  const inputElement = $refs.input
  inputElement.focus()
}

// Cleanup on component destruction
onBeforeUnmount(() => {
  if (socket) {
    socket.close() // Close WebSocket on component destruction
  }
})

// Initialize WebSocket connection when component is mounted
onMounted(() => {
  connectWebSocket()
})
</script>

<style scoped>
.terminal-container {
  max-width: 100%;
  height: 100vh;
  margin: 5px auto;
  padding: 2px;
  background-color: #1e1e1e;
  color: #0f0;
  font-family: monospace;
}

.terminal-output {
  background-color: #1e1e1e;
  color: #0f0;
  height: 100vh;
  overflow-y: auto;
  padding: 10px;
  border-radius: 5px;
  margin-bottom: 10px;
  white-space: pre-wrap; /* Keep text formatting */
  position: relative;
}

.prompt {
  display: inline; /* Keep the prompt inline with user input */
}

.input {
  outline: none; /* Remove default outline */
  caret-color: #0f0; /* Green caret */
}
</style>
