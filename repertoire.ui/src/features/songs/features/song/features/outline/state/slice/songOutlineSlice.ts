import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import OutlineView from '../../types/enums/OutlineView.ts'

interface SongOutlineState {
  view: OutlineView
  showDetails: boolean
}

const initialState: SongOutlineState = {
  view: OutlineView.Sections,
  showDetails: false
}

export const songOutlineSlice = createSlice({
  name: 'songOutline',
  initialState,
  reducers: {
    setView: (state, action: PayloadAction<OutlineView>) => {
      state.view = action.payload
    },
    setShowDetails: (state, action: PayloadAction<boolean>) => {
      state.showDetails = action.payload
    }
  }
})

export const { setView, setShowDetails } = songOutlineSlice.actions

export default songOutlineSlice.reducer
