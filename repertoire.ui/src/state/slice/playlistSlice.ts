import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface AuthState {
  songsTotalCount?: number
}

const initialState: AuthState = {
  songsTotalCount: undefined
}

export const playlistSlice = createSlice({
  name: 'playlist',
  initialState,
  reducers: {
    setSongsTotalCount: (state, action: PayloadAction<number | undefined>) => {
      state.songsTotalCount = action.payload
    }
  }
})

export const { setSongsTotalCount } = playlistSlice.actions

export default playlistSlice.reducer
