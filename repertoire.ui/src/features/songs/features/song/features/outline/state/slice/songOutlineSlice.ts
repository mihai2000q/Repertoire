import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import OutlineView from '../../types/enums/OutlineView.ts'
import LocalStorageKeys from '../../../../../../../../types/enums/keys/LocalStorageKeys.ts'

interface SongOutlineState {
  view: OutlineView
  showDetails: boolean
}

const initialState: SongOutlineState = {
  view: readPersistedView(),
  showDetails: false
}

export const songOutlineSlice = createSlice({
  name: 'songOutline',
  initialState,
  reducers: {
    setView: (state, action: PayloadAction<OutlineView>) => {
      state.view = action.payload
      localStorage.setItem(LocalStorageKeys.SongOutlineView, JSON.stringify(action.payload))
    },
    setShowDetails: (state, action: PayloadAction<boolean>) => {
      state.showDetails = action.payload
    }
  }
})

function readPersistedView(): OutlineView {
  try {
    const item = localStorage.getItem(LocalStorageKeys.SongOutlineView)
    if (item === null) return OutlineView.Sections
    const parsed = JSON.parse(item)
    if (parsed === OutlineView.Sections || parsed === OutlineView.Parts) return parsed
    return OutlineView.Sections
  } catch {
    return OutlineView.Sections
  }
}

export const { setView, setShowDetails } = songOutlineSlice.actions

export default songOutlineSlice.reducer
