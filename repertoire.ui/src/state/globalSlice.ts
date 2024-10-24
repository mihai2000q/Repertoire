import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface GlobalState {
  userId?: string | undefined
  errorPath?: string | undefined
}

const initialState: GlobalState = {}

export const globalSlice = createSlice({
  name: 'global',
  initialState,
  reducers: {
    setUserId: (state, action: PayloadAction<string | undefined>) => {
      state.userId = action.payload
    },
    setErrorPath: (state, action: PayloadAction<string | undefined>) => {
      state.errorPath = action.payload
    }
  }
})

export const { setUserId, setErrorPath } = globalSlice.actions

export default globalSlice.reducer
