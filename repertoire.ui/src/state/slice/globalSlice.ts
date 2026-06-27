import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface GlobalState {
  userId?: string | undefined
  documentTitle?: string | undefined
  errorPath?: string | undefined
}

const initialState: GlobalState = {
  userId: undefined
}

export const globalSlice = createSlice({
  name: 'global',
  initialState,
  reducers: {
    setUserId: (state, action: PayloadAction<string>) => {
      state.userId = action.payload
    },
    setDocumentTitle: (state, action: PayloadAction<string | undefined>) => {
      state.documentTitle = action.payload
    },
    setErrorPath: (state, action: PayloadAction<string | undefined>) => {
      state.errorPath = action.payload
    }
  }
})

export const { setUserId, setDocumentTitle, setErrorPath } = globalSlice.actions

export default globalSlice.reducer
