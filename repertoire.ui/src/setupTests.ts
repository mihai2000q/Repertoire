import '@testing-library/jest-dom/vitest'
import { WebSocket } from 'ws'
import dayjs from 'dayjs'
import isToday from 'dayjs/plugin/isToday'
import isYesterday from 'dayjs/plugin/isYesterday'
import 'dayjs/locale/en-gb'

global.WebSocket = WebSocket as never

dayjs.extend(isToday)
dayjs.extend(isYesterday)
dayjs.locale('en-gb')
