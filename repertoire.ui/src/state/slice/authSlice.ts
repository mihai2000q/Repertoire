import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { api } from '../api.ts'
import { authApi } from '../authApi.ts'

export interface AuthState {
  token: string | null
  historyOnSignIn: { index: number; justSignedIn: boolean }
}

const initialState: AuthState = {
  token: localStorage.getItem('access_token'),
  historyOnSignIn: { index: 0, justSignedIn: false }
}

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setToken: (state, action: PayloadAction<string>) => {
      state.token = action.payload
      localStorage.setItem('access_token', state.token)
    },
    signIn: (state, action: PayloadAction<string>) => {
      state.token = action.payload
      localStorage.setItem('access_token', state.token)

      api.util.resetApiState()
      authApi.util.resetApiState()
      state.historyOnSignIn = { index: (history.state?.idx ?? 0) + 1, justSignedIn: true }
    },
    signOut: (state) => {
      state.token = null
      localStorage.removeItem('access_token')
    },
    resetHistoryOnSignIn: (state) => {
      state.historyOnSignIn = { ...state.historyOnSignIn, justSignedIn: false }
    }
  }
})

export const { setToken, signIn, signOut, resetHistoryOnSignIn } = authSlice.actions

export default authSlice.reducer
