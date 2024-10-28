import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface AuthState {
  token: string | null
}

const initialState: AuthState = {
  token: localStorage.getItem('access_token')
}

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setToken: (state, action: PayloadAction<string>) => {
      state.token = action.payload
      localStorage.setItem('access_token', state.token)
    },
    signOut: (state) => {
      state.token = null
      localStorage.removeItem('access_token')
    }
  }
})

export const { setToken, signOut } = authSlice.actions

export default authSlice.reducer
