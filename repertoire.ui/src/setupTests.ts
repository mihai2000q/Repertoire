import '@testing-library/jest-dom/vitest'
import { WebSocket } from 'ws'

global.WebSocket = WebSocket as never
