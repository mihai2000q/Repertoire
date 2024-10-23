import axios from 'axios'
import { Mutex } from 'async-mutex'
import { useNavigate } from 'react-router-dom'

export const api = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL
})

// Authorization
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  config.headers.Authorization = `Bearer ${token}`
  return config
})

// Refresh Token
const mutex = new Mutex()
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    await mutex.waitForUnlock()
    console.log(error)
    if (error.response?.status === 401) {
      if (!mutex.isLocked()) {
        const release = await mutex.acquire()
        try {
          const token = localStorage.getItem('token')
          const refresh = await axios.put('/auth/refresh/', { token })
          const data = refresh.data as { token: string } | undefined
          if (data) {
            localStorage.setItem('token', token)
            error = await axios.request(error.config)
          } else {
            localStorage.removeItem('token')
          }
        } finally {
          release()
        }
      } else {
        await mutex.waitForUnlock()
        error = await axios.request(error.config)
      }
    }
    return error
  }
)

const errorCodeToPathname = new Map<number | string, string>([
  [403, '401'],
  [404, '404']
])

// Error Redirection
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const errorStatus = error.response?.status ?? 0
    const navigate = useNavigate()
    if (errorCodeToPathname.has(errorStatus)) {
      navigate(errorCodeToPathname.get(errorStatus))
    }
    return error
  }
)

function isAuthRequest(args: string): boolean {
  return typeof args === 'string' && args.includes('auth')
}
