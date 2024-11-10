import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface AuthState {
  songId: string | undefined
}

const initialState: AuthState = {
  songId: undefined
}

export const songsSlice = createSlice({
  name: 'songs',
  initialState,
  reducers: {
    setSongId: (state, action: PayloadAction<string | undefined>) => {
      state.songId = action.payload
    }
  }
})

export const { setSongId } = songsSlice.actions

export default songsSlice.reducer
