import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { ThemeMode } from '@renderer/theme/theme'

export interface GlobalState {
  mode: ThemeMode
  userId: string | undefined
  errorPath: string | undefined
}

const initialState: GlobalState = {
  mode: 'dark',
  userId: undefined,
  errorPath: undefined
}

export const globalSlice = createSlice({
  name: 'global',
  initialState,
  reducers: {
    setMode: (state) => {
      state.mode = state.mode === 'light' ? 'dark' : 'light'
    },
    setUserId: (state, action: PayloadAction<string | undefined>) => {
      state.userId = action.payload
    },
    setErrorPath: (state, action: PayloadAction<string | undefined>) => {
      state.errorPath = action.payload
    }
  }
})

export const { setMode, setUserId, setErrorPath } = globalSlice.actions

export default globalSlice.reducer
